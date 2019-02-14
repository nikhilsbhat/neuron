// This package takes care of registering flags,subcommands and returns the
// command to the function who creates or holds the root command.
package commands

import (
	"encoding/json"
	"fmt"
	svcreate "github.com/nikhilsbhat/neuron/cloudoperations/server/create"
	svdelete "github.com/nikhilsbhat/neuron/cloudoperations/server/delete"
	svget "github.com/nikhilsbhat/neuron/cloudoperations/server/get"
	svupdate "github.com/nikhilsbhat/neuron/cloudoperations/server/update"
	err "github.com/nikhilsbhat/neuron/error"
	"github.com/spf13/cobra"
	"strings"
)

var (
	servercmds map[string]*cobra.Command
	createsv   = svcreate.New()
	deletesv   = svdelete.New()
	updatesv   = svupdate.New()
	getsv      = svget.New()
)

// The function that helps in registering the subcommands with the respective main command.
// Make sure you call this, and this is the only way to register the subcommands.
func svRegister(name string, fn *cobra.Command) {
	if servercmds == nil {
		servercmds = make(map[string]*cobra.Command)
	}

	if servercmds[name] != nil {
		panic(fmt.Sprintf("Command %s is already registered", name))
	}
	servercmds[name] = fn
}

// The only way to create server command is to call this function and
// package commands will take care of calling this.
func getServCmds() *cobra.Command {

	// Creating "server" happens here.
	var cmdServer = &cobra.Command{
		Use:   "server [flags]",
		Short: "command to carry out server activities",
		Long:  `This will help user to create/update/get/delete server in cloud.`,
		Run:   cc.echoServer,
	}
	registersvFlags("server", cmdServer)

	for cmdname, cmd := range servercmds {
		cmdServer.AddCommand(cmd)
		registersvFlags(cmdname, cmd)
	}
	return cmdServer
}

// Registering all the flags to the subcommands and command netwrok itself.
func registersvFlags(cmdname string, cmd *cobra.Command) {

	switch strings.ToLower(cmdname) {
	case "servercreate":
		cmd.Flags().StringVarP(&createsv.InstanceName, "name", "n", "", "give a name to your network.")
		cmd.Flags().Int64VarP(&createsv.Count, "count", "", 1, "specify the number of servers that has to be provisioned.")
		cmd.Flags().StringVarP(&createsv.ImageId, "imageid", "i", "", "ID of the base image from which the new server has to be provisioned.")
		cmd.Flags().StringVarP(&createsv.SubnetId, "subnetid", "s", "", "ID of the subnet in which servers has to be created.")
		cmd.Flags().StringVarP(&createsv.KeyName, "keyname", "k", "", "name of the kay-pair which has to be assigned to instances so it will be helpful while logging into it (works only with aws).")
		cmd.Flags().StringVarP(&createsv.Flavor, "flavor", "f", "", "flavor/configuration of the vm that has to be created. (checkout 'neuron flavor list' for the list of available flavors.)")
		cmd.Flags().StringVarP(&createsv.UserData, "userdata", "", "echo 'from neuron'", "if in case you need to execute certain scripts such as shell,ruby on the startup of server.? pass it from this flag.")
		cmd.Flags().BoolVarP(&createsv.AssignPubIp, "assignpublicip", "", false, "turnn this flag on if you need public ip for the machines which will be created.")
	case "serverdelete":
		cmd.Flags().StringVarP(&deletesv.VpcId, "vpcid", "v", "", "pass ID of vpc, from which servers has to be deleted")
	case "serverupdate":
		cmd.Flags().StringVarP(&updatesv.Action, "action", "", "", "action to be performed on the instances (supports start/stop).")
	case "serverget":
		cmd.Flags().StringSliceVarP(&getsv.VpcIds, "vpcids", "v", nil, "ID's of vpcs/vnets, pass comma separated value (if this flag is on which means you'll get servers in vpcs you mentioned)")
		cmd.Flags().StringSliceVarP(&getsv.SubnetIds, "subnetids", "s", nil, "ID's of subnets to filter the servers. pass comma separated value.")
	case "server":
		cmd.PersistentFlags().StringSliceVarP(&getsv.VpcIds, "instanceids", "", nil, "ID's of servers/instances, pass comma separated value.")
	}
}

