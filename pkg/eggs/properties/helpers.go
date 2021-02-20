package properties

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
)

const (
	fileSizeLimit = 1 * 1024 * 1024 // 1 MiByte
)

func loadPropertiesFromFilesystem(root string) (map[string]string, error) {
	props := make(map[string]string)

	visit := func(path string, f os.FileInfo, err error) error {
		if err != nil {
			return fmt.Errorf("loadPropertiesFromFilesystem: could not walk path '%s': %+v", path, err)
		}
		if path != root && len(path) > len(root)+1 {
			stat, err := os.Stat(path)
			if err != nil {
				return fmt.Errorf("loadPropertiesFromFilesystem: could not stat '%s': %+v", path, err)
			}

			if stat.Size() > fileSizeLimit {
				return fmt.Errorf("loadPropertiesFromFilesystem: '%s' is too large (%d > %d)", path, stat.Size(), fileSizeLimit)
			}

			content, err := ioutil.ReadFile(path)
			if err != nil {
				return fmt.Errorf("loadPropertiesFromFilesystem: could not read '%s': %+v", path, err)
			}
			propName := path[len(root)+1:]
			props[propName] = string(content)
		}
		return nil
	}

	if err := filepath.Walk(root, visit); err != nil {
		return nil, fmt.Errorf("loadPropertiesFromFilesystem: could not walk root '%s': %+v", root, err)
	}
	return props, nil
}
