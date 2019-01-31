// The package which makes the tool cloud agnostic for deleting network and its components.
// The decision will be made here to route the request to respective package based on input.
package networkDelete

// The struct which impliments method DeleteNetwork.
type DeleteNetworkInput struct {

	// Ids or names of VPC's which has to be deleted.
	VpcIds []string

	// Ids or names of SUBNET's which has to be deleted
	SubnetIds []string

	// Ids or name of Internet Gateways which has to be deleted.
	IgwIds []string

	// Ids or name of Security Groups which has to be deletd.
	SecurityIds []string

	// Pass the cloud in which the resource has to be created. usage: "aws","azure" etc.
	Cloud string

	// Along with cloud, pass region in which resource has to be created.
	Region string

	// Passing the profile is important, because this will help in fetching the the credentials
	// of cloud stored along with user details.
	Profile string

	// Use this option if in case you need unfiltered output from cloud.
	GetRaw bool
}

//Nothing much from this file. This file contains only the structs for network/create
