package neuroncli

import (
	"fmt"
	"github.com/spf13/cobra"
	command "github.com/nikhilsbhat/neuron/cli/commands"
	"os"
)

var (
	cmd *cobra.Command
)

func init() {
	cmd = command.SetNeuronCmds()
}

func CliMain() {
	err := Execute(os.Args[1:])
	if err != nil {
		fmt.Println("An error occured")
		os.Exit(1)
	}
}

func Execute(args []string) error {

	cmd.SetArgs(args)
	_, err := cmd.ExecuteC()
	if err != nil {
		return err
	}
	return nil
}
