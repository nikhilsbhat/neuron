package imageCreate

type CreateImageInput struct {
	InstanceIds []string `json:"InstanceIds"`
	GetRaw      bool     `json:"GetRaw"`
	Cloud       string   `json:"Cloud"`
	Region      string   `json:"Region"`
	Profile     string   `json:"Profile"`
}

//Nothing much from this file. This file contains only the structs for image/create
