package imageDelete

import (
	"fmt"
	auth "github.com/nikhilsbhat/neuron/cloud/aws/interface"
	image "github.com/nikhilsbhat/neuron/cloud/aws/operations/image"
	awssess "github.com/nikhilsbhat/neuron/cloud/aws/sessions"
	common "github.com/nikhilsbhat/neuron/cloudoperations/common"
	support "github.com/nikhilsbhat/neuron/cloudoperations/support"
	log "github.com/nikhilsbhat/neuron/logger"
	"strings"
)

type DeleteImageResponse struct {
	AwsResponse     []image.ImageResponse `json:"AwsResponse,omitempty"`
	AzureResponse   string                `json:"AzureResponse,omitempty"`
	DefaultResponse string                `json:"Response,omitempty"`
}

// being create_loadbalancer my job is to create required loadbalancer and give back the response who called me
func (img *DeleteImageInput) DeleteImage() (DeleteImageResponse, error) {

	if status := support.DoesCloudSupports(strings.ToLower(img.Cloud.Name)); status != true {
		return DeleteImageResponse{}, fmt.Errorf(common.DefaultCloudResponse + "DeleteImage")
	}

	switch strings.ToLower(img.Cloud.Name) {
	case "aws":

		creds, crderr := common.GetCredentials(&common.GetCredentialsInput{Profile: img.Cloud.Profile, Cloud: img.Cloud.Name})
		if crderr != nil {
			return DeleteImageResponse{}, crderr
		}
		// I will establish session so that we can carry out the process in cloud
		session_input := awssess.CreateSessionInput{Region: img.Cloud.Region, KeyId: creds.KeyId, AcessKey: creds.SecretAccess}
		sess := session_input.CreateAwsSession()

		//authorizing to request further
		authinpt := auth.EstablishConnectionInput{Region: img.Cloud.Region, Resource: "ec2", Session: sess}

		delimages := new(image.DeleteImageInput)
		delimages.ImageIds = img.ImageIds
		result, err := delimages.DeleteImage(authinpt)
		if err != nil {
			return DeleteImageResponse{}, err
		}
		response := make([]image.ImageResponse, 0)
		response = append(response, result)
		return DeleteImageResponse{AwsResponse: response}, nil

	case "azure":
		return DeleteImageResponse{}, fmt.Errorf(common.DefaultAzResponse)
	case "gcp":
		return DeleteImageResponse{}, fmt.Errorf(common.DefaultGcpResponse)
	case "openstack":
		return DeleteImageResponse{}, fmt.Errorf(common.DefaultOpResponse)
	default:

		log.Info("")
		log.Error(common.DefaultCloudResponse + "DeleteImage")
		log.Info("")
		return DeleteImageResponse{}, fmt.Errorf(common.DefaultCloudResponse + "GetImage")
	}
}
