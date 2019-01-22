// This package takes care of registering flags,subcommands and returns the
// command to the function who creates or holds the root command.
package versioncmds

import (
	"fmt"
	"github.com/spf13/cobra"
	neuron "neuron/version"
)

// The only way to create version command is to call this function and
// package commands will take care of calling this.
func GetVersionCmds() *cobra.Command {

	// Creating "init" happens here.
	var cmdInit = &cobra.Command{
		Use:   "version [To configure neuron]",
		Short: "command to fetch the version of neuron installed",
		Long:  `This will help user to find what version of neuron he/she installed inthe machine.`,
		Run:   versionNeuron,
	}

	return cmdInit
	// since this is a subcommand to get version, it does not need flags.
	// hence not registering any flags.
}

func versionNeuron(cmd *cobra.Command, args []string) {
	fmt.Println("Neuron", neuron.GetVersion())
}
