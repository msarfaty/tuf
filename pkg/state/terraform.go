package state

const (
	DEFAULT_STATE_PULL_COMMAND = ""
	DEFAULT_STATE_FILE_NAME    = "terraform.tfstate"
)

// Metadata for interacting with Terraform
type TerraformMetadata struct {
	// the command to pull terraform state for a given workspace
	StatePullCommand string `yaml:"statePullCommand"`
	// the name of the state file within a given workspace
	StateFileName string `yaml:"stateFileName"`
}

// Initializes a new, default TerraformMetadata
func NewTerraformMetadata() *TerraformMetadata {
	return &TerraformMetadata{
		StatePullCommand: DEFAULT_STATE_PULL_COMMAND,
		StateFileName:    DEFAULT_STATE_FILE_NAME,
	}
}
