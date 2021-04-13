package auth

// Run an auth service, returning auth info needed to run a kproxy/kprobe without talking to kentik.com
import (
	"crypto/rand"
	"encoding/json"
	"io/ioutil"
	"math/big"
	"net"
	"net/http"
	"strconv"

	"github.com/kentik/ktranslate/pkg/eggs/kmux"
	"github.com/kentik/ktranslate/pkg/eggs/logger"
	"github.com/kentik/ktranslate/pkg/snmp"
)

type Server struct {
	Devices     map[string]*Device
	DevicesByIP map[string]*Device
	log         logger.ContextL
}

const (
	API     = "/api"
	TSDB    = "/tsdb"
	API_INT = "/api/internal"
)

func NewServer(deviceFile string, snmpFile string, log logger.ContextL) (*Server, error) {
	devices, err := loadDevices(deviceFile)
	if err != nil {
		return nil, err
	}
	s := &Server{
		log:         log,
		Devices:     devices,
		DevicesByIP: make(map[string]*Device),
	}

	if snmpFile != "" {
		snmp, err := snmp.ParseConfig(snmpFile)
		if err != nil {
			return nil, err
		}

		// Pick a device to copy things from
		var root *Device
		for _, db := range devices {
			root = db
			break
		}
		if root != nil {
			// Now, itterate over this snmp file, adding in all the devices we have listed here.
			nextID := root.ID + 100
			for _, d := range snmp.Devices {
				nd := &Device{
					ID:          nextID,
					Name:        d.DeviceName,
					Type:        root.Type,
					Subtype:     root.Subtype,
					Description: d.Description,
					IP:          net.ParseIP(d.DeviceIP),
					SendingIps:  []net.IP{net.ParseIP(d.DeviceIP)},
					SampleRate:  int(d.SampleRate),
					BgpType:     root.BgpType,
					Plan:        root.Plan,
					CdnAttr:     root.CdnAttr,
					MaxFlowRate: root.MaxFlowRate,
					CompanyID:   root.CompanyID,
					Customs:     root.Customs,
					CustomStr:   root.CustomStr,
				}
				if nd.SampleRate == 0 {
					nd.SampleRate = 1
				}
				s.Devices[strconv.Itoa(nd.ID)] = nd
				nextID += 100
			}
		}
	}

	log.Infof("API server running %d devices", len(s.Devices))
	for _, d := range s.Devices {
		for _, ip := range d.SendingIps {
			s.DevicesByIP[ip.String()] = d
		}
	}

	return s, nil
}

func (s *Server) RegisterRoutes(r *kmux.Router) {
	r.HandleFunc(API+"/device/{did}", s.wrap(s.device))
	r.HandleFunc(API+"/device/", s.wrap(s.create))
	r.HandleFunc(API+"/device/{did}/interfaces", s.wrap(s.interfaces))
	r.HandleFunc(API+"/company/{cid}/device/{did}/tags/snmp", s.wrap(s.update))
	r.HandleFunc(API+"/devices", s.wrap(s.devices))
	r.HandleFunc(API_INT+"/device/{did}", s.wrap(s.device))
	r.HandleFunc(API_INT+"/devices", s.wrap(s.devices))
}

func (s *Server) getDevice(query string) *Device {
	// Try finding this device directly by its ID
	device, ok := s.Devices[query]
	if ok {
		return device
	}

	// Else, can we find this by its IP?
	ipr := net.ParseIP(query)
	if ipr != nil {
		return s.DevicesByIP[ipr.String()]
	}

	return nil
}

func (s *Server) device(w http.ResponseWriter, r *http.Request) {
	id := kmux.Vars(r)["did"]
	device := s.getDevice(id)
	if device == nil {
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
