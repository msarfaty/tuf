package state

import (
	"errors"
	"fmt"
	"os"
	"path"
	"path/filepath"

	"github.com/google/uuid"
	"mikesarfaty.com/tuf/pkg/file"
)

type WorkspaceFile struct {
	name string `yaml:"name"`
	md5  string `yaml:"md5"`
}

type Workspace struct {
	uuid    string           `yaml:"guid"`
	abspath string           `yaml:"absolutePath"`
	files   []*WorkspaceFile `yaml:"files"`
}

type Workspaces struct {
	workspaces []*Workspace `yaml:"workspaces"`
}

func (wss *Workspaces) AddWorkspace(path string) (*Workspace, error) {
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

	md5s, err := md5ForTerraformFiles(ret.abspath)
	if err != nil {
		return nil, fmt.Errorf("generating md5 for terraform files: %w", err)
	}
	ret.files = []*WorkspaceFile{}
	for terraformFilePath, md5 := range md5s {
		ret.files = append(ret.files, &WorkspaceFile{
			name: filepath.Base(terraformFilePath),
			md5:  md5,
		})
	}

	ret.uuid = uuid.NewString()
	return &ret, nil
}

// validates a workspace to ensure that the stored state of the workspace matches the current working state of the workspace
func (ws *Workspace) Validate() error {
	actualMd5s, err := md5ForTerraformFiles(ws.abspath)
	if err != nil {
		return fmt.Errorf("generating md5 for terraform files in %s: %w", ws.abspath, err)
	}
	errs := []error{}

	for _, wsFile := range ws.files {
		if actualMd5s[path.Join(ws.abspath, wsFile.name)] != wsFile.md5 {
			errs = append(errs, fmt.Errorf("validation failed for %s (stored md5 %s does not match actual md5 %s)",
				wsFile.name,
				wsFile.md5,
				actualMd5s[path.Join(ws.abspath, wsFile.name)]))
		}
	}
	if len(ws.files) != len(actualMd5s) {
		errs = append(errs, fmt.Errorf(": expected %d files but found %d", len(ws.files), len(actualMd5s)))
	}

	return errors.Join(errs...)
}

// the culmination of multiple workspace validations
func (wss *Workspaces) Validate() error {
	errs := []error{}

	for _, ws := range wss.workspaces {
		err := ws.Validate()
		if err != nil {
			errs = append(errs, fmt.Errorf("failed validation for workspace %s: %w", ws.uuid, err))
		}
	}

	return errors.Join(errs...)
}

// finds all terraform files in a directory and generates their md5s, returning the mapping from abspath:md5
func md5ForTerraformFiles(dir string) (map[string]string, error) {
	absPath, err := filepath.Abs(dir)
	if err != nil {
		return nil, fmt.Errorf("getting abspath for %s: %w", dir, err)
	}

	terraformFileAbsPaths, err := file.GetAllTerraformFilesInDirectory(absPath)
	if err != nil {
		return nil, fmt.Errorf("getting terraform files in directory=[%s]: %w", absPath, err)
	}
	md5s, err := file.GenerateMd5ForFiles(terraformFileAbsPaths)
	if err != nil {
		return nil, fmt.Errorf("generating md5 for terraform files: %w", err)
	}
	return md5s, nil
}
