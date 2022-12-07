package enrich

import (
	"bytes"
	"context"
	"crypto/sha256"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"strconv"
	"strings"

	"go.starlark.net/lib/math"
	"go.starlark.net/lib/time"
	"go.starlark.net/resolve"
	"go.starlark.net/starlark"
	"go.starlark.net/starlarkjson"

	"github.com/kentik/ktranslate/pkg/eggs/logger"
	"github.com/kentik/ktranslate/pkg/kt"
)

const (
	EnrichUrlHashSrcIP = "hash_src_ip"
	EnrichUrlHashDstIP = "hash_dst_ip"
	EnrichUrlHashAllIP = "hash_ip"
)

type Enricher struct {
	logger.ContextL
	url     string
	source  string
	script  string
	client  *http.Client
	doSrc   bool
	doDst   bool
	salt    []byte
	thread  *starlark.Thread
	globals starlark.StringDict
}

func NewEnricher(url string, source string, script string, log logger.Underlying) (*Enricher, error) {
	e := Enricher{
		ContextL: logger.NewContextLFromUnderlying(logger.SContext{S: "Enricher"}, log),
		url:      url,
		source:   source,
		script:   script,
		client:   &http.Client{},
		doSrc:    strings.HasPrefix(url, EnrichUrlHashSrcIP) || strings.HasPrefix(url, EnrichUrlHashAllIP),
		doDst:    strings.HasPrefix(url, EnrichUrlHashDstIP) || strings.HasPrefix(url, EnrichUrlHashAllIP),
	}

	if e.doSrc || e.doDst {
		var salt string
		if strings.HasPrefix(url, EnrichUrlHashAllIP) {
			salt = url[len(EnrichUrlHashAllIP):]
		} else {
			salt = url[len(EnrichUrlHashSrcIP):] // same # chars src and dst.
		}
		e.salt = []byte(salt)
	} else {
		if url != "" {
			e.Infof("Enriching at remote url %s", url)
		} else {
			if source == "" && script == "" {
				return nil, fmt.Errorf("Source, Script or URL required for enrichment.")
			}

			// Try loading as an a local file.
			thread := &starlark.Thread{
				Print: func(_ *starlark.Thread, msg string) { e.Infof("%s", msg) },
				Name:  "kentik enrich",
				Load:  e.LoadFunc,
			}
			builtins := starlark.StringDict{
				"catch":           starlark.NewBuiltin("catch", catch),
				"findAllSubmatch": starlark.NewBuiltin("findAllSubmatch", findAllSubmatch),
			}

			program, err := e.sourceProgram(builtins)
			if err != nil {
				return nil, err
			}

			// Execute source
			globals, err := program.Init(thread, builtins)
			if err != nil {
				return nil, err
			}

			// Place to store state across runs.
			globals["state"] = starlark.NewDict(0)
			globals.Freeze()

			e.thread = thread
			e.globals = globals
			if script != "" {
				e.Infof("Enriching via a starlark script at %s", script)
			} else {
				e.Infof("Enriching via a starlark program")
			}
		}
	}

	e.Infof("Enriching at %s. Source: %v, Dest: %v, Salt %s", url, e.doSrc, e.doDst, string(e.salt))
	return &e, nil
}

func (e *Enricher) EnrichMib(idx string, key string, attr map[string]interface{}, lm *kt.LastMetadata) {
	mm := &Mib{ContextL: e}
	mm.Wrap(idx, key, attr, lm)
	_, err := starlark.Call(e.thread, e.globals["main"], starlark.Tuple{mm}, nil)
	if err != nil {
		e.Errorf("Cannot run enrich mib script: %v", err)
		return
	}
}

func (e *Enricher) EnrichMetric(idx string, key string, ints map[string]int64, strs map[string]string, metrics map[string]kt.MetricInfo) {
	mm := &MibMetric{ContextL: e}
	mm.Wrap(idx, key, ints, strs, metrics)
	_, err := starlark.Call(e.thread, e.globals["main"], starlark.Tuple{mm}, nil)
	if err != nil {
		e.Errorf("Cannot run enrich mib script: %v", err)
		return
	}
}

func (e *Enricher) Enrich(ctx context.Context, msgs []*kt.JCHF) ([]*kt.JCHF, error) {
	if e.doSrc || e.doDst {
		return e.hashIP(ctx, msgs)
	} else if e.globals != nil {
		return e.runScript(ctx, msgs)
	}

	target, err := json.Marshal(msgs) // Has to be an array here, no idea why.
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, "POST", e.url, bytes.NewBuffer(target))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	resp, err := e.client.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		e.client = &http.Client{}
		return nil, err
	} else if resp.StatusCode >= 300 {
		err = fmt.Errorf("HTTP error: status code %d, body %v", resp.StatusCode, string(body))
		return nil, err
	}

	err = json.Unmarshal(body, &msgs)
	return msgs, err
}

func (e *Enricher) hashIP(ctx context.Context, msgs []*kt.JCHF) ([]*kt.JCHF, error) {
	h := sha256.New()
	for _, msg := range msgs {
		if e.doSrc {
			h.Write(e.salt)
			h.Write([]byte(msg.SrcAddr))
			msg.SrcAddr = net.IP(h.Sum(nil)[0:16]).String()
			msg.CustomStr["src_endpoint"] = msg.SrcAddr + ":" + strconv.Itoa(int(msg.L4SrcPort))
			h.Reset()
		}
		if e.doDst {
			h.Write(e.salt)
			h.Write([]byte(msg.DstAddr))
			msg.DstAddr = net.IP(h.Sum(nil)[0:16]).String()
			msg.CustomStr["dst_endpoint"] = msg.DstAddr + ":" + strconv.Itoa(int(msg.L4DstPort))
			h.Reset()
		}
	}

	return msgs, nil
}

func (e *Enricher) runScript(ctx context.Context, msgs []*kt.JCHF) ([]*kt.JCHF, error) {
	inputs := []starlark.Value{}
	for _, msg := range msgs {
		lm := msg
		jf := &JCHF{}
		jf.Wrap(lm)
		inputs = append(inputs, jf)
	}
	rv, err := starlark.Call(e.thread, e.globals["main"], starlark.Tuple{starlark.NewList(inputs)}, nil)
	if err != nil {
		return nil, err
	}
	switch rv := rv.(type) {
	case starlark.NoneType:
		return nil, nil
	case starlark.Int:
		e.Infof("RC %d", rv)
	}

	return msgs, nil
}

func (e *Enricher) sourceProgram(builtins starlark.StringDict) (*starlark.Program, error) {
	var src interface{}
	if e.source != "" {
		src = e.source
	}
	_, program, err := starlark.SourceProgram(e.script, src, builtins.Has)
	return program, err
}

func (e *Enricher) LoadFunc(thread *starlark.Thread, module string) (starlark.StringDict, error) {
	switch module {
	case "json.star":
		return starlark.StringDict{
			"json": starlarkjson.Module,
		}, nil
	case "math.star":
		return starlark.StringDict{
			"math": math.Module,
		}, nil
	case "time.star":
		return starlark.StringDict{
			"time": time.Module,
		}, nil
	default:
		return nil, fmt.Errorf("module %s is not available", module)
	}
}

func init() {
	// https://github.com/bazelbuild/starlark/issues/20
	resolve.AllowNestedDef = true
	resolve.AllowLambda = true
	resolve.AllowFloat = true
	resolve.AllowSet = true
	resolve.AllowGlobalReassign = true
	resolve.AllowRecursion = true
}
