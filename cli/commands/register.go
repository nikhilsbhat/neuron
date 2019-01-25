package commands

import (
	"fmt"
	"github.com/spf13/cobra"
	err "neuron/error"
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
	neucmd.commands = append(neucmd.commands, getInitCmds())
	neucmd.commands = append(neucmd.commands, getNetCmds())
	neucmd.commands = append(neucmd.commands, getServCmds())
	neucmd.commands = append(neucmd.commands, getVersionCmds())
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

func SetNeuronCmds() *cobra.Command {
	cmd := getCmds()
	return cmd
}

func getNeuronCmds() *cobra.Command {

	var neuronCmd = &cobra.Command{
		Use:   "neuron [command]",
		Short: "command to deal with neuron and its activities",
		Long:  `This will help user to use neuron to get things done in various aspects including cloud/ci/database etc.`,
		RunE:  cc.echoNeuron,
	}
	neuronCmd.SetUsageTemplate(getUsageTemplate())
	registerFlags(neuronCmd)
	return neuronCmd
}

func (cm *cliMeta) echoNeuron(cmd *cobra.Command, args []string) error {
	if cm.CliSet == false {
		return err.CliNoStart()
	}
	cmd.Usage()
	return nil
}

// This function will return the custom template for usage function,
// only functions/methods inside this package can call this.

func getUsageTemplate() string {
	return `Usage:{{if .Runnable}}
  {{.UseLine}}{{end}}{{if gt (len .Aliases) 0}}{{printf "\n" }}
Aliases:
  {{.NameAndAliases}}{{end}}{{if .HasExample}}{{printf "\n" }}
Examples:
{{.Example}}{{end}}{{if .HasAvailableSubCommands}}{{printf "\n"}}
Available Commands:{{range .Commands}}{{if (or .IsAvailableCommand (eq .Name "help"))}}
  {{rpad .Name .NamePadding }} {{.Short}}{{end}}{{end}}{{end}}{{if .HasAvailableLocalFlags}}{{printf "\n"}}
Flags:
{{.LocalFlags.FlagUsages | trimTrailingWhitespaces}}{{end}}{{if .HasAvailableInheritedFlags}}{{printf "\n"}}
Global Flags:
{{.InheritedFlags.FlagUsages | trimTrailingWhitespaces}}{{end}}{{if .HasHelpSubCommands}}{{printf "\n"}}
Additional help topics:{{range .Commands}}{{if .IsAdditionalHelpTopicCommand}}
  {{rpad .CommandPath .CommandPathPadding}} {{.Short}}{{end}}{{end}}{{end}}{{if .HasAvailableSubCommands}}{{printf "\n"}}
Use "{{.CommandPath}} [command] --help" for more information about a command.{{end}}"
{{printf "\n"}}`
}
