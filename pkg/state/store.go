package state

import "fmt"

// A State is simply a way to track all the operations that happen during a tuf run.
// This could be tracking how a block moves throughout workspaces, which modules exist in which directory, etc.
type State interface {
	// Initialize a new State, which tracks the Terraform workspace(s) and how they are changed
	Initialize() error
}

// DiskState is a state implementation where information is stored on disk
type DiskState struct {
	workspaces []string
	State
}

func NewDiskState(workspaces []string) *DiskState {
	return &DiskState{workspaces: workspaces}
}

// Initializes initial state
func (ds *DiskState) Initialize() error {
	seen := map[string]string{}

	for _, w := range ds.workspaces {
		if _, ok := seen[w]; ok {
			return fmt.Errorf("could not initialize state manager: workspace=[%s] was initialized multiple times", w)
		}
	}

	return nil
}
