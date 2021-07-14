package util

import (
	"context"
	"fmt"
	"time"

	"github.com/kentik/gosnmp"
	"github.com/kentik/ktranslate/pkg/eggs/logger"
	"github.com/kentik/ktranslate/pkg/kt"
)

const (
	SNMP_PORT       = uint16(161)
	MAX_CONNECT_TRY = 10
)

func parseV3Config(v3config *kt.V3SNMPConfig) (*gosnmp.UsmSecurityParameters, gosnmp.SnmpV3MsgFlags, string, string, error) {

	if v3config == nil {
		return nil, gosnmp.AuthNoPriv, "", "", fmt.Errorf("invalid nil v3 config passed: %v", v3config)
	}

	params := gosnmp.UsmSecurityParameters{
		UserName:                 v3config.UserName,
		AuthenticationPassphrase: v3config.AuthenticationPassphrase,
		PrivacyPassphrase:        v3config.PrivacyPassphrase,
	}

	flags := gosnmp.AuthPriv

	switch v3config.PrivacyProtocol {
	case "NoPriv":
		flags = gosnmp.AuthNoPriv
		params.PrivacyProtocol = gosnmp.NoPriv
	case "DES":
		params.PrivacyProtocol = gosnmp.DES
	case "AES":
		params.PrivacyProtocol = gosnmp.AES
	case "AES192":
		params.PrivacyProtocol = gosnmp.AES192
	case "AES256":
		params.PrivacyProtocol = gosnmp.AES256
	case "AES192C":
		params.PrivacyProtocol = gosnmp.AES192C
	case "AES256C":
		params.PrivacyProtocol = gosnmp.AES256C
	default:
		return nil, gosnmp.AuthNoPriv, "", "", fmt.Errorf("invalid v3 privacy_protocol: %s. valid options: NoPriv|DES|AES|AES192|AES256|AES192C|AES256C", v3config.PrivacyProtocol)
	}

	switch v3config.AuthenticationProtocol {
	case "NoAuth":
		flags = gosnmp.NoAuthNoPriv
		params.AuthenticationProtocol = gosnmp.NoAuth
	case "MD5":
		params.AuthenticationProtocol = gosnmp.MD5
	case "SHA":
		params.AuthenticationProtocol = gosnmp.SHA
	default:
		return nil, gosnmp.AuthNoPriv, "", "", fmt.Errorf("invalid v3 authentication_protocol: %s. valid options: NoAuth|MD5|SHA", v3config.AuthenticationProtocol)
	}

	return &params, flags, v3config.ContextEngineID, v3config.ContextName, nil
}

type logWrapper struct {
	print  func(v ...interface{})
	printf func(format string, v ...interface{})
}

func (l logWrapper) Print(v ...interface{}) {
	l.print(v...)
}

func (l logWrapper) Printf(format string, v ...interface{}) {
	l.printf(format, v...)
}

func InitSNMP(device *kt.SnmpDeviceConfig, connectTimeout time.Duration, retries int, log logger.ContextL) (*gosnmp.GoSNMP, error) {

	if (device.Community == "" && device.V3 == nil) || device.DeviceIP == "" {
		return nil, fmt.Errorf("community or server IP not set")
	}

	// If these are set at a device level, use these instead of the passed in values.
	if device.TimeoutMS > 0 {
		connectTimeout = time.Duration(device.TimeoutMS) * time.Millisecond
		log.Infof("Device Level timeout of %v", connectTimeout)
	}
	if device.Retries > 0 {
		retries = device.Retries
		log.Infof("Device Level retries of %v", retries)
	}

	port := SNMP_PORT
	if device.Port != 0 {
		port = device.Port
	}

	server := &gosnmp.GoSNMP{
		Port:               port,
		Transport:          "udp",
		Timeout:            connectTimeout,
		Retries:            retries,
		Target:             device.DeviceIP,
		MaxOids:            gosnmp.MaxOids,
		ExponentialTimeout: true,
		Context:            context.Background(),
	}

	if device.V3 == nil {
		server.Community = device.Community
		if device.UseV1 {
			log.Infof("Running with SNMP v1")
			server.Version = gosnmp.Version1
		} else {
			log.Infof("Running with SNMP v2c")
			server.Version = gosnmp.Version2c
		}
	} else {
		params, flags, contextEngineID, contextName, err := parseV3Config(device.V3)
		if err != nil {
			return nil, err
		}

		log.Infof("Running with SNMP v3")
		server.Version = gosnmp.Version3
		server.ContextEngineID = contextEngineID
		server.ContextName = contextName
		server.SecurityModel = gosnmp.UserSecurityModel // Only one supported.
		server.MsgFlags = flags
		server.SecurityParameters = params
	}

	if device.Debug {
		server.Logger = logWrapper{
			print: func(v ...interface{}) {
				log.Debugf("GoSNMP:" + fmt.Sprint(v...))
			},
			printf: func(format string, v ...interface{}) {
				log.Debugf("GoSNMP:  "+format, v...)
			},
		}
	}

	// We have everything we need -- start connect.
	times := 0
	for {
		err := connectSNMP(server)
		if err == nil {
			return server, nil
		}

		if times > MAX_CONNECT_TRY {
			return nil, err
		}

		log.Warnf("Could not connect to SNMP -- take %d", times)
		time.Sleep(SNMP_POLL_SLEEP_TIME)
		times++
	}
}

func connectSNMP(x *gosnmp.GoSNMP) error {
	return x.Connect()
}
