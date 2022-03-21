package http

import (
	"github.com/aristanetworks/goeapi"
	"github.com/aristanetworks/goeapi/module"

	"github.com/kentik/ktranslate/pkg/eggs/logger"
)

type MyShowVlan struct {
	SourceDetail string          `json:"sourceDetail"`
	Vlans        map[string]Vlan `json:"vlans"`
}

type Vlan struct {
	Status     string               `json:"status"`
	Name       string               `json:"name"`
	Interfaces map[string]Interface `json:"interfaces"`
	Dynamic    bool                 `json:"dynamic"`
}

type Interface struct {
	Annotation      string `json:"annotation"`
	PrivatePromoted bool   `json:"privatePromoted"`
}

func (s *MyShowVlan) GetCmd() string {
	return "show vlan configured-ports"
}

type MyMlag struct {
	SourceDetail string `json:"sourceDetail"`
}

func (s *MyMlag) GetCmd() string {
	return "show mlag detail"
}

func connect(log logger.ContextL) error {
	// connect to our device
	node, err := goeapi.ConnectTo("arista1")
	if err != nil {
		return err
	}

	/**
	// get the running config and print it
	conf := node.RunningConfig()
	log.Infof("Running Config:\n%s\n", conf)

	// get api system module
	sys := module.System(node)
	// change the host name to "Ladie"
	if ok := sys.SetHostname("Ladie"); !ok {
		log.Infof("SetHostname Failed\n")
	}
	// get system info
	sysInfo := sys.Get()
	log.Infof("Sysinfo: %#v\n", sysInfo.HostName())
	*/

	sv := &MyShowVlan{}
	ml := &MyMlag{}

	handle, _ := node.GetHandle("json")
	handle.AddCommand(sv)
	handle.AddCommand(ml)
	if err := handle.Call(); err != nil {
		panic(err)
	}

	for k, v := range sv.Vlans {
		log.Infof("Vlan:%s\n", k)
		log.Infof("  Name  : %s\n", v.Name)
		log.Infof("  Status: %s\n", v.Status)
	}

	log.Infof("XXX %v", ml.SourceDetail)

	// get api bpg module
	bgp := module.Show(node)
	sum, err := bgp.ShowIPBGPSummary()
	log.Infof("XXX %v %v", sum, err)

	// get api mlag module
	mlag := module.Mlag(node)
	log.Infof("XXX %v", mlag.Get())

	env, err := bgp.ShowEnvironmentPower()
	log.Infof("XXX %v %v", env, err)

	return nil
}
