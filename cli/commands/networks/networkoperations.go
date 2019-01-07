package networkcmds

import (
	"encoding/json"
	"fmt"
	"github.com/spf13/cobra"
	nwcreate "neuron/cloudoperations/network/create"
	nwdelete "neuron/cloudoperations/network/delete"
	nwget "neuron/cloudoperations/network/get"
	nwupdate "neuron/cloudoperations/network/update"
	"os"
)

var (
	createnw = nwcreate.New()
	deletenw = nwdelete.New()
	updatenw = nwupdate.New()
	getnw    = nwget.New()
)

func CreateNetwork(cmd *cobra.Command, args []string) {
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
}

func DeleteNetwork(cmd *cobra.Command, args []string) {
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
}

func GetNetwork(cmd *cobra.Command, args []string) {
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
}

func UpdateNetwork(cmd *cobra.Command, args []string) {
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
}

func echoNetwork(cmd *cobra.Command, args []string) {
	fmt.Printf("I will do nothing, all I do is with the help of my flags.")
	fmt.Printf("Please do pass flags to get the help of this.")
}

func getRegion(cmd *cobra.Command) string {
	reg, regrr := cmd.Flags().GetString("region")
	if regrr != nil {
		fmt.Println("region not passed for the cloud selected")
	}
	return reg
}

func getCloud(cmd *cobra.Command) string {
	cld, clderr := cmd.Flags().GetString("cloud")
	if clderr != nil {
		fmt.Println("flag cloud is empty")
	}
	return cld
}

func getProfile(cmd *cobra.Command) string {
	prf, prferr := cmd.Flags().GetString("profile")
	if prferr != nil {
		fmt.Println("flag profile not passed")
	}
	return prf
}

func getGetRaw(cmd *cobra.Command) bool {
	raw, rwerr := cmd.Flags().GetBool("getraw")
	if rwerr != nil {
		fmt.Println("flag getraw not used")
	}
	return raw
}
