package imageDelete

import (
	cmn "github.com/nikhilsbhat/neuron/cloudoperations"
)

type DeleteImageInput struct {
	ImageIds []string `json:"imageids"`
	Cloud    cmn.Cloud
}

//Nothing much from this file. This file contains only the structs for image/delete
