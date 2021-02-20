package file

import (
	"context"
	"flag"
	"fmt"
	"io/ioutil"
	"math/rand"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/kentik/ktranslate/pkg/eggs/logger"
	go_metrics "github.com/kentik/go-metrics"
	"github.com/kentik/ktranslate/pkg/formats"
	"github.com/kentik/ktranslate/pkg/kt"
)

type FileSink struct {
	logger.ContextL
	doWrite  bool
	location string
}

var (
	FileDir   = flag.String("file_out", "./", "Write flows seen to log to this directory if set")
	FileWrite = flag.Bool("file_on", false, "If true, start writting to file sink right away. Otherwise, wait for a USR1 signal")
)

func NewSink(log logger.Underlying, registry go_metrics.Registry) (*FileSink, error) {
	rand.Seed(time.Now().UnixNano())
	return &FileSink{
		ContextL: logger.NewContextLFromUnderlying(logger.SContext{S: "fileSink"}, log),
		doWrite:  *FileWrite,
	}, nil
}

func (s *FileSink) Init(ctx context.Context, format formats.Format, compression kt.Compression) error {
	s.location = *FileDir

	s.Infof("File out -- Write is now: %v", s.doWrite)

	go func() {
		// Listen for signals to print or not.
		sigCh := make(chan os.Signal, 2)
		signal.Notify(sigCh, syscall.SIGUSR1)

		for {
			select {
			case sig := <-sigCh:
				switch sig {
				case syscall.SIGUSR1: // Toggles print.
					s.doWrite = !s.doWrite
					s.Infof("Write is now: %v", s.doWrite)
				}
			case <-ctx.Done():
				s.Infof("fileSink Done")
				return
			}
		}
	}()

	s.Infof("Writing files to %s, PID=%d", s.location, os.Getpid())
	return nil
}

func (s *FileSink) Send(ctx context.Context, payload []byte) {
	if s.doWrite {
		err := ioutil.WriteFile(fmt.Sprintf("%s/%d_%d", s.location, time.Now().Unix(), rand.Intn(100000)), payload, 0644)
		if err != nil {
			s.Infof("Cannot write to %s, %v", s.location, err)
		}
	}
}

func (s *FileSink) Close() {}

func (s *FileSink) HttpInfo() map[string]float64 {
	doWrite := float64(0.)
	if s.doWrite {
		doWrite = 1.0
	}

	return map[string]float64{
		"Write": doWrite,
	}
}
