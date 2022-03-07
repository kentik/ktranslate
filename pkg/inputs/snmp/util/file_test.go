package util

import (
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

	tests := map[string][]byte{
		":foo":               nil,
		svr.URL:              content,
		"s3://foo/bar/a/one": nil,
		"S3://foo/bar/a/two": nil,
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

	// Save test data to local.
	file, err := ioutil.TempFile("", "")
	if err != nil {
		t.FailNow()
	}

	defer os.Remove(file.Name())
	tests := map[string]bool{
		":foo":      false,
		file.Name(): true,
	}

	ctx := context.Background()
	for target, good := range tests {
		err := WriteFile(ctx, target, content, 0644)
		if !good {
			assert.Error(err)
		} else {
			assert.NoError(err)
			c, _ := ioutil.ReadFile(target) // This one assumes that we're writting locally.
			assert.Equal(string(content), string(c), "failed %s", target)
		}
	}
}
