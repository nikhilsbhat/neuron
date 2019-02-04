package imageCreate

import (
	cmn "neuron/cloudoperations"
)

type CreateImageInput struct {
	InstanceIds []string `json:"instanceids"`
	Cloud cmn.Cloud
}

//Nothing much from this file. This file contains only the structs for image/create
