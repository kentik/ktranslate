package logger

import (
	"bytes"
	"context"
	"io"
	"regexp"
	"testing"
)

// TestNilLogger tests that you can safely call log methods on a nil logger.
// This is convenient, for example, when you'd like to test code without
// creating and passing in a logger.
func TestNilLogger(t *testing.T) {
	var log *Logger

	log.Debugf("prefix", "Hello %s", "there")
	log.Infof("prefix", "Hello %s", "there")
	log.Warnf("prefix", "Hello %s", "there")
	log.Errorf("prefix", "Hello %s", "there")
	log.Panicf("prefix", "Hello %s", "there")
}

func TestRemoveNewline(t *testing.T) {
	defer func(origstdhdl io.Writer) { stdhdl = origstdhdl }(stdhdl)
	buf := bytes.Buffer{}
	stdhdl = &buf

	log := New(Levels.Debug)

	log.Debugf("", "testing")
	log.Debugf("", "testing\n")
	log.Debugf("", "testing\n\n\n\n")

	Drain()

	if !regexp.MustCompile("^[^\n]*testing\n[^\n]*testing\n[^\n]*testing\n$").Match(buf.Bytes()) {
		t.Error("Expected testing\\n * 3")
	}
}

func TestClose(t *testing.T) {
	defer func() { setup() }() // Set everything up again since we call Close()
	buf := bytes.Buffer{}
	stdhdl = &buf

	log := New(Levels.Debug)
	log.Debugf("", "testing123")
	Close(context.Background())
	if !regexp.MustCompile("^[^\n]*testing123\n$").Match(buf.Bytes()) {
		t.Fatalf("Expected testing123\\n to be written but got '%s'", buf.String())
	}

	didRecover := false
	logAndRecover := func() {
		defer func() {
			if r := recover(); r != nil {
				didRecover = true
			}
		}()
		log.Debugf("", "asdf")
	}
	logAndRecover()
	if !didRecover {
		t.Fatalf("Expected log to panic, but it didn't.")
	}
}
