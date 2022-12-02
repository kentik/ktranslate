package cat

import (
	"github.com/kentik/ktranslate"

	"github.com/fsnotify/fsnotify"
)

func (kc *KTranslate) monitorConf(cfg *ktranslate.ServerConfig, shutdown func(string)) error {
	kc.log.Infof("Monitoring %s for changes", cfg.CfgPath)

	// Create new watcher.
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return err
	}
	defer watcher.Close()

	// Start listening for events.
	go func() {
		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					return
				}
				if event.Has(fsnotify.Write) {
					kc.log.Warnf("Write detected on %s, shutting down", cfg.CfgPath)
					shutdown("Config file changed")
				}
			case err, ok := <-watcher.Errors:
				if !ok {
					return
				}
				kc.log.Infof("error:", err)
			}
		}
	}()

	// Add a path.
	err = watcher.Add(cfg.CfgPath)
	if err != nil {
		return err
	}

	// Block forever.
	<-make(chan struct{})

	return nil
}
