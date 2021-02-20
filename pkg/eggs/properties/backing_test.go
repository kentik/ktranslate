package properties

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStaticPropertyBacking(t *testing.T) {
	propertyMap := map[string]string{
		"string.property": "cool",
	}

	backing := NewStaticMapPropertyBacking(propertyMap)
	doTest(t, backing)
}

func TestFileSystemPropertyBacking(t *testing.T) {
	dir, err := ioutil.TempDir("", "example")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(dir) // nolint: errcheck

	content := []byte("cool")
	tmpfn := filepath.Join(dir, "string.property")
	if err := ioutil.WriteFile(tmpfn, content, 0666); err != nil {
		t.Fatal(err)
	}

	backing := NewFileSystemPropertyBacking(dir)
	doTest(t, backing)

	/// update value, look again (should see no change)
	if err := ioutil.WriteFile(tmpfn, []byte("new value"), 0666); err != nil {
		t.Fatal(err)
	}
	doTest(t, backing)

	// refresh, look again (should see new value)
	backing.Refresh()
	value, present := backing.GetProperty("string.property")
	assert.Equal(t, value, "new value")
	assert.True(t, present)
}

func doTest(t *testing.T, backing PropertyBacking) {
	value, present := backing.GetProperty("string.property")
	assert.Equal(t, value, "cool")
	assert.True(t, present)

	value, present = backing.GetProperty("nosuch.property")
	assert.Equal(t, value, "")
	assert.False(t, present)
}
