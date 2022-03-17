package util

import (
	"bytes"
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLoadFile(t *testing.T) {
	assert := assert.New(t)
	content := []byte("aaaa") // Set some test content.

	// Some test server
	svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, string(content))
	}))
	defer svr.Close()

	sMock = &mockS3{
		lastContent: content,
	}

	tests := map[string][]byte{
		":foo":                nil,
		svr.URL:               content,
		"s3://foo/bar/a/one":  nil,
		"S3://foo/bar/a/two":  nil,
		"s3m://foo/bar/a/one": content,
	}

	// Save test data to local.
	file, err := ioutil.TempFile("", "")
	if err != nil {
		t.FailNow()
	}
	if _, err := file.Write(content); err != nil {
		t.FailNow()
	}
	defer os.Remove(file.Name())
	tests[file.Name()] = content

	ctx := context.Background()
	for in, expt := range tests {
		res, err := LoadFile(ctx, in)
		if expt == nil {
			assert.Error(err)
		} else {
			assert.NoError(err)
			assert.Equal(string(expt), string(res), "failed %s", in)
		}
	}
}

func TestWriteFile(t *testing.T) {
	assert := assert.New(t)
	content := []byte("aaaa") // Set some test content.

	fileWeb, err := ioutil.TempFile("", "")
	if err != nil {
		t.FailNow()
	}
	defer os.Remove(fileWeb.Name())

	// Some test server
	svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body := r.Body
		defer body.Close()
		bb := []byte{}
		bodyBuffer := bytes.NewBuffer(bb)
		_, err = bodyBuffer.ReadFrom(body)
		ioutil.WriteFile(fileWeb.Name(), bodyBuffer.Bytes(), 0644)
	}))
	defer svr.Close()

	// Save test data to local.
	file, err := ioutil.TempFile("", "")
	if err != nil {
		t.FailNow()
	}

	sMock = &mockS3{
		lastContent: nil,
	}

	defer os.Remove(file.Name())
	tests := map[string]string{
		":foo":                "",
		file.Name():           file.Name(),
		svr.URL:               fileWeb.Name(),
		"s3m://foo/bar/a/one": string(content),
	}

	ctx := context.Background()
	for target, local := range tests {
		err := WriteFile(ctx, target, content, 0644)
		if local == "" {
			assert.Error(err)
		} else {
			assert.NoError(err)
			c, err := ioutil.ReadFile(local) // This one assumes that we're writting locally.
			if err != nil {
				c = sMock.lastContent
			}
			assert.Equal(string(content), string(c), "failed %s", target)
		}
	}
}
