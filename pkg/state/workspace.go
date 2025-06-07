package state

import (
	"errors"
	"fmt"
	"path"
	"path/filepath"
	"strings"

	"github.com/msarfaty/tuf/pkg/file"
)

type WorkspaceFile struct {
	Name string `yaml:"name"`
	Md5  string `yaml:"md5"`
}

type Workspace struct {
	Uuid    string           `yaml:"guid"`
	Abspath string           `yaml:"absolutePath"`
	Files   []*WorkspaceFile `yaml:"files"`
}

func (ws *Workspace) String() string {
	var sb strings.Builder
	sb.WriteString("Workspace{")
	sb.WriteString(fmt.Sprintf("uuid=%s abspath=%s", ws.Uuid, ws.Abspath))
	files := []string{}
	for _, f := range ws.Files {
		files = append(files, fmt.Sprintf("%v", f))
	}
	sb.WriteString(fmt.Sprintf(" files=[%s]", strings.Join(files, ", ")))
	sb.WriteString("}")
	return sb.String()
}

func (wsf *WorkspaceFile) String() string {
	return fmt.Sprintf("name=%s md5=%s", wsf.Name, wsf.Md5)
}

// validates a workspace to ensure that the stored state of the workspace matches the current working state of the workspace
func (ws *Workspace) Validate() error {
	actualMd5s, err := md5ForTerraformFiles(ws.Abspath)
	if err != nil {
		return fmt.Errorf("generating md5 for terraform files in %s: %w", ws.Abspath, err)
	}
	errs := []error{}

	for _, wsFile := range ws.Files {
		if actualMd5s[path.Join(ws.Abspath, wsFile.Name)] != wsFile.Md5 {
			errs = append(errs, fmt.Errorf("validation failed for %s (stored md5 %s does not match actual md5 %s)",
				wsFile.Name,
				wsFile.Md5,
				actualMd5s[path.Join(ws.Abspath, wsFile.Name)]))
		}
	}
	if len(ws.Files) != len(actualMd5s) {
		errs = append(errs, fmt.Errorf(": expected %d files but found %d", len(ws.Files), len(actualMd5s)))
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
