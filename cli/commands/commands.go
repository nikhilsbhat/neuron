package commands

import (
	"github.com/spf13/cobra"
	/*ci "neuron/cli/commands/ci"
	  data "neuron/cli/commands/database"
	  image "neuron/cli/commands/images"
	  load "neuron/cli/commands/loadbalancers"*/
	network "neuron/cli/commands/networks"
	//server "neuron/cli/commands/servers"
)

func init() {

	network.NwRegister("networkCreate", &cobra.Command{
		Use:   "create [create a network]",
		Short: "command to create a complete network",
		Long:  `This will help user to create network in cloud which he wants.`,
		Run:   network.CreateNetwork,
	})
	network.NwRegister("networkDelete", &cobra.Command{
		Use:   "delete [delete a network]",
		Short: "command to delete a complete network or its components",
		Long:  `This will help user to delete network in cloud which he wants.`,
		Run:   network.DeleteNetwork,
	})
	network.NwRegister("networkGet", &cobra.Command{
		Use:   "get [get a network]",
		Short: "command to get the details of network and its components",
		Long:  `This will help user to get network in cloud which he wants.`,
		Run:   network.GetNetwork,
	})
	network.NwRegister("networkUpdate", &cobra.Command{
		Use:   "update [update a network]",
		Short: "command to update the network and its components",
		Long: `This will help user to update network in cloud which he wants,
                                by letting one to create various components in it.`,
		Run: network.UpdateNetwork,
	})
}
