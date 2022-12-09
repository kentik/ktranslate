package cat

import (
	"github.com/kentik/ktranslate"
)

// Callback for when theres a config managment service which detects a change.
func (kc *KTranslate) newConfig(newC *ktranslate.Config) error {
	kc.log.Warnf("Write detected on %s, shutting down", newC.Server.CfgPath)
	kc.shutdown("Config file changed")
	return nil
}
