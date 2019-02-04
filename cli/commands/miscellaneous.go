// This package takes care of registering flags,subcommands and returns the
// command to the function who creates or holds the root command.
package commands

import (
	"encoding/json"
	"fmt"
	"github.com/spf13/cobra"
	misc "neuron/cloudoperations/miscellaneous"
	err "neuron/error"
	"strings"
)

var (
	misccmds map[string]*cobra.Command
	getrg    = misc.New()
)

// The function that helps in registering the subcommands with the respective main command.
// Make sure you call this, and this is the only way to register the subcommands.
func msRegister(name string, fn *cobra.Command) {
	if misccmds == nil {
		misccmds = make(map[string]*cobra.Command)
	}

	if misccmds[name] != nil {
		panic(fmt.Sprintf("Command %s is already registered", name))
	}
	misccmds[name] = fn
}

// The only way to create common command is to call this function and
// package commands will take care of calling this.
func getMiscCmds() *cobra.Command {

	// Creating "common" happens here.
	var cmdMisc = &cobra.Command{
		Use:   "common [flags]",
		Short: "command for miscellaneous operation",
		Long:  `This will help you to perform miscellaneous operation which we call on the cloud you wish.`,
		Run:   cc.echoCommon,
	}
	registermiscFlags("server", cmdMisc)

	for cmdname, cmd := range misccmds {
		cmdMisc.AddCommand(cmd)
		registermiscFlags(cmdname, cmd)
	}
	return cmdMisc
}

// Registering all the flags to the subcommands and command common itself.
func registermiscFlags(cmdname string, cmd *cobra.Command) {

	switch strings.ToLower(cmdname) {
	case "regionget":
	case "common":
	}
}

func (cm *cliMeta) getRegions(cmd *cobra.Command, args []string) {
	if cm.CliSet == false {
		cm.NeuronSaysItsError(err.CliNoStart().Error())
	}

	getrg.Cloud.Name = cm.getCloud(cmd)
	getrg.Cloud.Region = cm.getRegion(cmd)
	getrg.Cloud.Profile = cm.getProfile(cmd)
	getrg.Cloud.GetRaw = cm.getGetRaw(cmd)

	get_regions_response, reg_get_err := getrg.GetRegions()
	if reg_get_err != nil {
		cm.NeuronSaysItsError(reg_get_err.Error())
	} else {
		json_val, _ := json.MarshalIndent(get_regions_response, "", " ")
		cm.NeuronSaysItsInfo(string(json_val))
	}
}

func (cm *cliMeta) echoCommon(cmd *cobra.Command, args []string) {
	if cm.CliSet == false {
		cm.NeuronSaysItsError(err.CliNoStart().Error())
	}
	cm.printMessage()
	cmd.Usage()
}
