package cmetrics

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	metrics "github.com/kentik/go-metrics"
)

const (
	influxMaxSendTries          = 2
	influxSendSleep             = 1 * time.Second
	influxClientResponseTimeout = 5 * time.Second
	influxClientKeepAlive       = 60 * time.Second
	influxClientTLSTimeout      = 5 * time.Second
	influxContentType           = "application/x-www-form-urlencoded"
)

type INFLUXMetricSet struct {
	metrics []*INFLUXMetric
}

type INFLUXMetric struct {
	Metric    string
	Timestamp int64
	Value     int64
	Type      string
	Tags      map[string]string
}

// OpenINFLUXConfig provides a container with configuration parameters for
// the OpenINFLUX exporter
type OpenINFLUXConfig struct {
	Addr               string                // Network address to connect to
	Registry           metrics.Registry      // Registry to be exported
	FlushInterval      time.Duration         // Flush interval
	DurationUnit       time.Duration         // Time conversion unit for durations
	Prefix             string                // Prefix to be prepended to metric names
	Debug              bool                  // write to stdout for debug
	Quiet              bool                  // silence all comms
	Tags               map[string]string     // add these tags to each metric writen
	Send               chan *INFLUXMetricSet // manage # of outstanding http requests here.
	MaxHttpOutstanding int
	ProxyUrl           string
	Extra              map[string]string
}

// OpenINFLUX is a blocking exporter function which reports metrics in r
// to a INFLUX server located at addr, flushing them every d duration
// and prepending metric names with prefix.
func OpenINFLUX(r metrics.Registry, d time.Duration, prefix string, addr string, maxOutstanding int) {
	OpenINFLUXWithConfig(OpenINFLUXConfig{
		Addr:               addr,
		Registry:           r,
		FlushInterval:      d,
		DurationUnit:       time.Nanosecond,
		Prefix:             prefix,
		Debug:              false,
		MaxHttpOutstanding: maxOutstanding,
		Send:               make(chan *INFLUXMetricSet, maxOutstanding),
		Tags:               map[string]string{},
		Extra:              nil,
	})
}

// OpenINFLUXWithConfig is a blocking exporter function just like OpenINFLUX,
// but it takes a OpenINFLUXConfig instead.
func OpenINFLUXWithConfig(c OpenINFLUXConfig) {
	go c.runSend()

	for _ = range time.Tick(c.FlushInterval) {
		if err := openINFLUX(&c); nil != err {
			log.Println(err)
		}
	}
}

func (c *OpenINFLUXConfig) runSend() {
	if strings.HasPrefix(c.Addr, "http") {
		c.runSendViaHTTP()
	} else if strings.HasPrefix(c.Addr, "tcp") || strings.HasPrefix(c.Addr, "udp") {
		c.runSendViaSocket()
	}
}

func (c *OpenINFLUXConfig) runSendViaSocket() {
	for r := range c.Send {
		var w *bufio.Writer
		var conn net.Conn = nil

		if c.Debug {
			w = bufio.NewWriter(os.Stdout)
		} else {
			var err error
			pts := strings.Split(c.Addr, "://")
			if len(pts) == 2 {
				conn, err = net.Dial(pts[0], pts[1])
				if nil != err {
					if !c.Quiet {
						fmt.Printf("Invalid metrics address: %s, %v\n", c.Addr, err)
					}
					continue
				}
				w = bufio.NewWriter(conn)
			} else {
				if !c.Quiet {
					fmt.Printf("Invalid metrics address: %s\n", c.Addr)
				}
				continue
			}
		}

		if ebytes, err := r.ToWire(); err != nil {
			if !c.Quiet {
				fmt.Printf("Error encoding to wire format: %v\n", err)
			}
			continue
		} else {
			w.Write(ebytes)
			w.Flush()
		}

		if !c.Debug {
			conn.Close()
		}
	}
}

func (c *OpenINFLUXConfig) runSendViaHTTP() {
	tr := &http.Transport{
		DisableCompression: false,
		DisableKeepAlives:  false,
		Dial: (&net.Dialer{
			Timeout:   influxClientResponseTimeout,
			KeepAlive: influxClientKeepAlive,
		}).Dial,
		TLSHandshakeTimeout: influxClientTLSTimeout,
	}

	// Add a proxy if needed.
	if c.ProxyUrl != "" {
		proxyUrl, err := url.Parse(c.ProxyUrl)
		if err != nil {
			if !c.Quiet {
				fmt.Printf("Error setting proxy: %v\n", err)
			}
		} else {
			tr.Proxy = http.ProxyURL(proxyUrl)
			if !c.Quiet {
				fmt.Printf("Set outbound proxy: %s\n", c.ProxyUrl)
			}
		}
	}

	client := &http.Client{Transport: tr, Timeout: influxClientResponseTimeout}

	for r := range c.Send {
		if ebytes, err := r.ToWire(); err != nil {
			if !c.Quiet {
				fmt.Printf("Error encoding to wire format: %v\n", err)
			}
			continue
		} else {
			if c.Debug {
				fmt.Printf("Metrics: %v", string(ebytes))
			} else {
				for i := 0; i < influxMaxSendTries; i++ {
					req, err := http.NewRequest("POST", c.Addr, bytes.NewBuffer(ebytes))
					if err != nil {
						if !c.Quiet {
							fmt.Printf("Error Creating Request: %v\n", err)
						}
						continue
					}
					req.Header.Add("Content-Type", influxContentType)

					resp, err := client.Do(req)
					if err != nil {
						if !c.Quiet {
							if i == influxMaxSendTries-1 {
								fmt.Printf("Warn Posting to %s, giving up: %v\n", c.Addr, err)
							} else {
								fmt.Printf("Retry Posting to %s: %v\n", c.Addr, err)
							}
						}
						client = &http.Client{Transport: tr, Timeout: influxClientResponseTimeout}
					} else {
						// Fire and forget
						io.Copy(ioutil.Discard, resp.Body)
						resp.Body.Close()
						break
					}
				}
			}
		}
	}
}

