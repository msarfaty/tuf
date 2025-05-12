package state

import (
	"fmt"
	"os"
	"path/filepath"
)

type WorkspaceFile struct {
	name string `yaml:"name"`
	md5  string `yaml:"md5"`
}

type Workspace struct {
	guid    string           `yaml:"guid"`
	abspath string           `yaml:"absolutePath"`
	files   []*WorkspaceFile `yaml:files`
}

func NewWorkspace(path string) (*Workspace, error) {
	ret := Workspace{}

	abspath, err := filepath.Abs(path)
	if err != nil {
		return nil, fmt.Errorf("failed to get absolute path for %s: %w", path, err)
	}
	stat, err := os.Stat(abspath)
	if err != nil {
		return nil, fmt.Errorf("failed to stat %s: %w", abspath, err)
	}
	if !stat.IsDir() {
		return nil, fmt.Errorf("non-directory workspace provided (%s)", abspath)
	}
	ret.abspath = abspath

	return &ret, nil
}
