package snmp

import (
	"context"
	"os"

	"github.com/kentik/ktranslate/pkg/eggs/logger"
	"github.com/kentik/ktranslate/pkg/kt"
)

var (
	snmpConfPerms os.FileMode = 0664
)

func initOutputFile(ctx context.Context, log logger.ContextL, outputFile string, conf *kt.SnmpConfig, origSnmpFile string) error {
	log.Infof("Writing snmp config file to %s.", outputFile)
	data, err := os.ReadFile(origSnmpFile)
	if err != nil {
		return err
	}
	return os.WriteFile(outputFile, data, snmpConfPerms)
}
