// This package takes care of registering flags,subcommands and returns the
// command to the function who creates or holds the root command.
package commands

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
func getInitCmds() *cobra.Command {

        // Creating "init" happens here.
        var cmdInit = &cobra.Command{
                Use:   "init [To initialize neuron]",
                Short: "command to initialize/configure neuron",
                Long:  `This will help user to bring up neuron with the help of configuration file and make it usable.`,
                Run:   initNeuron,
        }
        registerinitFlags(cmdInit)

        return cmdInit
}

// The only way to create config command is to call this function and
// package commands will take care of calling this.
func getConfigCmds() *cobra.Command {

        // Creating "init" happens here.
        var cmdConfig = &cobra.Command{
                Use:   "config [To configure neuron]",
                Short: "command to configure neuron",
                Long:  `This will help user to configure neuron with the help of configuration file to make is callable from cli.`,
                Run:   configNeuron,
        }

        return cmdConfig
}

func initNeuron(cmd *cobra.Command, args []string) {
        config, prferr := cmd.Flags().GetString("config")
        if prferr != nil {
                fmt.Println("flag profile not passed")
        }
        start.ConfigPath = config
        err := start.Init()
        if err != nil {
                fmt.Println(err)
        }
}

func configNeuron(cmd *cobra.Command, args []string) {
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
func registerinitFlags(cmd *cobra.Command) {
        cmd.Flags().BoolVarP(&start.EnableAPI, "enableapi", "e", false, "enable this flag if you wish to enable api for neuron.")
}
