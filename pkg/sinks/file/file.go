package file

import (
	"context"
	"flag"
	"fmt"
	"math/rand"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	go_metrics "github.com/kentik/go-metrics"
	"github.com/kentik/ktranslate/pkg/eggs/logger"
	"github.com/kentik/ktranslate/pkg/formats"
	"github.com/kentik/ktranslate/pkg/kt"
)

type FileSink struct {
	logger.ContextL
	doWrite  bool
	location string
	fd       *os.File
	mux      sync.RWMutex
}

var (
	FileDir     = flag.String("file_out", "./", "Write flows seen to log to this directory if set")
	FileWrite   = flag.Bool("file_on", false, "If true, start writting to file sink right away. Otherwise, wait for a USR1 signal")
	FlushDurSec = flag.Int("file_flush_sec", 60, "Create a new output file every this many seconds")
)

func NewSink(log logger.Underlying, registry go_metrics.Registry) (*FileSink, error) {
	rand.Seed(time.Now().UnixNano())
	fs := &FileSink{
		ContextL: logger.NewContextLFromUnderlying(logger.SContext{S: "fileSink"}, log),
		doWrite:  *FileWrite,
	}
	return fs, nil
}

func (s *FileSink) getName() string {
	return fmt.Sprintf("%s/%d_%d", s.location, time.Now().Unix(), rand.Intn(100000))
}

func (s *FileSink) Init(ctx context.Context, format formats.Format, compression kt.Compression, fmtr formats.Formatter) error {
	s.location = *FileDir
	_, err := os.Stat(*FileDir)
	if err != nil {
		return err
	}

	// Set up a file first.
	if s.doWrite {
		name := s.getName()
		f, err := os.Create(name)
		if err != nil {
			return err
		}
		s.fd = f
	}

	go func() {
		// Listen for signals to print or not.
		sigCh := make(chan os.Signal, 2)
		signal.Notify(sigCh, syscall.SIGUSR1)
		dumpTick := time.NewTicker(time.Duration(*FlushDurSec) * time.Second)
		s.Infof("File out -- Write is now: %v, dumping on %v", s.doWrite, time.Duration(*FlushDurSec)*time.Second)
		defer dumpTick.Stop()

		for {
			select {
			case sig := <-sigCh:
				switch sig {
				case syscall.SIGUSR1: // Toggles print.
					s.doWrite = !s.doWrite
					s.Infof("Write is now: %v", s.doWrite)
					if s.doWrite {
						s.mux.Lock()
						name := s.getName()
						f, err := os.Create(name)
						if err != nil {
							s.Errorf("Cannot create file %s -> %v", name, err)
						} else {
							s.fd = f
						}
						s.mux.Unlock()
					}
				}

			case _ = <-dumpTick.C:
				if !s.doWrite {
					continue
				}

				s.mux.Lock()
				if s.fd != nil {
					s.fd.Sync()
					s.fd.Close()
				}
				name := s.getName()
				f, err := os.Create(name)
				if err != nil {
					s.Errorf("Cannot create file %s -> %v", name, err)
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
	}()

	s.Infof("Writing files to %s, PID=%d", s.location, os.Getpid())
	return nil
}

func (s *FileSink) Send(ctx context.Context, payload *kt.Output) {
	s.mux.Lock()
	defer s.mux.Unlock()
	if s.doWrite && s.fd != nil {
		_, err := s.fd.Write(payload.Body)
		if err != nil {
			s.Infof("Cannot write to %s, %v", s.location, err)
		}
	}
}

func (s *FileSink) Close() {
	if s.fd != nil {
		s.fd.Close()
	}
}

func (s *FileSink) HttpInfo() map[string]float64 {
	doWrite := float64(0.)
	if s.doWrite {
		doWrite = 1.0
	}

	return map[string]float64{
		"Write": doWrite,
	}
}
