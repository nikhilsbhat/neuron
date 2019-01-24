// This package takes care of registering flags,subcommands and returns the
// command to the function who creates or holds the root command.
package commands

import (
	"encoding/json"
	"fmt"
	"github.com/spf13/cobra"
	nwcreate "neuron/cloudoperations/network/create"
	nwdelete "neuron/cloudoperations/network/delete"
	nwget "neuron/cloudoperations/network/get"
	nwupdate "neuron/cloudoperations/network/update"
	err "neuron/error"
	"os"
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
		RunE:  cc.echoNetwork,
	}
	registernwFlags("network", cmdNetwork)

	for cmdname, cmd := range networkcmds {
		cmdNetwork.AddCommand(cmd)
		registernwFlags(cmdname, cmd)
	}
	return cmdNetwork
}

// Registering all the flags to the subcommands and command netwrok itself.
func registernwFlags(cmdname string, cmd *cobra.Command) {

	switch strings.ToLower(cmdname) {
	case "networkcreate":
		cmd.Flags().StringVarP(&createnw.Name, "name", "n", "", "give a name to your network")
		cmd.Flags().StringVarP(&createnw.VpcCidr, "vpcidr", "v", "", "pass CIDR for the vpc to be created ")
		cmd.Flags().StringSliceVarP(&createnw.SubCidr, "subcidr", "s", nil, "pass the CIDR for the subnet to be created, has to be passed in a array. Can pass multiple CIDR's if in case of requirement of multiple subnets")
		cmd.Flags().StringVarP(&createnw.Type, "type", "t", "", "type of network to be created. [public/private are the valid inputs]")
		cmd.Flags().StringSliceVarP(&createnw.Ports, "ports", "p", []string{"22"}, "pass the ports in an array to be opened for the network")
	case "networkdelete":
		cmd.Flags().StringSliceVarP(&deletenw.VpcIds, "vpcids", "v", nil, "pass ID's of vpc's, pass comma separated value.")
		cmd.Flags().StringSliceVarP(&deletenw.SubnetIds, "subnetids", "s", nil, "pass ID's of subnets, pass comma separated value")
		cmd.Flags().StringSliceVarP(&deletenw.IgwIds, "igwids", "i", nil, "pass ID's of internet gateway's, pass comma separated value.")
		cmd.Flags().StringSliceVarP(&deletenw.SecurityIds, "secids", "", nil, "pass ID's of secutiry group's, pass comma separated value")
	case "networkupdate":
		cmd.Flags().StringVarP(&updatenw.Catageory.Resource, "resource", "r", "", "type of resource that has to be updated. vpc/subnet/igw/secutiry etc.")
		cmd.Flags().StringVarP(&updatenw.Catageory.Name, "name", "n", "", "give a name to your network")
		cmd.Flags().StringVarP(&updatenw.Catageory.VpcCidr, "vpcidr", "v", "", "pass CIDR for the vpc to be created ")
		cmd.Flags().StringSliceVarP(&updatenw.Catageory.SubCidrs, "subcidr", "s", nil, "pass the CIDR for the subnet to be created. Can pass multiple CIDR's if in case of requirement of multiple subnets. Pass comma separated value")
		cmd.Flags().StringVarP(&updatenw.Catageory.Type, "type", "t", "", "type of network to be created. [public/private are the valid inputs]")
		cmd.Flags().StringSliceVarP(&updatenw.Catageory.Ports, "ports", "p", []string{"22"}, "pass the ports in an array to be opened for the network")
		cmd.Flags().StringVarP(&updatenw.Catageory.VpcId, "vpcid", "i", "", "ID of vpc, pass comma separated value.")
		cmd.Flags().StringVarP(&updatenw.Catageory.Zone, "zone", "z", "", "zone in which subnet has to be created")
		cmd.Flags().StringVarP(&updatenw.Catageory.Action, "action", "a", "", "action that has to be performed on the resource. Ex: create/delete")
	case "networkget":
		cmd.Flags().StringSliceVarP(&getnw.VpcIds, "vpcids", "v", nil, "ID's of vpc's, pass comma separated value.")
		cmd.Flags().StringSliceVarP(&getnw.SubnetIds, "subnetids", "s", nil, "ID's of subnets, pass comma separated value")
	case "network":
		// As of now nothing to go here but as time arrives this case will be filled.
	}
}

func (cm *cliMeta) createNetwork(cmd *cobra.Command, args []string) error {
	if cm.CliSet == false {
		return err.CliNoStart()
	}
	createnw.Cloud = getCloud(cmd)
	createnw.Region = getRegion(cmd)
	createnw.Profile = getProfile(cmd)
	createnw.GetRaw = getGetRaw(cmd)
	server_response, ser_resp_err := createnw.CreateNetwork()
	if ser_resp_err != nil {
		fmt.Fprintf(os.Stdout, "%v\n", ser_resp_err)
	} else {
		json_val, _ := json.MarshalIndent(server_response, "", " ")
		fmt.Fprintf(os.Stdout, "%v\n", string(json_val))
	}
	return nil
}

func (cm *cliMeta) deleteNetwork(cmd *cobra.Command, args []string) error {
	if cm.CliSet == false {
		return err.CliNoStart()
	}
	deletenw.Cloud = getCloud(cmd)
	deletenw.Region = getRegion(cmd)
	deletenw.Profile = getProfile(cmd)
	deletenw.GetRaw = getGetRaw(cmd)
	delete_network_response, net_err := deletenw.DeleteNetwork()
	if net_err != nil {
		fmt.Fprintf(os.Stdout, "%v\n", net_err)
	} else {
		json_val, _ := json.MarshalIndent(delete_network_response, "", " ")
		fmt.Fprintf(os.Stdout, "%v\n", string(json_val))
	}
	return nil
}

func (cm *cliMeta) getNetwork(cmd *cobra.Command, args []string) error {
	if cm.CliSet == false {
		return err.CliNoStart()
	}
	getnw.Cloud = getCloud(cmd)
	getnw.Region = getRegion(cmd)
	getnw.Profile = getProfile(cmd)
	getnw.GetRaw = getGetRaw(cmd)
	get_network_response, net_get_err := getnw.GetNetworks()
	if net_get_err != nil {
		fmt.Fprintf(os.Stdout, "%v\n", net_get_err)
	} else {
		json_val, _ := json.MarshalIndent(get_network_response, "", " ")
		fmt.Fprintf(os.Stdout, "%v\n", string(json_val))
	}
	return nil
}

func (cm *cliMeta) updateNetwork(cmd *cobra.Command, args []string) error {
	if cm.CliSet == false {
		return err.CliNoStart()
	}
	updatenw.Cloud = getCloud(cmd)
	updatenw.Region = getRegion(cmd)
	updatenw.Profile = getProfile(cmd)
	updatenw.GetRaw = getGetRaw(cmd)
	net_update_response, net_up_err := updatenw.UpdateNetwork()
	if net_up_err != nil {
		fmt.Fprintf(os.Stdout, "%v\n", net_up_err)
	} else {
		json_val, _ := json.MarshalIndent(net_update_response, "", " ")
		fmt.Fprintf(os.Stdout, "%v\n", string(json_val))
	}
	return nil
}

func (cm *cliMeta) echoNetwork(cmd *cobra.Command, args []string) error {
	if cm.CliSet == false {
		return err.CliNoStart()
	}
	fmt.Printf("I will do nothing, all I do is with the help of my flags.")
	fmt.Printf("Please do pass flags to get the help of this.")
	return nil
}
