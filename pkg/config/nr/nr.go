package nr

import (
	"context"
	"time"

	"github.com/kentik/ktranslate"
	"github.com/kentik/ktranslate/pkg/eggs/logger"
	"github.com/kentik/ktranslate/pkg/kt"
)

type NRConfig struct {
	logger.ContextL
	currentConfig *ktranslate.Config
}

func NewConfig(log logger.Underlying, cfg *ktranslate.Config) (*NRConfig, error) {
	nr := NRConfig{
		ContextL:      logger.NewContextLFromUnderlying(logger.SContext{S: "nrConfig"}, log),
		currentConfig: cfg,
	}

	return &nr, nil
}

func (nr *NRConfig) Run(ctx context.Context, cb func(*ktranslate.Config) error) {
	checkTicker := time.NewTicker(time.Second * time.Duration(nr.currentConfig.CfgManager.PollTimeSec))
	defer checkTicker.Stop()

	nr.Infof("config checker running")
	for {
		select {
		case <-checkTicker.C:
			// Get config
			newConfig, newVersion, err := nr.getConfig()
			if err != nil {
				nr.Errorf("Cannot load new config: %v", err)
			}

			if newConfig != nil && newVersion {
				nr.currentConfig = newConfig
				err := cb(newConfig)
				if err != nil {
					nr.Errorf("Cannot update to new config: %v", err)
				}
			}
		case <-ctx.Done():
			nr.Infof("config checker done")
			return
		}
	}
}

func (nr *NRConfig) DeviceDiscovery(devices kt.DeviceMap) {

}

func (nr *NRConfig) Close() {

}

func (nr *NRConfig) getConfig() (*ktranslate.Config, bool, error) {
	return nil, false, nil
}
