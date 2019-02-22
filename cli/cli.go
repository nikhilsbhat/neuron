// Package neuroncli will initialize cli for neuron.
package neuroncli

import (
	"fmt"
	command "github.com/nikhilsbhat/neuron/cli/commands"
	"github.com/spf13/cobra"
	"os"
)

var (
	cmd *cobra.Command
)

func init() {
	cmd = command.SetNeuronCmds()
}

// CliMain will take the workload of executing/starting the cli, when the command is passed to it.
func CliMain() {
	err := Execute(os.Args[1:])
	if err != nil {
		fmt.Println("An error occurred")
		os.Exit(1)
	}
}

// Execute will actually execute the cli by taking the arguments passed to cli.
func Execute(args []string) error {

	cmd.SetArgs(args)
	_, err := cmd.ExecuteC()
	if err != nil {
		return err
	}
	return nil
}
