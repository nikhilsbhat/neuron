package imageDelete

type DeleteImageInput struct {
	ImageIds []string `json:"ImageId"`
	Cloud    string   `json:"Cloud"`
	Region   string   `json:"Region"`
	Profile  string   `json:"Profile"`
}

//Nothing much from this file. This file contains only the structs for image/delete
