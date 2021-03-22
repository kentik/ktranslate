package apc

import (
	"encoding/xml"
	"os"

	"github.com/kentik/ktranslate/pkg/eggs/logger"
)

func ParseApcMib(file string, log logger.ContextL) (*APC, error) {
	data, err := os.ReadFile(file)
	if err != nil {
		return nil, err
	}

	var apc APC
	err = xml.Unmarshal(data, &apc)
	if err != nil {
		return nil, err
	}

	return &apc, nil
}
