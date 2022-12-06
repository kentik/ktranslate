package config

/**
Interface to manage configs.
*/

import (
	"context"
	"flag"
	"fmt"

	"github.com/kentik/ktranslate"
	"github.com/kentik/ktranslate/pkg/config/local"
	"github.com/kentik/ktranslate/pkg/config/nr"
	"github.com/kentik/ktranslate/pkg/eggs/logger"
	"github.com/kentik/ktranslate/pkg/kt"
)

type ConfigManager interface {
	Run(context.Context, func(*ktranslate.Config) error) // Run takes a context and a callback function to call whenever there is a new update to process
	DeviceDiscovery(kt.DeviceMap)                        // called whenever there is a new snmp device discovery to parse.
	Close()                                              // Called on shutdown of ktrans.
}

type ConfigProvider string

const (
	NewRelicConfig ConfigProvider = "new_relic"
	LocalConfig    ConfigProvider = "local"
	NoConfig       ConfigProvider = ""
)

var (
	configProvider string
)

func init() {
	flag.StringVar(&configProvider, "config_provider", "", "Implementation of which provider controls the config process. Can be one of (new_relic,local)")
}

func NewConfig(prov ConfigProvider, log logger.Underlying, config *ktranslate.Config) (ConfigManager, error) {
	switch prov {
	case NewRelicConfig:
		return nr.NewConfig(log, config)
	case LocalConfig:
		return local.NewConfig(log, config)
	case NoConfig:
		return nil, nil
	default:
		return nil, fmt.Errorf("Unknown config provider %v", prov)
	}
}
