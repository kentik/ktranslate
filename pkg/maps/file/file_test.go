package file

import (
	"bytes"
	"compress/gzip"
	"io/ioutil"
	"path/filepath"
	"testing"

	"github.com/kentik/ktranslate/pkg/eggs/logger"
	lt "github.com/kentik/ktranslate/pkg/eggs/logger/testing"
	"github.com/kentik/ktranslate/pkg/kt"

	"github.com/stretchr/testify/assert"
)

var testTags = []byte(`
100324,C_MARKET_SRC,1344420230,BHN - Bakersfield
100323,C_MARKET_DST,1353464128,LCHTR - Michigan
100324,C_MARKET_SRC,1344420636,LCHTR - Los Angeles
100323,C_MARKET_DST,1344420636,LCHTR - Los Angeles
100323,C_MARKET_DST,1353464487,LCHTR - SLO
100323,C_MARKET_DST,1353464485,LCHTR - Louisiana
100323,C_MARKET_DST,1353465119,LCHTR - Nevada
22,C_DST_CUSTOMERNETWORK,44,FOO
33,C_SRC_CUSTOMERNETWORK,44,BAR
101199,DST_SUBSCRIBER_ID,to_hex
0,
1,-
2,South Brisbane
3,Fuzhou
4,Melbourne
5,Guangzhou
6,Tokyo
`)

func TestGenMap(t *testing.T) {
	assert := assert.New(t)
	l := lt.NewTestContextL(logger.NilContext, t)

	dir := t.TempDir()
	tmpfn := filepath.Join(dir, "test_tags")
	if err := ioutil.WriteFile(tmpfn, testTags, 0666); err != nil {
		t.Fatal(err)
	}

	_, err := NewFileTagMapper(l.GetLogger().GetUnderlyingLogger(), "/bad/path")
	assert.Error(err)

	f, err := NewFileTagMapper(l.GetLogger().GetUnderlyingLogger(), tmpfn)
	assert.NoError(err)

	_, _, ok := f.LookupTagValue(kt.Cid(10), 0, "")
	assert.False(ok)

	k, v, ok := f.LookupTagValue(kt.Cid(10), 1344420636, "100323")
	assert.True(ok)
	assert.Equal(k, "c_market_dst")
	assert.Equal(v, "lchtr_-_los_angeles")

	k, v, ok = f.LookupTagValue(kt.Cid(10), 1344420230, "100324")
	assert.True(ok)
	assert.Equal(k, "c_market_src")
	assert.Equal(v, "bhn_-_bakersfield")

	k, v, ok = f.LookupTagValueBig(kt.Cid(10), int64(242124693101048), "101199")
	assert.True(ok)
	assert.Equal(k, "dst_subscriber_id")
	assert.Equal(v, "dc360c52d9f8")

	k, v, ok = f.LookupTagValue(kt.Cid(10), 44, "22")
	assert.True(ok)
	assert.Equal(k, "c_dst_customernetwork")
	assert.Equal(v, "foo")

	k, v, ok = f.LookupTagValue(kt.Cid(10), 44, "33")
	assert.True(ok)
	assert.Equal(k, "c_src_customernetwork")
	assert.Equal(v, "bar")

	v = f.LookupKV(6)
	assert.Equal(v, "Tokyo")
}

func TestGenMapWithGzip(t *testing.T) {
	assert := assert.New(t)
	l := lt.NewTestContextL(logger.NilContext, t)

	var buf bytes.Buffer
	zw := gzip.NewWriter(&buf)
	_, err := zw.Write(testTags)
	assert.NoError(err)
	zw.Close()

	dir := t.TempDir()
	tmpfn := filepath.Join(dir, "test_tags.gz")
	if err := ioutil.WriteFile(tmpfn, buf.Bytes(), 0666); err != nil {
		t.Fatal(err)
	}

	f, err := NewFileTagMapper(l.GetLogger().GetUnderlyingLogger(), tmpfn)
	if err != nil {
		t.Fatal(err)
	}

	_, _, ok := f.LookupTagValue(kt.Cid(10), 0, "")
	assert.False(ok)

	k, v, ok := f.LookupTagValue(kt.Cid(10), 1344420636, "100323")
	assert.True(ok)
	assert.Equal(k, "c_market_dst")
	assert.Equal(v, "lchtr_-_los_angeles")
}
