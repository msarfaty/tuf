package parser

import (
	"fmt"
	"strings"

	"github.com/hashicorp/hcl/v2"
)

// description for any generic hcl block
type BlockDescription interface {
	// Determine if this BlockDescription matches the given HCL block
	Matches(block hcl.Block) bool
	// Determines a useful destination file name for the block
	DestinationFileName() string
	// reverse the address name
	address() string
}

// descriptve characteristics of a module block
type ModuleBlockDescription struct {
	BlockDescription
	// The name of the module
	name string
}

type ResourceBlockDescription struct {
	BlockDescription
	// resource type (ie aws_security_group)
	rType string
	// resource name
	name string
}

// determines if the given hcl block matches the description of this ModuleBlockDescription
func (m *ModuleBlockDescription) Matches(block hcl.Block) bool {
	if block.Type != "module" {
		return false
	}

	if len(block.Labels) != 1 {
		// unexpected case, but do not throw error
		return false
	}

	return block.Labels[0] == m.name
}

func (m *ModuleBlockDescription) DestinationFileName() string {
	return fmt.Sprintf("module_%s.tuf.tf", m.name)
}

func (m *ModuleBlockDescription) address() string {
	return fmt.Sprintf("module.%s", m.name)
}

// determines if the given hcl block matches the description of this ResourceBlockDescription
func (m *ResourceBlockDescription) Matches(block hcl.Block) bool {
	if block.Type != "resource" {
		return false
	}

	if len(block.Labels) != 2 {
		// expect a resource type label and name label, but do not throw error
		return false
	}

	return block.Labels[0] == m.rType && block.Labels[1] == m.name
}

func (m *ResourceBlockDescription) DestinationFileName() string {
	return "resources.tuf.tf"
}

func (m *ResourceBlockDescription) address() string {
	return fmt.Sprintf("%s.%s", m.rType, m.name)
}

// Creates a BlockDescription for module address calls
func newModuleBlockDescription(address string) (*ModuleBlockDescription, error) {
	parts := strings.Split(address, ".")
	if len(parts) != 2 {
		return nil, fmt.Errorf("wrong number of parts to describe address for module %s", address)
	}

	if parts[0] != "module" {
		return nil, fmt.Errorf("cannot make module block description from invalid module address %s", address)
	}

	return &ModuleBlockDescription{name: parts[1]}, nil
}

func newResourceBlockDescription(address string) (*ResourceBlockDescription, error) {
	parts := strings.Split(address, ".")

	if len(parts) != 2 {
		return nil, fmt.Errorf("wrong number of parts to describe address of resource (%s)", address)
	}

	return &ResourceBlockDescription{rType: parts[0], name: parts[1]}, nil
}

// creates a new BlockDescription to aid in finding terraform blocks
func New(address string) (BlockDescription, error) {
	parts := strings.Split(address, ".")

	if len(parts) <= 1 {
		return nil, fmt.Errorf("not enough parts for resource %s", address)
	}

	var bd BlockDescription
	var err error
	switch parts[0] {
	case "module":
		bd, err = newModuleBlockDescription(address)
	default:
		// resource address do not have a static starting path
		bd, err = newResourceBlockDescription(address)
	}

	if err != nil {
		return nil, fmt.Errorf("failed to find a valid description for %s: %w", address, err)
	}

	return bd, nil
}
