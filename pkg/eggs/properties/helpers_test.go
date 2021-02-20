package properties

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLoadPropertiesFromFilesystem(t *testing.T) {
	props, err := loadPropertiesFromFilesystem("/no/such")
	assert.Nil(t, props)
	assert.Error(t, err)

	dir, err := ioutil.TempDir("", "example")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(dir) // nolint: errcheck

	content := []byte("cool value")
	tmpfn := filepath.Join(dir, "string.property")
	if err := ioutil.WriteFile(tmpfn, content, 0666); err != nil {
		t.Fatal(err)
	}

	props, err = loadPropertiesFromFilesystem(dir)
	assert.NoError(t, err)

	value, present := props["string.property"]
	assert.True(t, present)
	assert.Equal(t, value, "cool value")
}
