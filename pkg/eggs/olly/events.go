// Package olly provides helpers and wrappers for observability (aka o11y)
//

package olly

import (
	"fmt"
	"sync"
	"time"

	"github.com/kentik/ktranslate/pkg/eggs/concurrent"

	"github.com/google/uuid"
	"github.com/honeycombio/libhoney-go"
)

const (
	CloseFlushTimeout = 1 * time.Second
)

// Builder wraps libhoney in case we want to use something else
type Builder struct {
	*libhoney.Builder
}

// NewBuilder constructs new builder object
func NewBuilder() *Builder {
	return &Builder{Builder: libhoney.NewBuilder()}
}

// NewEvent type wrapper for libhoney Event
func (b *Builder) NewEvent() *Event {
	return &Event{Event: b.Builder.NewEvent()}
}

// Clone convenience wrapper for libhoney.Builder.Clone()
func (b *Builder) Clone() *Builder {
	return &Builder{Builder: b.Builder.Clone()}
}

// Event wraps libhoney.Event in case we want to use something else
type Event struct {
	*libhoney.Event
}

// Context helps provide scope builders
type Context interface {
	OllyBuilder() *Builder
}

// FieldHolder or functions that need to work on both libhoney Builders and Events
type FieldHolder interface {
	AddField(name string, val interface{})
}

var initialized bool

// Init prepares event sending engine
func Init(svcName string, version string, cfg libhoney.Config, data ...interface{}) {
	if initialized {
		return
	}

	if err := libhoney.Init(cfg); err != nil {
		return
	}

	globalFields := data2map(data...)
	globalFields["service_name"] = svcName
	globalFields["version"] = version
	if err := libhoney.Add(globalFields); err != nil {
		panic(fmt.Errorf("olly: init add failed: %+v", err))
	}

	initialized = true
}

// Close shuts down event sending engine
func Close() {
	if initialized {
		wg := &sync.WaitGroup{}
		wg.Add(1)
		go func() {
			libhoney.Close()
			wg.Done()
		}()
		concurrent.WgWaitTimeout(wg, CloseFlushTimeout)
	}
}

// QuickC creates and sends an event in given context
func QuickC(ollyContext Context, op ollyOp, data ...interface{}) {
	Send(newEvent(ollyContext.OllyBuilder(), op, data...))
}

// PrepareC prepares an event in given context
func PrepareC(ollyContext Context, op ollyOp, data ...interface{}) *Event {
	return newEvent(ollyContext.OllyBuilder(), op, data...)
}

// AddContext Adds fields from context to event
func AddContext(ev *Event, ollyContext Context) *Event {
	_ = ev.Add(ollyContext.OllyBuilder().Fields()) //
	return ev
}

// AddUuid inserts newly generated uuid to a builder or event under the given field name
func AddUuid(fh FieldHolder, fieldName string) {
	fh.AddField(fieldName, uuid.New().String())
}

// AddErr adds "err" field with given value to event
func AddErr(event *Event, err error) {
	event.AddField("err", err.Error())
}

// Send an event
func Send(event *Event) {
	if !initialized {
		return
	}
	if err := event.Send(); err != nil {
		// TODO: log
		return
	}
}

func newEvent(builder *Builder, op ollyOp, data ...interface{}) *Event {
	if builder == nil {
		panic("olly: newEvent: builder cannot be nil")
	}

	dataMap := data2map(data...)
	dataMap["op"] = op.string
	ev := builder.NewEvent()
	if err := ev.Add(dataMap); err != nil {
		panic(fmt.Errorf("olly: could not add data: %+v", err))
	}
	return ev
}

func data2map(data ...interface{}) map[string]interface{} {
	if len(data)%2 != 0 {
		panic("olly: newEvent: data must be k/v pairs")
	}
	dataMap := make(map[string]interface{}, len(data)/2+1)
	for i := 0; i < len(data); i += 2 {
		if k, ok := data[i].(string); ok {
			dataMap[k] = data[i+1]
		}
	}
	return dataMap
}
