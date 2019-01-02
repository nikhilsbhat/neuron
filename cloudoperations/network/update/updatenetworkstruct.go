package networkUpdate

type NetworkUpdateInput struct {
	Catageory Catageory
	Cloud     string
	Region    string
	Profile   string
	GetRaw    bool
}

type Catageory struct {
	Resource string
	Action   string
	Name     string
	VpcCidr  string
	SubCidrs []string
	Type     string
	Ports    []string
	VpcId    string
	Zone     string
}

//Nothing much from this file. This file contains only the structs for network/update
