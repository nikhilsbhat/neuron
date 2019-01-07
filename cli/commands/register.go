package commands

import (
	"github.com/spf13/cobra"
	/*ci "neuron/cli/commands/ci"
	  data "neuron/cli/commands/database"
	  image "neuron/cli/commands/images"
	  load "neuron/cli/commands/loadbalancers"*/
	"fmt"
	network "neuron/cli/commands/networks"
	//server "neuron/cli/commands/servers"
)

var (
	cmds map[string]*cobra.Command
)

type neucmds struct {
	commands []*cobra.Command
}

func Register(name string, fn *cobra.Command) {
	if cmds == nil {
		cmds = make(map[string]*cobra.Command)
	}

	if cmds[name] != nil {
		panic(fmt.Errorf("Command %q is already registered", name))
	}
	cmds[name] = fn
}

func getCmds() *cobra.Command {
	neucmd := new(neucmds)
	neucmd.commands = append(neucmd.commands, network.GetNetCmds())
	//future subcommands will go here

	// This gets the full and final command with all subcommands and flags for neuron
	cmd := neucmd.prepareCmds()
	return cmd
}

func (c *neucmds) prepareCmds() *cobra.Command {
	rootCmd := getNeuronCmds()
	for _, cm := range c.commands {
		rootCmd.AddCommand(cm)
	}
	return rootCmd
}

func SetCmds() *cobra.Command {
	cmd := getCmds()
	return cmd
}

func getNeuronCmds() *cobra.Command {

	var neuronCmd = &cobra.Command{
		Use:   "neuron [to deal with neuron]",
		Short: "command to deal with neuron and its activities",
		Long:  `This will help user to use neuron to get things done in various aspects including cloud/ci/database etc.`,
		Run:   echoNeuron,
	}
	registerFlags(neuronCmd)
	return neuronCmd
}

func echoNeuron(cmd *cobra.Command, args []string) {
	fmt.Println("Welcome to Neuron")
}
