package networkCreate

type NetworkCreateInput struct {
	Name    string   `json:"Name"`
	VpcCidr string   `json:"VpcCidr"`
	SubCidr []string `json:"SubCidr"`
	Type    string   `json:"Type"`
	Ports   []string `json:"Ports"`
	Cloud   string   `json:"Cloud"`
	Region  string   `json:"Region"`
	Profile string   `json:"Profile"`
	GetRaw  bool     `json:"GetRaw"`
}

//Nothing much from this file. This file contains only the structs for network/create
