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
	workspaces []*Workspace `yaml:"workspaces"`
}

// represents this workspacemgr as a string
func (wsmgr *WorkspaceMgr) String() string {
	workspaces := []string{}
	for _, ws := range wsmgr.workspaces {
		workspaces = append(workspaces, fmt.Sprintf("%v", ws))
	}

	return fmt.Sprintf("WorkspaceMgr{workspaces=[%s]}", strings.Join(workspaces, ", "))
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

	err = os.WriteFile(TUF_STATE_FILE, data, os.FileMode(os.O_CREATE))

	return nil
}
