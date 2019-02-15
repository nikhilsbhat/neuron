// This package takes care of registering flags,subcommands and returns the
// command to the function who creates or holds the root command.
package commands

import (
	"encoding/json"
	"fmt"
	nwcreate "github.com/nikhilsbhat/neuron/cloudoperations/network/create"
	nwdelete "github.com/nikhilsbhat/neuron/cloudoperations/network/delete"
	nwget "github.com/nikhilsbhat/neuron/cloudoperations/network/get"
	nwupdate "github.com/nikhilsbhat/neuron/cloudoperations/network/update"
	err "github.com/nikhilsbhat/neuron/error"
	"github.com/spf13/cobra"
	"strings"
)

var (
	networkcmds map[string]*cobra.Command
	createnw    = nwcreate.New()
	deletenw    = nwdelete.New()
	updatenw    = nwupdate.New()
	getnw       = nwget.New()
)

// The function that helps in registering the subcommands with the respective main command.
// Make sure you call this, and this is the only way to register the subcommands.
func nwRegister(name string, fn *cobra.Command) {
	if networkcmds == nil {
		networkcmds = make(map[string]*cobra.Command)
	}

	if networkcmds[name] != nil {
		panic(fmt.Sprintf("Command %s is already registered", name))
	}
	networkcmds[name] = fn
}

// The only way to create network command is to call this function and
// package commands will take care of calling this.
func getNetCmds() *cobra.Command {

	// Creating "network" happens here.
	var cmdNetwork = &cobra.Command{
		Use:   "network [network related activities]",
		Short: "command to carry out network activities",
		Long:  `This will help user to create/update/get/delete network in cloud.`,
		Run:   cc.echoNetwork,
	}
	registernwFlags("network", cmdNetwork)

	for cmdname, cmd := range networkcmds {
		cmdNetwork.AddCommand(cmd)
		registernwFlags(cmdname, cmd)
	}
	return cmdNetwork
}

// Registering all the flags to the subcommands and command network itself.
func registernwFlags(cmdname string, cmd *cobra.Command) {

	switch strings.ToLower(cmdname) {
	case "networkcreate":
		cmd.Flags().StringVarP(&createnw.Name, "name", "n", "", "give a name to your network")
		cmd.Flags().StringVarP(&createnw.VpcCidr, "vpcidr", "v", "", "pass CIDR for the vpc to be created ")
		cmd.Flags().StringSliceVarP(&createnw.SubCidr, "subcidr", "s", nil, "pass the CIDR for the subnet to be created, has to be passed in a array. Can pass multiple CIDR's if in case of requirement of multiple subnets")
		cmd.Flags().StringVarP(&createnw.Type, "type", "t", "", "type of network to be created. [public/private are the valid inputs]")
		cmd.Flags().StringSliceVarP(&createnw.Ports, "ports", "", []string{"22"}, "pass the ports in an array to be opened for the network")
	case "networkdelete":
		cmd.Flags().StringSliceVarP(&deletenw.VpcIds, "vpcids", "v", nil, "pass ID's of vpc's, pass comma separated value.")
		cmd.Flags().StringSliceVarP(&deletenw.SubnetIds, "subnetids", "s", nil, "pass ID's of subnets, pass comma separated value")
		cmd.Flags().StringSliceVarP(&deletenw.IgwIds, "igwids", "i", nil, "pass ID's of internet gateway's, pass comma separated value.")
		cmd.Flags().StringSliceVarP(&deletenw.SecurityIds, "secids", "", nil, "pass ID's of secutiry group's, pass comma separated value")
	case "networkupdate":
		cmd.Flags().StringVarP(&updatenw.Resource, "resource", "r", "", "type of resource that has to be updated. vpc/subnet/igw/secutiry etc.")
		cmd.Flags().StringVarP(&updatenw.Name, "name", "n", "", "give a name to your network")
		cmd.Flags().StringVarP(&updatenw.VpcCidr, "vpcidr", "v", "", "pass CIDR for the vpc to be created ")
		cmd.Flags().StringSliceVarP(&updatenw.SubCidrs, "subcidr", "s", nil, "pass the CIDR for the subnet to be created. Can pass multiple CIDR's if in case of requirement of multiple subnets. Pass comma separated value")
		cmd.Flags().StringVarP(&updatenw.Type, "type", "t", "", "type of network to be created. [public/private are the valid inputs]")
		cmd.Flags().StringSliceVarP(&updatenw.Ports, "ports", "p", []string{"22"}, "pass the ports in an array to be opened for the network")
		cmd.Flags().StringVarP(&updatenw.VpcId, "vpcid", "i", "", "ID of vpc, pass comma separated value.")
		cmd.Flags().StringVarP(&updatenw.Zone, "zone", "z", "", "zone in which subnet has to be created")
		cmd.Flags().StringVarP(&updatenw.Action, "action", "a", "", "action that has to be performed on the resource. Ex: create/delete")
	case "networkget":
		cmd.Flags().StringSliceVarP(&getnw.VpcIds, "vpcids", "v", nil, "ID's of vpc's, pass comma separated value.")
		cmd.Flags().StringSliceVarP(&getnw.SubnetIds, "subnetids", "s", nil, "ID's of subnets, pass comma separated value")
	case "network":
		// As of now nothing to go here but as time arrives this case will be filled.
	}
}

