// The package which makes the tool cloud agnostic for updating network components.
// The decision will be made here to route the request to respective package based on input.
package networkUpdate

// The struct which impliments method GetNetworks, GetSubnets.
type NetworkUpdateInput struct {

	// The type of resources and the action to be performed in it
	// goes here the detailed inputs is in below struct
	Catageory Catageory

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

type Catageory struct {
	// The type of resource that has to be updated.
	Resource string

	// Select the action that has to be performed on the resource
	// passed in above option.
	Action string

	// Pass the name here for the resource that has to be created.
	Name string

	// The CIDR block which will be used to create VPC and this
	// contains info that how many IP should be present in the network
	// so decide that in prior before calling this.
	VpcCidr string

	// The CIDR for the subnet that has to be created in the VPC.
	// Pass an array of CIDR's and neuron will take care of creating
	// appropriate number of subnets and attaching to created VPC
	SubCidrs []string

	// The type of the network that has to be created, public or private.
	// Accordingly IGW will be created and attached.
	Type string

	// The ports that has to be opened for the network,
	// if not passed, by default 22 will be made open so that
	// one can access machines that will be created inside the created network.
	Ports []string

	// Pass the Id of the vpc here if you select to update a resource inside it.
	VpcId string

	// Pass the zone here if you need to create subnet in the required zone.
	Zone string
}

//Nothing much from this file. This file contains only the structs for network/update
