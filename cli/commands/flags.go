package commands

import (
	cmn "github.com/nikhilsbhat/neuron/cloudoperations"
	"github.com/spf13/cobra"
)

type cloudGeneral struct {
	Cloud   cmn.Cloud
	Config  string
	Version string
	all     bool
}

// Registering all the flags to the command neuron itself.
func registerFlags(cmd *cobra.Command) {

	neu := new(cloudGeneral)
	cmd.PersistentFlags().StringVarP(&neu.Cloud.Name, "cloud", "", "", "name of the cloud in which resource has to be created")
	cmd.PersistentFlags().StringVarP(&neu.Cloud.Region, "region", "r", "", "the region of the cloud selected, because we need the region where the resource has to be created. Because we cannot create things in a random cloud just like that")
	cmd.PersistentFlags().StringVarP(&neu.Cloud.Profile, "profile", "p", "", "name of the cloud profile saved, so that we can fetch the credentials saved")
	cmd.PersistentFlags().BoolVarP(&neu.Cloud.GetRaw, "filter", "", false, "enable this flag if you prefer to get filtered response, filtered result will give you a crisp information of the resource unlike raw output")
	cmd.PersistentFlags().StringVarP(&neu.Config, "config", "c", "", "pass the location of config file here, so that neuron gets configured as per the entries in config file")
	cmd.PersistentFlags().BoolVarP(&neu.all, "all", "a", false, "turn this on if you need get result of all. (Note: this is meant to work only for certain commands)")
}
