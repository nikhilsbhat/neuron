// This package takes care of registering flags,subcommands and returns the
// command to the function who creates or holds the root command.
package commands

import (
	"fmt"
	err "github.com/nikhilsbhat/neuron/error"
	neuron "github.com/nikhilsbhat/neuron/version"
	"github.com/spf13/cobra"
)

// The only way to create version command is to call this function and
// package commands will take care of calling this.
func getVersionCmds() *cobra.Command {

	// Creating "version" happens here.
	var cmdInit = &cobra.Command{
		Use:   "version [To configure neuron]",
		Short: "command to fetch the version of neuron installed",
		Long:  `This will help user to find what version of neuron he/she installed in her machine.`,
		RunE:  cc.versionNeuron,
	}

	return cmdInit
	// since this is a subcommand to get version, it does not need flags.
	// hence not registering any flags.
}

func (cm *cliMeta) versionNeuron(cmd *cobra.Command, args []string) error {
	if cm.CliSet == false {
		return err.CliNoStart()
	}
	fmt.Println("Neuron", neuron.GetVersion())
	return nil
}
