package state

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/google/uuid"
	"gopkg.in/yaml.v3"
)

const TUF_STATE_FILE = "tuf.state"

// A WorkspaceMgr keeps track of the tuf migration
type WorkspaceMgr struct {
	Workspaces        []*Workspace       `yaml:"workspaces"`
	TerraformMetadata *TerraformMetadata `yaml:"terraform"`
}

// represents this workspacemgr as a string
func (wsmgr *WorkspaceMgr) String() string {
	workspaces := []string{}
	for _, ws := range wsmgr.Workspaces {
		workspaces = append(workspaces, fmt.Sprintf("%v", ws))
	}

	return fmt.Sprintf("WorkspaceMgr{workspaces=[%s]}", strings.Join(workspaces, ", "))
}

// Add a workspace to the workspaces
func (wsmgr *WorkspaceMgr) AddWorkspace(path string) error {
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
	ws.Abspath = abspath

	md5s, err := md5ForTerraformFiles(ws.Abspath)
	if err != nil {
		return fmt.Errorf("generating md5 for terraform files: %w", err)
	}
	ws.Files = []*WorkspaceFile{}
	for terraformFilePath, md5 := range md5s {
		ws.Files = append(ws.Files, &WorkspaceFile{
			Name: filepath.Base(terraformFilePath),
			Md5:  md5,
		})
	}

	ws.Uuid = uuid.NewString()
	wsmgr.Workspaces = append(wsmgr.Workspaces, &ws)
	return nil
}

// the culmination of multiple workspace validations
func (wss *WorkspaceMgr) Validate() error {
	errs := []error{}

	for _, ws := range wss.Workspaces {
		err := ws.Validate()
		if err != nil {
			errs = append(errs, fmt.Errorf("failed validation for workspace %s: %w", ws.Uuid, err))
		}
	}

	return errors.Join(errs...)
}

// yields a new workspaces object with no workspaces
func NewWorkspaceMgr() *WorkspaceMgr {
	return &WorkspaceMgr{
		Workspaces:        []*Workspace{},
		TerraformMetadata: NewTerraformMetadata(),
	}
}

// Write the state of this WorkspaceMgr to disk
func (wsmgr *WorkspaceMgr) WriteToDisk() error {
	_, err := os.Stat(TUF_STATE_FILE)
	if err == nil {
		return fmt.Errorf("found existing tuf state file at %s", TUF_STATE_FILE)
	}
	if !errors.Is(err, fs.ErrNotExist) {
		return fmt.Errorf("could not check for existing tuf state file (unexpectedly): %v", err)
	}

	data, err := yaml.Marshal(wsmgr)
	if err != nil {
		return fmt.Errorf("could not marshal wsmgr to yaml: %v", err)
	}

	err = os.WriteFile(TUF_STATE_FILE, data, 0644)
	if err != nil {
		return fmt.Errorf("failed to write tuf state to file: %v", err)
	}

	return nil
}
