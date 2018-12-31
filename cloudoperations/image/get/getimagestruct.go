package imagesGet

type GetImagesInput struct {
	ImageIds []string `json:"ImageIds"`
	GetRaw   bool     `json:"GetRaw"`
	Cloud    string   `json:"Cloud"`
	Region   string   `json:"Region"`
	Profile  string   `json:"Profile"`
}

//Nothing much from this file. This file contains only the structs for image/get
