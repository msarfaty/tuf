/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

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

tuf init --workspaces .,../workspace-a,../workspace-b --terraform-state-pull-command="/path/to/my-terraform-pull-command.sh" --terraform-state-file="state.tfstate"

This is still assuming that in each workspace, your state is still being pulled to state.tfstate
`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("init called")
	},
}

func init() {
	rootCmd.AddCommand(initCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// initCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// initCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