/**
Write out additional tags
*/
func openINFLUX(c *OpenINFLUXConfig) error {

	shortHostnameBase := GetShortHostname()
	now := time.Now().UnixNano()
	sendBody := NewINFLUXMetricSet()
	du := float64(c.DurationUnit)

	c.Registry.Each(func(baseName string, i interface{}) {
		name, tags := ExpandTags(baseName, c.Prefix, shortHostnameBase, c.Tags, c.Extra)

		switch metric := i.(type) {
		case metrics.Counter:
			sendBody.Add(&INFLUXMetric{Metric: name, Timestamp: now, Value: metric.Count(), Tags: tags, Type: "count"})
		case metrics.Gauge:
			sendBody.Add(&INFLUXMetric{Metric: name, Timestamp: now, Value: metric.Value(), Tags: tags, Type: "value"})
		case metrics.Histogram:
			h := metric.Snapshot()
			ps := h.Percentiles([]float64{0.5, 0.75, 0.95, 0.99, 0.999})
			sendBody.Add(&INFLUXMetric{Metric: name, Timestamp: now, Value: h.Count(), Tags: tags, Type: "count"})
			sendBody.Add(&INFLUXMetric{Metric: name, Timestamp: now, Value: h.Min(), Tags: tags, Type: "min"})
			sendBody.Add(&INFLUXMetric{Metric: name, Timestamp: now, Value: h.Max(), Tags: tags, Type: "max"})
			sendBody.Add(&INFLUXMetric{Metric: name, Timestamp: now, Value: int64(h.Mean()), Tags: tags, Type: "mean"})
			sendBody.Add(&INFLUXMetric{Metric: name, Timestamp: now, Value: int64(ps[2]), Tags: tags, Type: "95-percentile"})
			sendBody.Add(&INFLUXMetric{Metric: name, Timestamp: now, Value: int64(ps[3]), Tags: tags, Type: "99-percentile"})
			metric.Clear()
		case metrics.Meter:
			m := metric.Snapshot()
			sendBody.Add(&INFLUXMetric{Metric: name, Timestamp: now, Value: m.Count(), Tags: tags, Type: "count"})
			sendBody.Add(&INFLUXMetric{Metric: name, Timestamp: now, Value: int64(m.Rate1()), Tags: tags, Type: "one-minute"})
		case metrics.Timer:
			t := metric.Snapshot()
			ps := t.Percentiles([]float64{0.5, 0.75, 0.95, 0.99, 0.999})
			sendBody.Add(&INFLUXMetric{Metric: name, Timestamp: now, Value: t.Count(), Tags: tags, Type: "count"})
			sendBody.Add(&INFLUXMetric{Metric: name, Timestamp: now, Value: t.Min() / int64(du), Tags: tags, Type: "min"})
			sendBody.Add(&INFLUXMetric{Metric: name, Timestamp: now, Value: t.Max() / int64(du), Tags: tags, Type: "max"})
			sendBody.Add(&INFLUXMetric{Metric: name, Timestamp: now, Value: int64(t.Mean() / du), Tags: tags, Type: "mean"})
			sendBody.Add(&INFLUXMetric{Metric: name, Timestamp: now, Value: int64(ps[2] / du), Tags: tags, Type: "95-percentile"})
			sendBody.Add(&INFLUXMetric{Metric: name, Timestamp: now, Value: int64(ps[3] / du), Tags: tags, Type: "99-percentile"})
			sendBody.Add(&INFLUXMetric{Metric: name, Timestamp: now, Value: int64(t.Rate1()), Tags: tags, Type: "one-minute"})
			sendBody.Add(&INFLUXMetric{Metric: name, Timestamp: now, Value: int64(t.Rate5()), Tags: tags, Type: "five-minute"})
			sendBody.Add(&INFLUXMetric{Metric: name, Timestamp: now, Value: int64(t.Rate15()), Tags: tags, Type: "fifteen-minute"})
			metric.Clear()
		}
	})

	if sendBody.Len() > 0 {
		if len(c.Send) < c.MaxHttpOutstanding {
			c.Send <- sendBody
		} else {
			if !c.Quiet {
				fmt.Printf("Dropping flow: Q at %d\n", len(c.Send))
			}
		}
	}

	return nil
}

func (m *INFLUXMetric) ToWire() []byte {
	tags := make([]string, (len(m.Tags))+2)
	tags[0] = m.Metric
	tags[1] = "type=" + m.Type
	i := 2
	for k, v := range m.Tags {
		tags[i] = k + "=" + v
		i++
	}

	return []byte(fmt.Sprintf("%s value=%d %d\n", strings.Join(tags, ","), m.Value, m.Timestamp))
}

func (m *INFLUXMetricSet) Len() int {
	return len(m.metrics)
}

func (m *INFLUXMetricSet) ToWire() ([]byte, error) {
	var buf bytes.Buffer
	for _, met := range m.metrics {
		buf.Write(met.ToWire())
	}
	return buf.Bytes(), nil
}

func (m *INFLUXMetricSet) Add(met *INFLUXMetric) {
	if met.Value > 0 {
		m.metrics = append(m.metrics, met)
	}
}

func NewINFLUXMetricSet() *INFLUXMetricSet {
	return &INFLUXMetricSet{
		metrics: make([]*INFLUXMetric, 0),
	}
}
