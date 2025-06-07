package init

import (
	"errors"
	"fmt"

	"github.com/msarfaty/tuf/pkg/state"
)

// options for initializing tuf
type Options struct {
	TerraformStatePullCommand string
	Workspaces                []string
	StateFileName             string
}

func (o *Options) validate() error {
	if o.Workspaces == nil {
		return errors.New("must provide at least one workspace for the migration")
	}

	if o.StateFileName == "" {
		return errors.New("must provide the statefile name that the pull command outputs to")
	}

	return nil
}

// initializes a tuf migration using the given options
func TufInit(o Options) error {
	err := o.validate()
	if err != nil {
		return fmt.Errorf("failed to initialize new tuf migration: %v", err)
	}
	wsmgr := state.NewWorkspaceMgr()

	for _, workspacePath := range o.Workspaces {
		if err := wsmgr.AddWorkspace(workspacePath); err != nil {
			return fmt.Errorf("failed to add workspace %s: %v", workspacePath, err)
		}
	}

	wsmgr.TerraformMetadata = &state.TerraformMetadata{
		StatePullCommand: o.TerraformStatePullCommand,
		StateFileName:    o.StateFileName,
	}

	if err = wsmgr.WriteToDisk(); err != nil {
		return err
	}

	return nil
}
