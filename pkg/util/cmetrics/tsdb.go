package cmetrics

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
	"strings"
	"time"

	metrics "github.com/kentik/go-metrics"
)

// OpenTSDBConfig provides a container with configuration parameters for
// the OpenTSDB exporter
type OpenTSDBConfig struct {
	Addr          *net.TCPAddr      // Network address to connect to
	Registry      metrics.Registry  // Registry to be exported
	FlushInterval time.Duration     // Flush interval
	DurationUnit  time.Duration     // Time conversion unit for durations
	Prefix        string            // Prefix to be prepended to metric names
	Debug         bool              // write to stdout for debug
	Tags          map[string]string // add these tags to each metric writen
	Extra         map[string]string // Extra tags added to tag list.
}

// OpenTSDB is a blocking exporter function which reports metrics in r
// to a TSDB server located at addr, flushing them every d duration
// and prepending metric names with prefix.
func OpenTSDB(r metrics.Registry, d time.Duration, prefix string, addr *net.TCPAddr) {
	OpenTSDBWithConfig(OpenTSDBConfig{
		Addr:          addr,
		Registry:      r,
		FlushInterval: d,
		DurationUnit:  time.Nanosecond,
		Prefix:        prefix,
		Debug:         false,
		Tags:          nil,
		Extra:         nil,
	})
}

// OpenTSDBWithConfig is a blocking exporter function just like OpenTSDB,
// but it takes a OpenTSDBConfig instead.
func OpenTSDBWithConfig(c OpenTSDBConfig) {
	for _ = range time.Tick(c.FlushInterval) {
		if err := openTSDB(&c); nil != err {
			log.Println(err)
		}
	}
}

