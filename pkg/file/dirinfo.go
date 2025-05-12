package file

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
)

const EXT_TERRAFORM = "tf"

// retrieves all terraform files in directory (marked by ".tf" extension)
func GetAllTerraformFilesInDirectory(dir string) ([]string, error) {
	stats, err := os.Stat(dir)
	if err != nil {
		return nil, fmt.Errorf("failed to stat %s: %w", dir, err)
	}
	if !stats.IsDir() {
		return nil, fmt.Errorf("cannot get terraform files because supplied path (%s) is not a directory", dir)
	}

	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, fmt.Errorf("failed to get directory contents of %s: %w", dir, err)
	}
	ret := []string{}
	for _, entry := range entries {
		if entry.Type().IsRegular() && filepath.Ext(entry.Name()) == EXT_TERRAFORM {
			ret = append(ret, path.Join(dir, entry.Name()))
		}
	}

	return ret, nil
}
