package testutils

import (
	"os"
	"path/filepath"
	"testing"
)

// Options for temporary directory creation
type TempDirOpts struct {
	// the contents of the temp dir; fileName=contents
	Contents map[string]string
}

// Makes a temp directory with whatever options and returns its name
func MakeDirectory(t *testing.T, tdo *TempDirOpts) string {
	ret := t.TempDir()
	if tdo == nil {
		tdo = &TempDirOpts{}
	}

	for name, contents := range tdo.Contents {
		func() {
			file, err := os.OpenFile(filepath.Join(ret, name), os.O_CREATE|os.O_WRONLY, 0644)
			if err != nil {
				t.Fatalf("failed to create file %s: %v", name, err)
			}
			defer file.Close()
			_, err = file.Write([]byte(contents))
			if err != nil {
				t.Fatalf("error writing contents to file: %v", err)
			}
		}()

	}

	return ret
}