/**
Write out additional tags
*/
func openTSDB(c *OpenTSDBConfig) error {

	var w *bufio.Writer
	shortHostnameBase := GetShortHostname()
	now := time.Now().Unix()
	du := float64(c.DurationUnit)
	var conn *net.TCPConn = nil

	if c.Debug {
		w = bufio.NewWriter(os.Stdout)
	} else {
		var err error
		conn, err = net.DialTCP("tcp", nil, c.Addr)
		if nil != err {
			return err
		}
		w = bufio.NewWriter(conn)
	}

	defer func() {
		if !c.Debug {
			conn.Close()
		}
	}()

	c.Registry.Each(func(baseName string, i interface{}) {
		name, tagsMap := ExpandTags(baseName, c.Prefix, shortHostnameBase, c.Tags, c.Extra)
		tags := make([]string, 0, len(tagsMap))
		for k, v := range tagsMap {
			tags = append(tags, fmt.Sprintf("%s=%s", k, v))
		}
		tagsSer := strings.Join(tags, " ")

		switch metric := i.(type) {
		case metrics.Counter:
			fmt.Fprintf(w, "put %s.%s %d %d %s type=count\n", c.Prefix, name, now, metric.Count(), tagsSer)
		case metrics.Gauge:
			fmt.Fprintf(w, "put %s.%s %d %d %s type=value\n", c.Prefix, name, now, metric.Value(), tagsSer)
		case metrics.Histogram:
			h := metric.Snapshot()
			ps := h.Percentiles([]float64{0.5, 0.75, 0.95, 0.99, 0.999})
			fmt.Fprintf(w, "put %s.%s %d %d %s type=count\n", c.Prefix, name, now, h.Count(), tagsSer)
			fmt.Fprintf(w, "put %s.%s %d %d %s type=min\n", c.Prefix, name, now, h.Min(), tagsSer)
			fmt.Fprintf(w, "put %s.%s %d %d %s type=max\n", c.Prefix, name, now, h.Max(), tagsSer)
			fmt.Fprintf(w, "put %s.%s %d %.2f %s type=mean\n", c.Prefix, name, now, h.Mean(), tagsSer)
			//fmt.Fprintf(w, "put %s.%s %d %.2f %s type=std-dev\n", c.Prefix, name, now, h.StdDev(), tagsSer)
			//fmt.Fprintf(w, "put %s.%s %d %.2f %s type=50-percentile\n", c.Prefix, name, now, ps[0], tagsSer)
			//fmt.Fprintf(w, "put %s.%s %d %.2f %s type=75-percentile\n", c.Prefix, name, now, ps[1], tagsSer)
			fmt.Fprintf(w, "put %s.%s %d %.2f %s type=95-percentile\n", c.Prefix, name, now, ps[2], tagsSer)
			fmt.Fprintf(w, "put %s.%s %d %.2f %s type=99-percentile\n", c.Prefix, name, now, ps[3], tagsSer)
			//fmt.Fprintf(w, "put %s.%s %d %.2f %s type=999-percentile\n", c.Prefix, name, now, ps[4], tagsSer)
			metric.Clear()
		case metrics.Meter:
			m := metric.Snapshot()
			fmt.Fprintf(w, "put %s.%s %d %d %s type=count\n", c.Prefix, name, now, m.Count(), tagsSer)
			fmt.Fprintf(w, "put %s.%s %d %.2f %s type=one-minute\n", c.Prefix, name, now, m.Rate1(), tagsSer)
			//fmt.Fprintf(w, "put %s.%s %d %.2f %s type=five-minute\n", c.Prefix, name, now, m.Rate5(), tagsSer)
			//fmt.Fprintf(w, "put %s.%s %d %.2f %s type=fifteen-minute\n", c.Prefix, name, now, m.Rate15(), tagsSer)
			//fmt.Fprintf(w, "put %s.%s %d %.2f %s type=mean\n", c.Prefix, name, now, m.RateMean(), tagsSer)
		case metrics.Timer:
			t := metric.Snapshot()
			ps := t.Percentiles([]float64{0.5, 0.75, 0.95, 0.99, 0.999})
			fmt.Fprintf(w, "put %s.%s %d %d %s type=count\n", c.Prefix, name, now, t.Count(), tagsSer)
			fmt.Fprintf(w, "put %s.%s %d %d %s type=min\n", c.Prefix, name, now, t.Min()/int64(du), tagsSer)
			fmt.Fprintf(w, "put %s.%s %d %d %s type=max\n", c.Prefix, name, now, t.Max()/int64(du), tagsSer)
			fmt.Fprintf(w, "put %s.%s %d %.2f %s type=mean\n", c.Prefix, name, now, t.Mean()/du, tagsSer)
			//fmt.Fprintf(w, "put %s.%s %d %.2f %s type=std-dev\n", c.Prefix, name, now, t.StdDev()/du, tagsSer)
			//fmt.Fprintf(w, "put %s.%s %d %.2f %s type=50-percentile\n", c.Prefix, name, now, ps[0]/du, tagsSer)
			//fmt.Fprintf(w, "put %s.%s %d %.2f %s type=75-percentile\n", c.Prefix, name, now, ps[1]/du, tagsSer)
			fmt.Fprintf(w, "put %s.%s %d %.2f %s type=95-percentile\n", c.Prefix, name, now, ps[2]/du, tagsSer)
			fmt.Fprintf(w, "put %s.%s %d %.2f %s type=99-percentile\n", c.Prefix, name, now, ps[3]/du, tagsSer)
			//fmt.Fprintf(w, "put %s.%s %d %.2f %s type=999-percentile\n", c.Prefix, name, now, ps[4]/du, tagsSer)
			fmt.Fprintf(w, "put %s.%s %d %.2f %s type=one-minute\n", c.Prefix, name, now, t.Rate1(), tagsSer)
			fmt.Fprintf(w, "put %s.%s %d %.2f %s type=five-minute\n", c.Prefix, name, now, t.Rate5(), tagsSer)
			fmt.Fprintf(w, "put %s.%s %d %.2f %s type=fifteen-minute\n", c.Prefix, name, now, t.Rate15(), tagsSer)
			//fmt.Fprintf(w, "put %s.%s %d %.2f %s type=mean-rate\n", c.Prefix, name, now, t.RateMean(), tagsSer)
			metric.Clear()
		}
		w.Flush()
	})
	return nil
}
