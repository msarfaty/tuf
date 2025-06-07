/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	tufinit "github.com/msarfaty/tuf/pkg/cli/init"
	"github.com/spf13/cobra"
)

var workspaces []string
var terraformStatePullCommand string
var terraformStateFile string

// initCmd represents the init command
var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize a tuf migration",
	Long: `Initializes multiple workspaces for moving resources and modules between states.
This will:
	- initialize a tuf.state file in the directory you are calling from
	- keep track of all files etc in these workspaces
	- keep track of all operations that occur between these workspaces
	- pull the initial terraform states for use in the migration

Examples:

tuf init --workspaces .,../workspace-a,../workspace-b \
  --terraform-state-pull-command="terraform state pull > state.tfstate" \
  --terraform-state-file="state.tfstate"

* initializes a tuf migration in the current directory, ../workspace-a, and ../workspace-b
* uses "terraform state pull > state.tfstate" as the command to retrieve terraform state
* the state file, in each workspace, will be called state.tfstate

If your states cannot be pulled with a consistent command, you can instead define your own script:

tuf init --workspaces .,../workspace-a,../workspace-b \
  --terraform-state-pull-command="/path/to/my-terraform-pull-command.sh" \
	--terraform-state-file="state.tfstate"

This is still assuming that in each workspace, your state is still being pulled to state.tfstate
`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return tufinit.TufInit(tufinit.Options{
			Workspaces:                workspaces,
			StateFileName:             terraformStateFile,
			TerraformStatePullCommand: terraformStatePullCommand,
		})
	},
}

func init() {
	rootCmd.AddCommand(initCmd)

	initCmd.Flags().StringArrayVar(&workspaces, "workspaces", []string{}, "workspaces to migrate between")
	initCmd.Flags().StringVar(&terraformStatePullCommand, "terraform-state-pull-command", "", "the command to use to pull terraform state")
	initCmd.Flags().StringVar(&terraformStateFile, "terraform-state-file", "terraform.tfstate", "the state file name that the pull command will create")
}
