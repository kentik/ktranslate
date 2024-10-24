package cat

import (
	"context"
)

func (kc *KTranslate) splitLogsForSinks(ctx context.Context) {
	if kc.logTee == nil || len(kc.logTeeSinks) == 0 {
		return
	}
	kc.log.Infof("splitLogsForSinks running with %d splits", len(kc.logTeeSinks))

	for {
		select {
		case log := <-kc.logTee:
			for _, c := range kc.logTeeSinks {
				c <- log
			}
		case <-ctx.Done():
			kc.log.Infof("splitLogsForSinks Done")
			return
		}
	}
}
