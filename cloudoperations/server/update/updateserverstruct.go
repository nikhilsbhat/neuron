package updateServers

type UpdateServersInput struct {
	InstanceIds []string `json:"InstanceIds"`
	Action      string   `json:"Action"`
	Cloud       string   `json:"Cloud"`
	Region      string   `json:"Region"`
	Profile     string   `json:"Profile"`
	GetRaw      bool     `json:"GetRaw"`
}

//Nothing much from this file. This file contains only the structs for server/update
