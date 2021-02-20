package features

import (
	"github.com/kentik/ktranslate/pkg/eggs/properties"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFeaturesWithNoBacking(t *testing.T) {
	props := properties.NewPropertyService()
	feats := NewFeatureService(props)
	assert.False(t, feats.EnabledGlobally("testing"))
}

func TestFeaturesWithStaticBacking(t *testing.T) {
	backing := properties.NewStaticMapPropertyBacking(map[string]string{
		"features.frobber-global":          "true",
		"features.frobber-cid_11":          "false",
		"features.frobber-cid_11-uid_9999": "true",
		"features.goobs-global":            "false",
		"features.goobs-cid_11":            "true",
	})
	props := properties.NewPropertyService(backing)
	feats := NewFeatureService(props)

	assert.True(t, feats.EnabledGlobally("frobber"))
	assert.True(t, feats.Enabled("frobber"))
	assert.True(t, feats.Enabled("frobber", "cid", 10, "uid", 222))
	assert.False(t, feats.Enabled("frobber", "cid", 11))
	assert.False(t, feats.Enabled("frobber", "cid", 11, "uid", 222))
	assert.True(t, feats.Enabled("frobber", "cid", 11, "uid", 9999))

	assert.False(t, feats.EnabledGlobally("goobs"))
	assert.False(t, feats.Enabled("goobs"))
	assert.False(t, feats.Enabled("goobs", "cid", 10))
	assert.True(t, feats.Enabled("goobs", "cid", 11, "uid", 0))
	assert.True(t, feats.Enabled("goobs", "cid", 11, "uid", 1234))
}
