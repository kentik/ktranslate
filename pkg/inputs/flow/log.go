package flow

import (
	"fmt"

	"github.com/kentik/ktranslate/pkg/eggs/logger"
)

type KentikLog struct {
	logger.ContextL
}

func (l *KentikLog) Printf(f string, vars ...interface{}) {
	l.Infof(f, vars)
}
func (l *KentikLog) Errorf(f string, vars ...interface{}) {
	l.Errorf(f, vars)
}
func (l *KentikLog) Warnf(f string, vars ...interface{}) {
	l.Warnf(f, vars)
}
func (l *KentikLog) Warn(vars ...interface{}) {
	l.Warnf("%v", vars)
}
func (l *KentikLog) Error(vars ...interface{}) {
	l.Errorf("%v", vars)
}
func (l *KentikLog) Debug(vars ...interface{}) {
	l.Debugf("%v", vars)
}
func (l *KentikLog) Debugf(f string, vars ...interface{}) {
	l.Debugf(f, vars)
}
func (l *KentikLog) Infof(f string, vars ...interface{}) {
	l.Infof(f, vars)
}
func (l *KentikLog) Fatalf(f string, vars ...interface{}) {
	l.Errorf(f, vars)
	panic(fmt.Sprintf(f, vars...))
}
