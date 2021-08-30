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
			return fmt.Errorf("There was an error when accessing the following path: %s: %v.", path, err)
		}
		if path != root && len(path) > len(root)+1 {
			stat, err := os.Stat(path)
			if err != nil {
				return fmt.Errorf("There was an error when accessing the following path: %s: %v.", path, err)
			}

			if stat.Size() > fileSizeLimit {
				return fmt.Errorf("The %s files is too big. The file is %d KB and the max is %d KB.", path, stat.Size(), fileSizeLimit)
			}

			content, err := ioutil.ReadFile(path)
			if err != nil {
				return fmt.Errorf("There was an error when accessing the following path: %s: %v.", path, err)
			}
			propName := path[len(root)+1:]
			props[propName] = string(content)
		}
		return nil
	}

	if err := filepath.Walk(root, visit); err != nil {
		return nil, fmt.Errorf("There was an error when accessing the following path: %s: %v.", root, err)
	}
	return props, nil
}
