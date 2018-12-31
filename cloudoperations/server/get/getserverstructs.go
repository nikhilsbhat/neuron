package getServers

type GetServersInput struct {
	InstanceIds []string `json:"InstanceIds"`
	VpcIds      []string `json:"VpcIds"`
	SubnetIds   []string `json:"SubnetIds"`
	Cloud       string   `json:"Cloud"`
	Region      string   `json:"Region"`
	Profile     string   `json:"Profile"`
	GetRaw      bool     `json:"GetRaw"`
}

//Nothing much from this file. This file contains only the structs for server/get
