// This package takes care of registering flags,subcommands and returns the
// command to the function who creates or holds the root command.
package networkcmds

import (
	"fmt"
	"github.com/spf13/cobra"
	"strings"
)

var (
	networkcmds map[string]*cobra.Command
)

// The function that helps in registering the subcommands with the respective main command.
// Make sure you call this, and this is the only way to register the subcommands.
func NwRegister(name string, fn *cobra.Command) {
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
func GetNetCmds() *cobra.Command {

	// Creating "network" happens here.
	var cmdNetwork = &cobra.Command{
		Use:   "network [network related activities]",
		Short: "command to carry out network activities",
		Long:  `This will help user to create/update/get/delete network in cloud.`,
		Run:   echoNetwork,
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
