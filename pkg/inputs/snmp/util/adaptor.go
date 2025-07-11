package util

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/gosnmp/gosnmp"
	"github.com/kentik/ktranslate/pkg/eggs/logger"
	"github.com/kentik/ktranslate/pkg/kt"
)

const (
	SNMP_PORT       = uint16(161)
	MAX_CONNECT_TRY = 10
)

func ParseV3Config(v3config *kt.V3SNMPConfig) (*gosnmp.UsmSecurityParameters, gosnmp.SnmpV3MsgFlags, string, string, error) {

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
	case "SHA224":
		params.AuthenticationProtocol = gosnmp.SHA224
	case "SHA256":
		params.AuthenticationProtocol = gosnmp.SHA256
	case "SHA384":
		params.AuthenticationProtocol = gosnmp.SHA384
	case "SHA512":
		params.AuthenticationProtocol = gosnmp.SHA512
	default:
		return nil, gosnmp.AuthNoPriv, "", "", fmt.Errorf("invalid v3 authentication_protocol: %s. valid options: NoAuth|MD5|SHA|SHA224|SHA256|SHA384|SHA512", v3config.AuthenticationProtocol)
	}

	// If there is an invalid AWS defined string, error here.
	if strings.HasPrefix(params.AuthenticationPassphrase, kt.AWSErrPrefix) {
		return nil, gosnmp.AuthNoPriv, "", "", fmt.Errorf("Invalid AuthenticationPassphrase: %s", params.AuthenticationPassphrase)
	}
	if strings.HasPrefix(params.PrivacyPassphrase, kt.AWSErrPrefix) {
		return nil, gosnmp.AuthNoPriv, "", "", fmt.Errorf("Invalid PrivacyPassphrase: %s", params.PrivacyPassphrase)
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

func InitSNMP(device *kt.SnmpDeviceConfig, connectTimeout time.Duration, retries int, posit string, log logger.ContextL) (*gosnmp.GoSNMP, error) {

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
			log.Infof("%s Running with SNMP v1", posit)
			server.Version = gosnmp.Version1
			device.NoUseBulkWalkAll = true // v1 snmp doesn't support this.
		} else {
			log.Infof("%s Running with SNMP v2c", posit)
			server.Version = gosnmp.Version2c
		}
	} else {
		params, flags, contextEngineID, contextName, err := ParseV3Config(device.V3)
		if err != nil {
			return nil, err
		}

		log.Infof("%s Running with SNMP v3: Priv: %s Auth: %s", posit, params.PrivacyProtocol, params.AuthenticationProtocol)
		server.Version = gosnmp.Version3
		server.ContextEngineID = contextEngineID
		server.ContextName = contextName
		server.SecurityModel = gosnmp.UserSecurityModel // Only one supported.
		server.MsgFlags = flags
		server.SecurityParameters = params
	}

	if device.Debug {
		server.Logger = gosnmp.NewLogger(logWrapper{
			print: func(v ...interface{}) {
				log.Debugf("GoSNMP: [hostname=" + device.DeviceName + "] " + fmt.Sprint(v...))
			},
			printf: func(format string, v ...interface{}) {
				log.Debugf("GoSNMP: [hostname="+device.DeviceName+"] "+format, v...)
			},
		})
	}

	if device.NoCheckIncreasing {
		server.AppOpts = map[string]interface{}{"c": true}
		log.Warnf("Turning off Increasing Oid Check. This can result in an infanite loop.")
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

		log.Warnf("It was not possible to connect to the SNMP devices after %d times.", times)
		time.Sleep(SNMP_POLL_SLEEP_TIME)
		times++
	}
}

func connectSNMP(x *gosnmp.GoSNMP) error {
	return x.Connect()
}
