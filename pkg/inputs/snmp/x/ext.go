package x

import (
	"context"
	"time"

	"github.com/kentik/ktranslate/pkg/eggs/logger"
	"github.com/kentik/ktranslate/pkg/inputs/snmp/x/arista"
	"github.com/kentik/ktranslate/pkg/inputs/snmp/x/meraki"
	"github.com/kentik/ktranslate/pkg/kt"
)

// Code to handle various vendor extensions to snmp.
type Extension interface {
	Run(context.Context, time.Duration)
	GetName() string
}

func NewExtension(jchfChan chan []*kt.JCHF, conf *kt.SnmpDeviceConfig, metrics *kt.SnmpDeviceMetric, log logger.ContextL) (Extension, error) {
	if conf.Ext == nil { // No extensions set.
		return nil, nil
	}

	if conf.Ext.EAPIConfig != nil {
		return arista.NewEAPIClient(jchfChan, conf, metrics, log)
	} else if conf.Ext.MerakiConfig != nil {
		return meraki.NewMerakiClient(jchfChan, conf, metrics, log)
	}

	return nil, nil
}
