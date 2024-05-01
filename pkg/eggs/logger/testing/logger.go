package testing

import (
	"fmt"
	"testing"

	"github.com/kentik/ktranslate/pkg/eggs/logger"
)

// Testing implementations:

func NewTestContextL(lc logger.Context, t *testing.T) logger.ContextL {
	return logger.NewContextLFromUnderlying(lc, &testUnderlyingL{logf: t.Logf})
}

func NewBenchContextL(lc logger.Context, b *testing.B) logger.ContextL {
	return logger.NewContextLFromUnderlying(lc, &testUnderlyingL{logf: b.Logf})
}

// Implements logger.Underlying
type Test struct {
	T *testing.T
}

func (l *Test) Debugf(lp string, f string, params ...interface{}) {
	l.T.Logf("%s DEBUG %s", lp, fmt.Sprintf(f, params...))
}
func (l *Test) Infof(lp string, f string, params ...interface{}) {
	l.T.Logf("%s INFO %s", lp, fmt.Sprintf(f, params...))
}
func (l *Test) Warnf(lp string, f string, params ...interface{}) {
	l.T.Logf("%s WARN %s", lp, fmt.Sprintf(f, params...))
}
func (l *Test) Errorf(lp string, f string, params ...interface{}) {
	l.T.Logf("%s ERROR %s", lp, fmt.Sprintf(f, params...))
}

// testUnderlyingL implements logger.Underlying slightly different way
type testUnderlyingL struct {
	logf func(string, ...interface{})
}

func (l *testUnderlyingL) Debugf(lp string, f string, params ...interface{}) {
	l.logf("%s [DEBUG] %s", lp, fmt.Sprintf(f, params...))
}
func (l *testUnderlyingL) Infof(lp string, f string, params ...interface{}) {
	l.logf("%s [INFO] %s", lp, fmt.Sprintf(f, params...))
}
func (l *testUnderlyingL) Warnf(lp string, f string, params ...interface{}) {
	l.logf("%s [WARN] %s", lp, fmt.Sprintf(f, params...))
}
func (l *testUnderlyingL) Errorf(lp string, f string, params ...interface{}) {
	l.logf("%s [ERROR] %s", lp, fmt.Sprintf(f, params...))
}
func (l *testUnderlyingL) GetLogLevel() string {
	return "debug"
}