func (cm *cliMeta) createServer(cmd *cobra.Command, args []string) {
	if cm.CliSet == false {
		cm.NeuronSaysItsError(err.CliNoStart().Error())
	}
	createsv.Cloud.Name = cm.getCloud(cmd)
	createsv.Cloud.Region = cm.getRegion(cmd)
	createsv.Cloud.Profile = cm.getProfile(cmd)
	createsv.Cloud.GetRaw = cm.getGetRaw(cmd)
	server_response, ser_resp_err := createsv.CreateServer()
	if ser_resp_err != nil {
		cm.NeuronSaysItsError(ser_resp_err.Error())
	} else {
		json_val, _ := json.MarshalIndent(server_response, "", " ")
		cm.NeuronSaysItsInfo(string(json_val))
	}
}

func (cm *cliMeta) deleteServer(cmd *cobra.Command, args []string) {
	if cm.CliSet == false {
		cm.NeuronSaysItsError(err.CliNoStart().Error())
	}
	deletesv.Cloud.Name = cm.getCloud(cmd)
	deletesv.Cloud.Region = cm.getRegion(cmd)
	deletesv.Cloud.Profile = cm.getProfile(cmd)
	deletesv.Cloud.GetRaw = cm.getGetRaw(cmd)
	delete_sv_response, sv_err := deletesv.DeleteServer()
	if sv_err != nil {
		cm.NeuronSaysItsError(sv_err.Error())
	} else {
		json_val, _ := json.MarshalIndent(delete_sv_response, "", " ")
		cm.NeuronSaysItsInfo(string(json_val))
	}
}

func (cm *cliMeta) getServer(cmd *cobra.Command, args []string) {
	if cm.CliSet == false {
		cm.NeuronSaysItsError(err.CliNoStart().Error())
	}

	getsv.Cloud.Name = cm.getCloud(cmd)
	getsv.Cloud.Region = cm.getRegion(cmd)
	getsv.Cloud.Profile = cm.getProfile(cmd)
	getsv.Cloud.GetRaw = cm.getGetRaw(cmd)

	if cm.isAll(cmd) {
		get_server_response, sv_get_err := getsv.GetAllServers()
		if sv_get_err != nil {
			cm.NeuronSaysItsError(sv_get_err.Error())
		} else {
			json_val, _ := json.MarshalIndent(get_server_response, "", " ")
			cm.NeuronSaysItsInfo(string(json_val))
		}
	}

	get_server_response, sv_get_err := getsv.GetServersDetails()
	if sv_get_err != nil {
		cm.NeuronSaysItsError(sv_get_err.Error())
	} else {
		json_val, _ := json.MarshalIndent(get_server_response, "", " ")
		cm.NeuronSaysItsInfo(string(json_val))
	}
}

func (cm *cliMeta) updateServer(cmd *cobra.Command, args []string) {
	if cm.CliSet == false {
		cm.NeuronSaysItsError(err.CliNoStart().Error())
	}
	updatesv.Cloud.Name = cm.getCloud(cmd)
	updatesv.Cloud.Region = cm.getRegion(cmd)
	updatesv.Cloud.Profile = cm.getProfile(cmd)
	updatesv.Cloud.GetRaw = cm.getGetRaw(cmd)
	sv_update_response, sv_up_err := updatesv.UpdateServers()
	if sv_up_err != nil {
		cm.NeuronSaysItsError(sv_up_err.Error())
	} else {
		json_val, _ := json.MarshalIndent(sv_update_response, "", " ")
		cm.NeuronSaysItsInfo(string(json_val))
	}
}

func (cm *cliMeta) echoServer(cmd *cobra.Command, args []string) {
	if cm.CliSet == false {
		cm.NeuronSaysItsError(err.CliNoStart().Error())
	}
	cm.printMessage()
	cmd.Usage()
}
