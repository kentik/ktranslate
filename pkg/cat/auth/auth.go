package auth

// Run an auth service, returning auth info needed to run a kproxy/kprobe without talking to kentik.com
import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/big"
	"net"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/gorilla/mux"

	"github.com/kentik/ktranslate/pkg/eggs/logger"
)

type Server struct {
	Host     net.IP
	Port     int
	Devices  map[string]*Device
	mux      *mux.Router
	listener net.Listener
	log      logger.ContextL
}

const (
	API     = "/api"
	TSDB    = "/tsdb"
	API_INT = "/api/internal"
)

func NewServer(host string, port int, tls bool, deviceFile string, log logger.ContextL) (*Server, error) {
	var listener net.Listener

	addr, err := net.ResolveTCPAddr("tcp", net.JoinHostPort(host, strconv.Itoa(port)))
	if err != nil {
		return nil, err
	}

	listener, err = net.ListenTCP("tcp", addr)
	if err != nil {
		return nil, err
	}

	if tls {
		listener, err = tlslistener(listener, host, addr)
		if err != nil {
			return nil, err
		}
	}

	addr = listener.Addr().(*net.TCPAddr)

	devices, err := loadDevices(deviceFile)
	if err != nil {
		return nil, err
	}

	log.Infof("API server running at %s:%d with %d devices", host, port, len(devices))

	return &Server{
		Host:     addr.IP,
		Port:     addr.Port,
		mux:      mux.NewRouter(),
		listener: listener,
		log:      log,
		Devices:  devices,
	}, nil
}

func (s *Server) Serve() error {
	s.mux.HandleFunc(API+"/device/{did}", s.wrap(s.device))
	s.mux.HandleFunc(API+"/device/", s.wrap(s.create))
	s.mux.HandleFunc(API+"/device/{did}/interfaces", s.wrap(s.interfaces))
	s.mux.HandleFunc(API+"/company/{cid}/device/{did}/tags/snmp", s.wrap(s.update))
	s.mux.HandleFunc(API+"/devices", s.wrap(s.devices))
	s.mux.HandleFunc(API_INT+"/device/{did}", s.wrap(s.device))

	return http.Serve(s.listener, s.mux)
}

func (s *Server) URL(path string) *url.URL {
	url, _ := url.Parse(fmt.Sprintf("http://%s:%d%s", s.Host, s.Port, path))
	return url
}

func (s *Server) device(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["did"]

	device, ok := s.Devices[id]
	if !ok {
		panic(http.StatusNotFound)
	}

	w.Header().Set("Content-Type", "application/json")
	err := json.NewEncoder(w).Encode(&DeviceWrapper{
		Device: device,
	})

	if err != nil {
		panic(http.StatusInternalServerError)
	}

	s.log.Infof("Lookup up device %d", device.ID)
}

func (s *Server) devices(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	devices := []*Device{}
	for _, d := range s.Devices {
		devices = append(devices, d)
	}

	err := json.NewEncoder(w).Encode(&AllDeviceWrapper{
		Devices: devices,
	})

	if err != nil {
		panic(http.StatusInternalServerError)
	}
}

func (s *Server) create(w http.ResponseWriter, r *http.Request) {
	wrapper := map[string]*DeviceCreate{"device": &DeviceCreate{}}

	if err := json.NewDecoder(r.Body).Decode(&wrapper); err != nil {
		panic(http.StatusInternalServerError)
	}

	create := wrapper["device"]

	plan := Plan{
		ID: uint64(create.PlanID),
	}

	var od *Device
	for _, d := range s.Devices {
		od = d
		break
	}

	id, _ := rand.Int(rand.Reader, big.NewInt(65535))
	device := &Device{
		ID:          int(id.Int64()),
		Name:        create.Name,
		Type:        create.Type,
		Description: create.Description,
		IP:          create.IPs[0],
		SampleRate:  create.SampleRate,
		BgpType:     create.BgpType,
		Plan:        plan,
		CdnAttr:     create.CdnAttr,
	}

	if od != nil {
		device.MaxFlowRate = od.MaxFlowRate
		device.CompanyID = od.CompanyID
		device.Customs = od.Customs
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	err := json.NewEncoder(w).Encode(&DeviceWrapper{
		Device: device,
	})

	if err != nil {
		panic(http.StatusInternalServerError)
	}

	s.log.Infof("Created device %d", device.ID)
	s.Devices[create.IPs[0].String()] = device // Save for later
}

func (s *Server) interfaces(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode([]Interface{})
}

func (s *Server) update(w http.ResponseWriter, r *http.Request) {
	// just ignore it
}

func (s *Server) wrap(f handler) handler {
	return func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if r := recover(); r != nil {
				if code, ok := r.(int); ok {
					http.Error(w, http.StatusText(code), code)
					return
				}
				panic(r)
			}
		}()

		if err := r.ParseForm(); err != nil {
			panic(http.StatusBadRequest)
		}

		f(w, r)
	}
}

func tlslistener(tcp net.Listener, host string, addr *net.TCPAddr) (net.Listener, error) {
	pri, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		return nil, err
	}
	pub := &pri.PublicKey

	sn, err := rand.Int(rand.Reader, new(big.Int).Lsh(big.NewInt(1), 128))
	if err != nil {
		return nil, err
	}

	template := x509.Certificate{
		SerialNumber:          sn,
		Subject:               pkix.Name{Organization: []string{"Kentik"}},
		IPAddresses:           []net.IP{addr.IP},
		DNSNames:              []string{host},
		NotBefore:             time.Now(),
		NotAfter:              time.Now().AddDate(10, 0, 0),
		KeyUsage:              x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		BasicConstraintsValid: true,
	}

	der, err := x509.CreateCertificate(rand.Reader, &template, &template, pub, pri)
	if err != nil {
		return nil, err
	}

	cert := tls.Certificate{
		Certificate: [][]byte{der},
		PrivateKey:  pri,
	}

	cfg := tls.Config{Certificates: []tls.Certificate{cert}}
	return tls.NewListener(tcp, &cfg), nil
}

type handler func(http.ResponseWriter, *http.Request)

func loadDevices(file string) (map[string]*Device, error) {
	ms := map[string]*Device{}

	// If the file is empty string, just continue and load 0 devices.
	if file == "" {
		return ms, nil
	}

	// Otherwise, we need to try and process it.
	by, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(by, &ms)
	if err != nil {
		return nil, err
	}

	return ms, nil
}
