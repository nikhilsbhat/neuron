// The package which makes the tool cloud agnostic for fetching network details.
// The decision will be made here to route the request to respective package based on input.
package networkGet

// The struct which impliments method GetNetworks, GetSubnets.
type GetNetworksInput struct {
	// Ids or names of VPC's of which the information has to be fetched.
	VpcIds []string `json:"VpcIds"`

	// Ids or names of the SUBNET's of which the informaion has to be fetched.
	SubnetIds []string `json:"SubnetIds"`

	// Pass the cloud in which the resource has to be created. usage: "aws","azure" etc.
	Cloud string `json:"Cloud"`

	// Along with cloud, pass region in which resource has to be created.
	Region string `json:"Region"`

	// Passing the profile is important, because this will help in fetching the the credentials
	// of cloud stored along with user details.
	Profile string `json:"Profile"`

	// Use this option if in case you need unfiltered output from cloud.
	GetRaw bool `json:"GetRaw"`
}

//Nothing much from this file. This file contains only the structs for network/get.
