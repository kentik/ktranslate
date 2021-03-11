package auth

// Run an auth service, returning auth info needed to run a kproxy/kprobe without talking to kentik.com
import (
	"crypto/rand"
	"encoding/json"
	"io/ioutil"
	"math/big"
	"net/http"

	"github.com/kentik/ktranslate/pkg/eggs/kmux"
	"github.com/kentik/ktranslate/pkg/eggs/logger"
)

type Server struct {
	Devices map[string]*Device
	log     logger.ContextL
}

const (
	API     = "/api"
	TSDB    = "/tsdb"
	API_INT = "/api/internal"
)

func NewServer(deviceFile string, log logger.ContextL) (*Server, error) {
	devices, err := loadDevices(deviceFile)
	if err != nil {
		return nil, err
	}

	log.Infof("API server running %d devices", len(devices))

	return &Server{
		log:     log,
		Devices: devices,
	}, nil
}

func (s *Server) RegisterRoutes(r *kmux.Router) {
	r.HandleFunc(API+"/device/{did}", s.wrap(s.device))
	r.HandleFunc(API+"/device/", s.wrap(s.create))
	r.HandleFunc(API+"/device/{did}/interfaces", s.wrap(s.interfaces))
	r.HandleFunc(API+"/company/{cid}/device/{did}/tags/snmp", s.wrap(s.update))
	r.HandleFunc(API+"/devices", s.wrap(s.devices))
	r.HandleFunc(API_INT+"/device/{did}", s.wrap(s.device))
}

func (s *Server) device(w http.ResponseWriter, r *http.Request) {
	id := kmux.Vars(r)["did"]

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
