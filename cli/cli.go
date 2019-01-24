package neuroncli

import (
	"fmt"
	command "neuron/cli/commands"
	"os"
)

func CliMain() {
	err := Execute(os.Args[1:])
	if err != nil {
		fmt.Println("An error occured")
		os.Exit(1)
	}
}

func Execute(args []string) error {

	cmd := command.SetNeuronCmds()
	cmd.SetArgs(args)
	_, err := cmd.ExecuteC()
	if err != nil {
		return err
	}
	return nil
}
