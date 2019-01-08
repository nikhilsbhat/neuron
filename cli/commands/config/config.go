// This package takes care of registering flags,subcommands and returns the
// command to the function who creates or holds the root command.
package configcmds

import (
	"fmt"
	"github.com/spf13/cobra"
	neuron "neuron/app"
)

var (
	configcmds map[string]*cobra.Command
	start      = new(neuron.Config)
)

// The only way to create init command is to call this function and
// package commands will take care of calling this.
func GetInitCmds() *cobra.Command {

	// Creating "init" happens here.
	var cmdInit = &cobra.Command{
		Use:   "init [To configure neuron]",
		Short: "command to initializa/configure neuron",
		Long:  `This will help user to bring up neuron with the help of configuration file and make is usable.`,
		Run:   initNeuron,
	}
	registernwFlags(cmdInit)

	return cmdInit
}

func initNeuron(cmd *cobra.Command, args []string) {
	config, prferr := cmd.Flags().GetString("config")
	if prferr != nil {
		fmt.Println("flag profile not passed")
	}
	start.ConfigPath = config
	err := start.ConfigureNeuron()
	if err != nil {
		fmt.Println(err)
	}
}

// Registering all the flags for init command.
func registernwFlags(cmd *cobra.Command) {
	cmd.Flags().BoolVarP(&start.EnableAPI, "enableapi", "e", false, "enable this flag if you wish to enable api for neuron.")
}
