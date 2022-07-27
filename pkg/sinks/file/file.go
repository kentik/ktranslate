package file

import (
	"context"
	"flag"
	"fmt"
	"math/rand"
	"os"
	"os/signal"
	"sync"
	"time"

	go_metrics "github.com/kentik/go-metrics"
	"github.com/kentik/ktranslate"
	"github.com/kentik/ktranslate/pkg/eggs/logger"
	"github.com/kentik/ktranslate/pkg/formats"
	"github.com/kentik/ktranslate/pkg/kt"
)

var (
	fileDir     string
	fileWrite   bool
	flushDurSec int
)

func init() {
	flag.StringVar(&fileDir, "file_out", "./", "Write flows seen to log to this directory if set")
	flag.BoolVar(&fileWrite, "file_on", false, "If true, start writting to file sink right away. Otherwise, wait for a USR1 signal")
	flag.IntVar(&flushDurSec, "file_flush_sec", 60, "Create a new output file every this many seconds")
}

type FileSink struct {
	logger.ContextL
	doWrite  bool
	location string
	fd       *os.File
	mux      sync.RWMutex
	suffix   string
	written  int
	config   *ktranslate.FileSinkConfig
}

func NewSink(log logger.Underlying, registry go_metrics.Registry, cfg *ktranslate.FileSinkConfig) (*FileSink, error) {
	rand.Seed(time.Now().UnixNano())
	fs := &FileSink{
		ContextL: logger.NewContextLFromUnderlying(logger.SContext{S: "fileSink"}, log),
		doWrite:  cfg.EnableImmediateWrite,
		config:   cfg,
	}
	return fs, nil
}

func (s *FileSink) getName() string {
	return fmt.Sprintf("%s/%d_%d%s", s.location, time.Now().Unix(), rand.Intn(100000), s.suffix)
}

func (s *FileSink) Init(ctx context.Context, format formats.Format, compression kt.Compression, fmtr formats.Formatter) error {
	s.location = s.config.Path
	_, err := os.Stat(s.config.Path)
	if err != nil {
		return err
	}

	switch format {
	case formats.FORMAT_JSON, formats.FORMAT_JSON_FLAT, formats.FORMAT_NRM, formats.FORMAT_NR, formats.FORMAT_ELASTICSEARCH:
		s.suffix = ".json"
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
		signal.Notify(sigCh, kt.SIGUSR1)
		dumpTick := time.NewTicker(time.Duration(s.config.FlushIntervalSeconds) * time.Second)
		s.Infof("Writing file at %s %v ...", s.location, s.doWrite)
		defer dumpTick.Stop()

		for {
			select {
			case sig := <-sigCh:
				switch sig {
				case kt.SIGUSR1: // Toggles print. Note -- doesn't work in windows.
					s.doWrite = !s.doWrite
					s.Infof("Writing file at %s %v ...", s.location, s.doWrite)
					if s.doWrite {
						s.mux.Lock()
						name := s.getName()
						f, err := os.Create(name)
						if err != nil {
							s.Errorf("There was an error when creating the %s file: %v.", name, err)
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
	}()

	s.Infof("Writing files to %s, PID=%d", s.location, os.Getpid())
	return nil
}

func (s *FileSink) Send(ctx context.Context, payload *kt.Output) {
	s.mux.Lock()
	defer s.mux.Unlock()
	if s.doWrite && s.fd != nil {
		written, err := s.fd.Write(payload.Body)
		if err != nil {
			s.Infof("Cannot write to %s, %v", s.location, err)
		}
		s.written += written
	}
}

func (s *FileSink) Close() {
	if s.fd != nil {
		oldName := s.fd.Name()
		s.fd.Close()
		if s.written == 0 {
			os.Remove(oldName)
		}
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
