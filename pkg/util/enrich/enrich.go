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
	url    string
	client *http.Client
	doSrc  bool
	doDst  bool
	salt   []byte
}

func NewEnricher(url string, log logger.Underlying) (*Enricher, error) {
	e := Enricher{
		ContextL: logger.NewContextLFromUnderlying(logger.SContext{S: "Enricher"}, log),
		url:      url,
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
	}

	e.Infof("Enriching at %s. Source: %v, Dest: %v, Salt %s", url, e.doSrc, e.doDst, string(e.salt))
	return &e, nil
}

func (e *Enricher) Enrich(ctx context.Context, msgs []*kt.JCHF) ([]*kt.JCHF, error) {
	if e.doSrc || e.doDst {
		return e.hashIP(ctx, msgs)
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
