package internal

import (
	"os"
	"path/filepath"
	"strings"
)

func ProcessDirectory(srcDir string) ([]string, error) {
	var yamlPaths []string

	err := filepath.Walk(srcDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() && strings.HasSuffix(info.Name(), ".yaml") {
			yamlPaths = append(yamlPaths, path)
		}

		return nil
	})

	return yamlPaths, err
}
