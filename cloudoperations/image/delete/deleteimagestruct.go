package imageDelete

import (
	cmn "neuron/cloudoperations"
)

type DeleteImageInput struct {
	ImageIds []string `json:"imageids"`
	Cloud cmn.Cloud
}

//Nothing much from this file. This file contains only the structs for image/delete
