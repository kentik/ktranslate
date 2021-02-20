package patricia

// TODO: fix outdated tests?
import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net"
	"os"
	"testing"

	"github.com/kentik/patricia"
	"github.com/stretchr/testify/assert"
)

const (
	FORCE_USE   = false
	IPV4ADDRLEN = 32
	IPV6ADDRLEN = 128
)

type ExpectedGeoData struct {
	IP      string
	Country uint32
	Region  uint32
	City    uint32
}

func TestPrefix(t *testing.T) {
	tests := []string{
		`10.10.10.10\n`,
		`10.10.10.10x`,
		``,
		`woofoo`,
		`woofoo/10`,
		`woofoo/woo`,
		`1011.1011.1011.1011/24`,
		`.../24`,
		`/`,
	}

	for _, s := range tests {
		v4, v6, err := patricia.ParseIPFromString(s)
		assert.Error(t, err, s)
		assert.Nil(t, v4)
		assert.Nil(t, v6)
	}
}

// @TOOD -- make work.
// Needs spanning geo file to work I think.
func TestGeo(t *testing.T) {
	search := []string{"1.0.199.1"}
	searchRange := []string{"1.0.199.0/0"}
	searchNum := []uint32{19282, 21843, 21843, 21843, 21843, 21843}
	searchNumCity := []uint32{117, 7201, 4294, 7436, 4294, 4304}
	searchNumRegion := []uint32{61, 460, 213, 495, 213, 213}

	var content bytes.Buffer
	for i, ip := range searchRange {
		content.WriteString(fmt.Sprintf("%s,%d,%d,%d\n", ip, searchNum[i], searchNumCity[i], searchNumRegion[i]))
	}

	tmpfile, err := ioutil.TempFile("", "testFile")
	assert.NoError(t, err)

	defer os.Remove(tmpfile.Name()) // clean up

	tmpfile.Write(content.Bytes())
	assert.NoError(t, err)

	err = tmpfile.Close()
	assert.NoError(t, err)

	geo, err := OpenGeo(tmpfile.Name(), true, nil)
	assert.NoError(t, err)

	defer geo.Close()

	for i, s := range search {
		v4Addr, v6Addr, err := patricia.ParseIPFromString(s)
		assert.NoError(t, err)
		var node *NodeGeo
		if v4Addr != nil {
			node, err = geo.SearchBestFromHostGeo(net.IPv4(byte(v4Addr.Address>>24), byte(v4Addr.Address>>16), byte(v4Addr.Address>>8), byte(v4Addr.Address)))
			assert.Error(t, err, s)
		} else if v6Addr != nil {
			node, err = geo.SearchBestFromHostGeo(net.ParseIP(s))
			assert.Error(t, err, s)
		}

		assert.Nil(t, node)
		if node != nil {
			cntry := GetCountry(node)
			region := GetRegion(node)
			city := GetCity(node)

			if cntry != searchNum[i] {
				t.Errorf("For %s, Country Got %d, Expecting %d", s, cntry, searchNum[i])
			}

			if region != searchNumRegion[i] {
				t.Errorf("For %s, Region Got %d, Expecting %d", s, region, searchNumRegion[i])
			}

			if city != searchNumCity[i] {
				t.Errorf("For %s, City Got %d, Expecting %d", s, city, searchNumCity[i])
			}
		}
	}
}

func makeIP(a, b, c, d uint32) uint32 {
	return a<<24 + b<<16 + c<<8 + d
}

func mustCIDR(cidr string) *net.IPNet {
	_, ipnet, err := net.ParseCIDR(cidr)
	if err != nil {
		panic(err)
	}
	return ipnet
}
