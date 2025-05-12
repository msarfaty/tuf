package file

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

// generates an md5 checksum for multiple files and returns them as a key/val pairing
func GenerateMd5ForFiles(paths []string) (map[string]string, error) {

	ret := map[string]string{}
	for _, path := range paths {
		stats, err := os.Stat(path)
		if err != nil {
			return nil, fmt.Errorf("stat failed for %s: %w", path, err)
		}
		if stats.IsDir() {
			return nil, fmt.Errorf("should not be generating md5 for %s (is a directory)", path)
		}

		file, err := os.Open(path)
		if err != nil {
			return nil, fmt.Errorf("issue opening file: %w", err)
		}

		hasher := md5.New()
		bytes, err := io.Copy(hasher, file)
		if err != nil {
			file.Close()
			return nil, fmt.Errorf("failed to copy file in md5 hasher: %w", err)
		}
		fmt.Printf("wrote %d bytes", bytes)

		file.Close()
		ret[filepath.Base(path)] = hex.EncodeToString(hasher.Sum(nil))
	}

	return ret, nil
}
