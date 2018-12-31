package networkDelete

type DeleteNetworkInput struct {
	VpcIds  []string `json:"VpcIds"`
	Cloud   string   `json:"Cloud"`
	Region  string   `json:"Region"`
	Profile string   `json:"Profile"`
	GetRaw  bool     `json:"GetRaw"`
}

//Nothing much from this file. This file contains only the structs for network/create
