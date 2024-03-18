//go:build windows
// +build windows

package file

import (
	"context"
	"os"
	"time"
)

// Listen for signals to print or not.
func (s *FileSink) loopAndListen(ctx context.Context) {
	dumpTick := time.NewTicker(time.Duration(s.config.FlushIntervalSeconds) * time.Second)
	s.Infof("Writing file at %s %v ...", s.location, s.doWrite)
	defer dumpTick.Stop()

	for {
		select {
		case _ = <-dumpTick.C:
			if !s.doWrite {
				continue
			}

			s.mux.Lock()
			oldName := s.fd.Name()
			if s.fd != nil {
				s.fd.Sync()
				s.fd.Close()
			}
			if s.written == 0 {
				os.Remove(oldName)
			}

			s.written = 0
			name := s.getName()
			f, err := os.Create(name)
			if err != nil {
				s.Errorf("There was an error when creating the %s file: %v.", name, err)
				s.fd = nil
			} else {
				s.fd = f
			}
			s.mux.Unlock()
			s.Debugf("New file: %s", name)

		case <-ctx.Done():
			s.Infof("fileSink Done")
			return
		}
	}
}