func (cm *cliMeta) createNetwork(cmd *cobra.Command, args []string) {
	if cm.CliSet == false {
		cm.NeuronSaysItsError(err.CliNoStart().Error())
	}
	createnw.Cloud.Name = cm.getCloud(cmd)
	createnw.Cloud.Region = cm.getRegion(cmd)
	createnw.Cloud.Profile = cm.getProfile(cmd)
	createnw.Cloud.GetRaw = cm.getGetRaw(cmd)
	server_response, ser_resp_err := createnw.CreateNetwork()
	if ser_resp_err != nil {
		cm.NeuronSaysItsError(ser_resp_err.Error())
	} else {
		json_val, _ := json.MarshalIndent(server_response, "", " ")
		cm.NeuronSaysItsInfo(string(json_val))
	}
}

func (cm *cliMeta) deleteNetwork(cmd *cobra.Command, args []string) {
	if cm.CliSet == false {
		cm.NeuronSaysItsError(err.CliNoStart().Error())
	}
	deletenw.Cloud.Name = cm.getCloud(cmd)
	deletenw.Cloud.Region = cm.getRegion(cmd)
	deletenw.Cloud.Profile = cm.getProfile(cmd)
	deletenw.Cloud.GetRaw = cm.getGetRaw(cmd)
	delete_network_response, net_err := deletenw.DeleteNetwork()
	if net_err != nil {
		cm.NeuronSaysItsError(net_err.Error())
	} else {
		json_val, _ := json.MarshalIndent(delete_network_response, "", " ")
		cm.NeuronSaysItsInfo(string(json_val))
	}
}

func (cm *cliMeta) getNetwork(cmd *cobra.Command, args []string) {
	if cm.CliSet == false {
		cm.NeuronSaysItsError(err.CliNoStart().Error())
	}

	getnw.Cloud.Name = cm.getCloud(cmd)
	getnw.Cloud.Region = cm.getRegion(cmd)
	getnw.Cloud.Profile = cm.getProfile(cmd)
	getnw.Cloud.GetRaw = cm.getGetRaw(cmd)
	get_network_response, net_get_err := getnw.GetNetworks()
	if net_get_err != nil {
		cm.NeuronSaysItsError(net_get_err.Error())
	} else {
		json_val, _ := json.MarshalIndent(get_network_response, "", " ")
		cm.NeuronSaysItsInfo(string(json_val))
	}
}

func (cm *cliMeta) updateNetwork(cmd *cobra.Command, args []string) {
	if cm.CliSet == false {
		cm.NeuronSaysItsError(err.CliNoStart().Error())
	}
	updatenw.Cloud.Name = cm.getCloud(cmd)
	updatenw.Cloud.Region = cm.getRegion(cmd)
	updatenw.Cloud.Profile = cm.getProfile(cmd)
	updatenw.Cloud.GetRaw = cm.getGetRaw(cmd)
	net_update_response, net_up_err := updatenw.UpdateNetwork()
	if net_up_err != nil {
		cm.NeuronSaysItsError(net_up_err.Error())
	} else {
		json_val, _ := json.MarshalIndent(net_update_response, "", " ")
		cm.NeuronSaysItsInfo(string(json_val))
	}
}

func (cm *cliMeta) echoNetwork(cmd *cobra.Command, args []string) {
	if cm.CliSet == false {
		cm.NeuronSaysItsError(err.CliNoStart().Error())
	}
	cm.printMessage()
	cmd.Usage()
}

func (cm *cliMeta) printMessage() {
	fmt.Printf("\n")
	cm.NeuronSaysItsInfo("I will do nothing, all I do is with the help of my flags.\n")
	cm.NeuronSaysItsInfo("Please do pass flags to get help out of me.\n")
}
