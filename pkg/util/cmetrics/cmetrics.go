package cmetrics

import (
	"log"
	"log/syslog"
	"net"
	"os"
	"strings"
	"time"

	metrics "github.com/kentik/go-metrics"
)

var (
	SYSLOG_FILE_PATH    = "/dev/log"
	MAX_HTTP_REQ        = 3 // # in-flight metric calls
	CH_HTTP_LOCAL_PROXY = "CH_HTTP_LOCAL_PROXY"
)

// Logger abstracts away the logging implementation
type Logger interface {
	Debugf(prefix, format string, v ...interface{})
	Infof(prefix, format string, v ...interface{})
	Errorf(prefix, format string, v ...interface{})
	Warnf(prefix, format string, v ...interface{})
}

func SetConf(conf string, l Logger, log_prefix string, tsdb_prefix string, tags []string, extra []string, apiEmail *string, apiPassword *string) {
	SetConfWithRegistry(conf, l, log_prefix, tsdb_prefix, tags, extra, apiEmail, apiPassword, metrics.DefaultRegistry)
}

func SetConfWithRegistry(conf string, l Logger, log_prefix string, tsdb_prefix string, tags []string, extra []string, apiEmail *string, apiPassword *string, registry metrics.Registry) {
	l.Infof(log_prefix, "Setting metrics: %s", conf)

	if conf != "none" {
		switch conf {
		case "syslog":
			if w, err := syslog.New(syslog.LOG_INFO, "metrics"); err == nil && w != nil {
				go metrics.Syslog(registry, 60e9, w)
			} else {
				l.Errorf(log_prefix, "Could not start syslog metrics: %v", err)
			}
		case "stderr":
			go metrics.Log(registry, 60e9, log.New(os.Stderr, "metrics: ", log.Lmicroseconds))
		default:
			dest := strings.SplitN(conf, ":", 2)
			switch dest[0] {
			case "graphite":
				l.Infof(log_prefix, "Metrics: Connecting to graphite: %s", dest[1])
				addr, _ := net.ResolveTCPAddr("tcp", dest[1])
				go metrics.Graphite(registry, 10e9, "metrics", addr)
			case "tsdb", "tsdb_debug":
				flushTime := 60 * time.Second
				if dest[0] == "tsdb_debug" {
					flushTime = 30 * time.Second
				}

				if strings.HasPrefix(dest[1], "http") {
					l.Infof(log_prefix, "Metrics: Connecting to [%s]: %s. [HTTP]", dest[0], dest[1])
					go OpenHTTPTSDBWithConfig(OpenHTTPTSDBConfig{
						Addr:               dest[1],
						Registry:           registry,
						FlushInterval:      flushTime,
						DurationUnit:       time.Millisecond,
						Prefix:             tsdb_prefix,
						Debug:              (dest[0] == "tsdb_debug"),
						Tags:               TagsMap(tags),
						Send:               make(chan []byte, MAX_HTTP_REQ),
						ProxyUrl:           os.Getenv(CH_HTTP_LOCAL_PROXY),
						MaxHttpOutstanding: MAX_HTTP_REQ,
						Extra:              TagsMap(extra),
						ApiEmail:           apiEmail,
						ApiPassword:        apiPassword,
					})
				} else {
					l.Infof(log_prefix, "Metrics: Connecting to [%s]: %s. [TCP]. Debug=%v", dest[0], dest[1], (dest[0] == "tsdb_debug"))
					addr, err := net.ResolveTCPAddr("tcp", dest[1])
					if err != nil {
						l.Errorf(log_prefix, "Could not resolve address: %s %v", dest[1], err)
					} else {
						go OpenTSDBWithConfig(OpenTSDBConfig{
							Addr:          addr,
							Registry:      registry,
							FlushInterval: flushTime,
							DurationUnit:  time.Millisecond,
							Prefix:        tsdb_prefix,
							Debug:         (dest[0] == "tsdb_debug"),
							Tags:          TagsMap(tags),
							Extra:         TagsMap(extra),
						})
					}
				}
			case "influx", "influx_debug", "influx_quiet":
				flushTime := 60 * time.Second
				if dest[0] == "influx_debug" {
					flushTime = 30 * time.Second
				}

				if strings.HasPrefix(dest[1], "http") || strings.HasPrefix(dest[1], "tcp") || strings.HasPrefix(dest[1], "udp") {
					l.Infof(log_prefix, "Metrics: Connecting Influx to [%s]: %s. [HTTP|TCP|UDP]", dest[0], dest[1])
					go OpenINFLUXWithConfig(OpenINFLUXConfig{
						Addr:               dest[1],
						Registry:           registry,
						FlushInterval:      flushTime,
						DurationUnit:       time.Millisecond,
						Prefix:             tsdb_prefix,
						Debug:              (dest[0] == "influx_debug"),
						Quiet:              (dest[0] == "influx_quiet"),
						Tags:               TagsMap(tags),
						Send:               make(chan *INFLUXMetricSet, MAX_HTTP_REQ),
						ProxyUrl:           os.Getenv(CH_HTTP_LOCAL_PROXY),
						MaxHttpOutstanding: MAX_HTTP_REQ,
						Extra:              TagsMap(extra),
					})
				} else {
					l.Errorf(log_prefix, "Only HTTP|TCP|UDP endpoint for influx currently supported")
				}
			}
		}
	}
}
