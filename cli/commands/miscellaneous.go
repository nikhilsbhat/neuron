package commands

import (
	"encoding/json"
	"fmt"
	misc "github.com/nikhilsbhat/neuron-cloudy/cloudoperations/miscellaneous"
	err "github.com/nikhilsbhat/neuron/error"
	"github.com/spf13/cobra"
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
		Use:    "common [flags]",
		Short:  "command for miscellaneous operation",
		Long:   `This will help you to perform miscellaneous operation which we call on the cloud you wish.`,
		Run:    cc.echoCommon,
		Hidden: true,
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

	getregionsresponse, regiongeterr := getrg.GetRegions()
	if regiongeterr != nil {
		cm.NeuronSaysItsError(regiongeterr.Error())
	} else {
		jsonval, _ := json.MarshalIndent(getregionsresponse, "", " ")
		cm.NeuronSaysItsInfo(string(jsonval))
	}
}

func (cm *cliMeta) echoCommon(cmd *cobra.Command, args []string) {
	if cm.CliSet == false {
		cm.NeuronSaysItsError(err.CliNoStart().Error())
	}
	cm.printMessage()
	cmd.Usage()
}
