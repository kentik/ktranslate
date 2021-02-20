package logger

import (
	"fmt"
)

// At the lowest level we have a logger.Underlying.

// Implemented by github.com/kentik/golog/logger
// Don't use this directly.
type Underlying interface {
	Debugf(lp string, f string, params ...interface{})
	Infof(lp string, f string, params ...interface{})
	Warnf(lp string, f string, params ...interface{})
	Errorf(lp string, f string, params ...interface{})
}

// ---------------
// Built on top of a logger.Underlying, we have a logger.L, which
// as a convenience lets you pass an interface in that becomes
// the log prefix.

// Logs and lets you specify the context.
type L interface {
	Debugf(lc Context, f string, params ...interface{})
	Infof(lc Context, f string, params ...interface{})
	Warnf(lc Context, f string, params ...interface{})
	Errorf(lc Context, f string, params ...interface{})

	GetUnderlyingLogger() Underlying
}

type LoggerImpl struct {
	UL Underlying
}

func (l *LoggerImpl) Debugf(c Context, f string, params ...interface{}) {
	l.UL.Debugf(c.GetLogPrefix()+" ", f, params...)
}
func (l *LoggerImpl) Infof(c Context, f string, params ...interface{}) {
	l.UL.Infof(c.GetLogPrefix()+" ", f, params...)
}
func (l *LoggerImpl) Warnf(c Context, f string, params ...interface{}) {
	l.UL.Warnf(c.GetLogPrefix()+" ", f, params...)
}
func (l *LoggerImpl) Errorf(c Context, f string, params ...interface{}) {
	l.UL.Errorf(c.GetLogPrefix()+" ", f, params...)
}
func (l *LoggerImpl) GetUnderlyingLogger() Underlying {
	return l.UL
}

// You need a Context to pass to the Logger.

type Context interface {
	GetLogPrefix() string
}

var NilContext Context = &NilContextS{}

type NilContextS struct{}

func (nc NilContextS) GetLogPrefix() string { return "" }

// A barebones implementation of a Context.
type SContext struct {
	S string
}

func (slc SContext) GetLogPrefix() string {
	return slc.S
}

func NewSubContext(parent Context, child Context) Context {
	return &SubContext{parent: parent, child: child}
}

// For nested, hierarchical Contexts.
type SubContext struct{ parent, child Context }

func (slc *SubContext) GetLogPrefix() string {
	return fmt.Sprintf("%s>%s", slc.parent.GetLogPrefix(), slc.child.GetLogPrefix())
}

// --------------------
// Built on top of a logger.L (on top of a logger.Underlying), we have
// a logger.ContextL, which already has a Context associated with it.
// It passes the Context to the L for you.

// Logs within its Context.
type ContextL interface {
	Context
	Debugf(f string, params ...interface{})
	Infof(f string, params ...interface{})
	Warnf(f string, params ...interface{})
	Errorf(f string, params ...interface{})

	// Methods accepting one-time subcontexts
	SDebugf(lc Context, f string, params ...interface{})
	SInfof(lc Context, f string, params ...interface{})
	SWarnf(lc Context, f string, params ...interface{})
	SErrorf(lc Context, f string, params ...interface{})

	GetLogger() L
	GetLogContext() Context
}

func NewContextL(lc Context, cl ContextL) ContextL {
	return &ContextLImpl{
		Context: lc,
		L:       cl.GetLogger(),
	}
}

type ContextLImpl struct {
	Context
	L L
}

func (l *ContextLImpl) Debugf(f string, params ...interface{}) { l.L.Debugf(l, f, params...) }
func (l *ContextLImpl) Infof(f string, params ...interface{})  { l.L.Infof(l, f, params...) }
func (l *ContextLImpl) Warnf(f string, params ...interface{})  { l.L.Warnf(l, f, params...) }
func (l *ContextLImpl) Errorf(f string, params ...interface{}) { l.L.Errorf(l, f, params...) }
func (l *ContextLImpl) SDebugf(lc Context, f string, params ...interface{}) {
	l.L.Debugf(NewSubContext(l, lc), f, params...)
}
func (l *ContextLImpl) SInfof(lc Context, f string, params ...interface{}) {
	l.L.Infof(NewSubContext(l, lc), f, params...)
}
func (l *ContextLImpl) SWarnf(lc Context, f string, params ...interface{}) {
	l.L.Warnf(NewSubContext(l, lc), f, params...)
}
func (l *ContextLImpl) SErrorf(lc Context, f string, params ...interface{}) {
	l.L.Errorf(NewSubContext(l, lc), f, params...)
}
func (l *ContextLImpl) GetLogger() L           { return l.L }
func (l *ContextLImpl) GetLogContext() Context { return l.Context }

func NewSubContextL(subContext Context, cl ContextL) ContextL {
	return &ContextLImpl{
		Context: NewSubContext(cl.GetLogContext(), subContext),
		L:       cl.GetLogger(),
	}
}

func NewContextLFromUnderlying(lc Context, ul Underlying) ContextL {
	return &ContextLImpl{
		Context: lc,
		L:       &LoggerImpl{UL: ul},
	}
}


