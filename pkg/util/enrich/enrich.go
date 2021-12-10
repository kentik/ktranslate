package enrich

import (
	"bytes"
	"context"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/kentik/ktranslate/pkg/eggs/logger"
	"github.com/kentik/ktranslate/pkg/kt"
)

type Enricher struct {
	logger.ContextL
	url    string
	client *http.Client
}

func NewEnricher(url string, log logger.Underlying) (*Enricher, error) {
	e := Enricher{
		ContextL: logger.NewContextLFromUnderlying(logger.SContext{S: "Enricher"}, log),
		url:      url,
		client:   &http.Client{},
	}

	e.Infof("Enriching at %s", url)
	return &e, nil
}

func (e *Enricher) Enrich(ctx context.Context, msgs []*kt.JCHF) ([]*kt.JCHF, error) {

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
