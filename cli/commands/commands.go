package commands

import (
	"fmt"
	"github.com/spf13/cobra"
	neuron "neuron/app"
	config "neuron/app/config"
	"os"
)

type cliMeta struct {
	*config.CliMeta
}

var (
	cc = &cliMeta{}
)

func init() {

	meta, clierr := neuron.NeuronCliMeta()
	if clierr != nil {
		fmt.Println(clierr)
		os.Exit(3)
	}
	cc = &cliMeta{&meta}

	nwRegister("networkCreate", &cobra.Command{
		Use:   "create [create a network]",
		Short: "command to create a complete network",
		Long:  `This will help user to create network in cloud which he wants.`,
		RunE:  cc.createNetwork,
	})
	nwRegister("networkDelete", &cobra.Command{
		Use:   "delete [delete a network]",
		Short: "command to delete a complete network or its components",
		Long:  `This will help user to delete network in cloud which he wants.`,
		RunE:  cc.deleteNetwork,
	})
	nwRegister("networkGet", &cobra.Command{
		Use:   "get [get a network]",
		Short: "command to get the details of network and its components",
		Long:  `This will help user to get network in cloud which he wants.`,
		RunE:  cc.getNetwork,
	})
	nwRegister("networkUpdate", &cobra.Command{
		Use:   "update [update a network]",
		Short: "command to update the network and its components",
		Long: `This will help user to update network in cloud which he wants,
                                by letting one to create various components in it.`,
		RunE: cc.updateNetwork,
	})
}
