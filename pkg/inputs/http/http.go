package http

import (
	"compress/gzip"
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"

	go_metrics "github.com/kentik/go-metrics"

	"github.com/kentik/ktranslate/pkg/api"
	"github.com/kentik/ktranslate/pkg/eggs/kmux"
	"github.com/kentik/ktranslate/pkg/eggs/logger"
	"github.com/kentik/ktranslate/pkg/kt"
)

type KentikHttpListener struct {
	logger.ContextL
	metrics  HttpListenerMetric
	apic     *api.KentikApi
	devices  map[string]*kt.Device
	jchfChan chan []*kt.JCHF
}

type HttpListenerMetric struct {
	Messages go_metrics.Meter
	Errors   go_metrics.Meter
}

const (
	CHAN_SLACK           = 10000
	DeviceUpdateDuration = 1 * time.Hour
	Listen               = "/input"
)

func NewHttpListener(ctx context.Context, host string, log logger.Underlying, registry go_metrics.Registry, jchfChan chan []*kt.JCHF, apic *api.KentikApi) (*KentikHttpListener, error) {
	ks := KentikHttpListener{
		ContextL: logger.NewContextLFromUnderlying(logger.SContext{S: "Http"}, log),
		jchfChan: jchfChan,
		metrics: HttpListenerMetric{
			Messages: go_metrics.GetOrRegisterMeter(fmt.Sprintf("http_messages^force=true"), registry),
			Errors:   go_metrics.GetOrRegisterMeter(fmt.Sprintf("http_errors^force=true"), registry),
		},
		apic:    apic,
		devices: apic.GetDevicesAsMap(0),
	}

	go ks.run(ctx)
	return &ks, nil
}

func (ks *KentikHttpListener) RegisterRoutes(r *kmux.Router) {
	r.HandleFunc(Listen+"/telegraf/standard", ks.wrap(ks.readStandard))
	r.HandleFunc(Listen+"/telegraf/batch", ks.wrap(ks.readBatch))
	r.HandleFunc(Listen+"/ktranslate/jchf", ks.wrap(ks.readJCHF))
}

type basic struct {
	Fields    map[string]float64 `json:"fields"`
	Name      string             `json:"name"`
	Tags      map[string]string  `json:"tags"`
	Timestamp int64              `json:"timestamp"`
}

type batch struct {
	Metrics []basic `json:"metrics"`
}

func (ks *KentikHttpListener) readBatch(w http.ResponseWriter, r *http.Request) {
	var wrapper batch

	// Decode body in gzip format if the request header is set this way.
	body := r.Body
	if r.Header.Get("Content-Encoding") == "gzip" {
		z, err := gzip.NewReader(r.Body)
		if err != nil {
			panic(http.StatusInternalServerError)
		}
		body = z
	}
	defer body.Close()

	if err := json.NewDecoder(body).Decode(&wrapper); err != nil {
		panic(http.StatusInternalServerError)
	}
	w.WriteHeader(http.StatusOK)

	remoteIP := getIP(r)
	out := make([]*kt.JCHF, len(wrapper.Metrics))
	for i, m := range wrapper.Metrics {
		out[i] = ks.getJCHF(&m, remoteIP)
	}

	ks.metrics.Messages.Mark(int64(len(out)))
	ks.jchfChan <- out
}

func (ks *KentikHttpListener) readStandard(w http.ResponseWriter, r *http.Request) {
	var wrapper basic

	// Decode body in gzip format if the request header is set this way.
	body := r.Body
	if r.Header.Get("Content-Encoding") == "gzip" {
		z, err := gzip.NewReader(r.Body)
		if err != nil {
			panic(http.StatusInternalServerError)
		}
		body = z
	}
	defer body.Close()

	if err := json.NewDecoder(body).Decode(&wrapper); err != nil {
		panic(http.StatusInternalServerError)
	}
	w.WriteHeader(http.StatusOK)

	remoteIP := getIP(r)
	ks.jchfChan <- []*kt.JCHF{ks.getJCHF(&wrapper, remoteIP)}
	ks.metrics.Messages.Mark(1)
}

func (ks *KentikHttpListener) getJCHF(wrapper *basic, remoteIP string) *kt.JCHF {
	in := kt.NewJCHF()
	in.CustomStr = make(map[string]string)
	in.CustomInt = make(map[string]int32)
	in.CustomBigInt = make(map[string]int64)
	in.EventType = strings.ReplaceAll(wrapper.Name, ".", "_")
	in.Provider = kt.ProviderHttpDevice
	in.SrcAddr = remoteIP

	// Use host for device_name if its set.
	if host, ok := wrapper.Tags["host"]; ok {
		remoteIP = host
	}

	if dev, ok := ks.devices[remoteIP]; ok {
		in.DeviceName = dev.Name // Copy in any of these info we get
		in.DeviceId = dev.ID
		in.CompanyId = dev.CompanyID
		in.SampleRate = dev.SampleRate
		dev.SetUserTags(in.CustomStr)
	}

	in.Timestamp = wrapper.Timestamp
	for t, v := range wrapper.Tags {
		in.CustomStr[t] = v
	}
	for f, v := range wrapper.Fields {
		in.CustomBigInt[f] = int64(v)
	}

	return in
}

// Get the JCHF content directly.
func (ks *KentikHttpListener) readJCHF(w http.ResponseWriter, r *http.Request) {
	wrapper := []*kt.JCHF{}

	// Decode body in gzip format if the request header is set this way.
	body := r.Body
	if r.Header.Get("Content-Encoding") == "gzip" {
		z, err := gzip.NewReader(r.Body)
		if err != nil {
			panic(http.StatusInternalServerError)
		}
		body = z
	}
	defer body.Close()

	if err := json.NewDecoder(body).Decode(&wrapper); err != nil {
		panic(http.StatusInternalServerError)
	}
	w.WriteHeader(http.StatusOK)

	ks.metrics.Messages.Mark(int64(len(wrapper)))
	for _, chf := range wrapper {
		chf.SetMap()
	}

	ks.jchfChan <- wrapper
}

func (ks *KentikHttpListener) Close() {}

func (ks *KentikHttpListener) HttpInfo() map[string]float64 {
	msgs := map[string]float64{
		"messages": ks.metrics.Messages.Rate1(),
		"errors":   ks.metrics.Errors.Rate1(),
	}
	return msgs
}

func (ks *KentikHttpListener) wrap(f handler) handler {
	return func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if r := recover(); r != nil {
				ks.metrics.Errors.Mark(1)
				if code, ok := r.(int); ok {
					http.Error(w, http.StatusText(code), code)
					return
				}
				panic(r)
			}
		}()

		if err := r.ParseForm(); err != nil {
			panic(http.StatusBadRequest)
		}

		f(w, r)
	}
}

type handler func(http.ResponseWriter, *http.Request)

func getIP(r *http.Request) string {
	res := r.Header.Get("X-FORWARDED-FOR")
	if res == "" {
		res = r.RemoteAddr
	}
	pts := strings.SplitN(res, ":", 2)

	return pts[0]
}

func (ks *KentikHttpListener) run(ctx context.Context) {
	deviceTicker := time.NewTicker(DeviceUpdateDuration)
	defer deviceTicker.Stop()

	ks.Infof("kentik http running, registered at %s", Listen)
	for {
		select {
		case <-deviceTicker.C:
			go func() {
				ks.Infof("Updating the network flow device list.")
				ks.devices = ks.apic.GetDevicesAsMap(0)
			}()
		case <-ctx.Done():
			ks.Infof("kentik http done")
			return
		}
	}
}
