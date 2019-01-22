package commands

import (
	"github.com/spf13/cobra"
)

type cloudGeneral struct {
	Cloud   string
	Region  string
	Profile string
	GetRaw  bool
	Config  string
	Version string
}

// Registering all the flags to the command neuron itself.
func registerFlags(cmd *cobra.Command) {

	neu := new(cloudGeneral)
	cmd.PersistentFlags().StringVarP(&neu.Cloud, "cloud", "", "", "name of the cloud in which resource has to be created")
	cmd.PersistentFlags().StringVarP(&neu.Region, "region", "r", "", "the region of the cloud selected, because we need the region where the resource has to be created. Because we cannot create things in a random cloud just like that")
	cmd.PersistentFlags().StringVarP(&neu.Profile, "profile", "p", "", "name of the cloud profile saved, so that we can fetch the credentials saved")
	cmd.PersistentFlags().BoolVarP(&neu.GetRaw, "getraw", "", false, "enable this flag if you prefer to get unfiltered response, filtered result will give you a crisp information of the resource")
	cmd.PersistentFlags().StringVarP(&neu.Config, "config", "c", "", "pass the location of config file here, so that neuron gets configured as per the entries in config file.")
	//cmd.PersistentFlags().StringVarP(&neu.Version, "version", "v", "", "Fetch the version of the neuron installed by just enabling this flag.")
}
