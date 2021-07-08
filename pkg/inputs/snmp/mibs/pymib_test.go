package mibs

// TODO -- do we still need/want this?
/**
import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/kentik/ktranslate/pkg/eggs/logger"
	lt "github.com/kentik/ktranslate/pkg/eggs/logger/testing"
)

func TestFullMyMib(t *testing.T) {
	l := lt.NewTestContextL(logger.NilContext, t)
	mdb, err := NewMibDB("", "", l)
	assert.NoError(t, err)
	defer mdb.Close()

	assert.Equal(t, 2, len(mdb.profiles))
	pro := mdb.FindProfile(".1.3.6.1.4.1.89.108")
	assert.NotNil(t, pro)

	dev, ifm := pro.GetMetrics([]string{"IF-MIB"})
	assert.Equal(t, 0, len(dev))
	assert.Equal(t, 26, len(ifm))

	devm, ifmm := pro.GetMetadata([]string{"IF-MIB"})
	assert.Equal(t, 0, len(devm))
	assert.Equal(t, 14, len(ifmm))

}
*/
