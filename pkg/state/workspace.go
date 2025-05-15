package state

import (
	"errors"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"strings"

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

type WorkspaceMgr struct {
	workspaces []*Workspace `yaml:"workspaces"`
}

func (wsmgr *WorkspaceMgr) String() string {
	workspaces := []string{}
	for _, ws := range wsmgr.workspaces {
		workspaces = append(workspaces, fmt.Sprintf("%v", ws))
	}

	return fmt.Sprintf("WorkspaceMgr{workspaces=[%s]}", strings.Join(workspaces, ", "))
}

func (ws *Workspace) String() string {
	var sb strings.Builder
	sb.WriteString("Workspace{")
	sb.WriteString(fmt.Sprintf("uuid=%s abspath=%s", ws.uuid, ws.abspath))
	files := []string{}
	for _, f := range ws.files {
		files = append(files, fmt.Sprintf("%v", f))
	}
	sb.WriteString(fmt.Sprintf(" files=[%s]", strings.Join(files, ", ")))
	sb.WriteString("}")
	return sb.String()
}

func (wsf *WorkspaceFile) String() string {
	return fmt.Sprintf("name=%s md5=%s", wsf.name, wsf.md5)
}

// Add a workspace to the workspaces
func (wsmgr *WorkspaceMgr) AddWorkspace(path string) error {
	fmt.Printf("Wsmgr has %d workspaces", len(wsmgr.workspaces))
	ws := Workspace{}

	abspath, err := filepath.Abs(path)
	if err != nil {
		return fmt.Errorf("failed to get absolute path for %s: %w", path, err)
	}
	stat, err := os.Stat(abspath)
	if err != nil {
		return fmt.Errorf("failed to stat %s: %w", abspath, err)
	}
	if !stat.IsDir() {
		return fmt.Errorf("non-directory workspace provided (%s)", abspath)
	}
	ws.abspath = abspath

	md5s, err := md5ForTerraformFiles(ws.abspath)
	if err != nil {
		return fmt.Errorf("generating md5 for terraform files: %w", err)
	}
	ws.files = []*WorkspaceFile{}
	for terraformFilePath, md5 := range md5s {
		ws.files = append(ws.files, &WorkspaceFile{
			name: filepath.Base(terraformFilePath),
			md5:  md5,
		})
	}

	ws.uuid = uuid.NewString()
	wsmgr.workspaces = append(wsmgr.workspaces, &ws)
	fmt.Printf("Wsmgr has %d workspaces", len(wsmgr.workspaces))
	return nil
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
func (wss *WorkspaceMgr) Validate() error {
	errs := []error{}

	for _, ws := range wss.workspaces {
		err := ws.Validate()
		if err != nil {
			errs = append(errs, fmt.Errorf("failed validation for workspace %s: %w", ws.uuid, err))
		}
	}

	return errors.Join(errs...)
}

// yields a new workspaces object with no workspaces
func NewWorkspaceMgr() *WorkspaceMgr {
	return &WorkspaceMgr{
		workspaces: []*Workspace{},
	}
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
