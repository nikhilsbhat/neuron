package imageCreate

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

type CreateImageResponse struct {
	AwsResponse     []image.ImageResponse `json:"AwsResponse,omitempty"`
	AzureResponse   string                `json:"AzureResponse,omitempty"`
	DefaultResponse string                `json:"DefaultResponse,omitempty"`
}

// being create_image my job is to capture image/take server backup and give back the response who called me
func (img *CreateImageInput) CreateImage() (CreateImageResponse, error) {

	if status := support.DoesCloudSupports(strings.ToLower(img.Cloud.Name)); status != true {
		return CreateImageResponse{}, fmt.Errorf(common.DefaultCloudResponse + "CreateImage")
	}

	switch strings.ToLower(img.Cloud.Name) {
	case "aws":

		creds, crderr := common.GetCredentials(&common.GetCredentialsInput{Profile: img.Cloud.Profile, Cloud: img.Cloud.Name})
		if crderr != nil {
			return CreateImageResponse{}, crderr
		}
		// I will establish session so that we can carry out the process in cloud
		session_input := awssess.CreateSessionInput{Region: img.Cloud.Region, KeyId: creds.KeyId, AcessKey: creds.SecretAccess}
		sess := session_input.CreateAwsSession()

		//authorizing to request further
		authinpt := auth.EstablishConnectionInput{Region: img.Cloud.Region, Resource: "ec2", Session: sess}

		response_image := make([]image.ImageResponse, 0)

		for _, id := range img.InstanceIds {
			imgcreate := new(image.ImageCreateInput)
			imgcreate.InstanceId = id
			imgcreate.GetRaw = img.Cloud.GetRaw
			response, imgerr := imgcreate.CreateImage(authinpt)
			if imgerr != nil {
				return CreateImageResponse{}, imgerr
			}
			response_image = append(response_image, response)
		}
		return CreateImageResponse{AwsResponse: response_image}, nil

	case "azure":
		return CreateImageResponse{}, fmt.Errorf(common.DefaultAzResponse)
	case "gcp":
		return CreateImageResponse{}, fmt.Errorf(common.DefaultGcpResponse)
	case "openstack":
		return CreateImageResponse{}, fmt.Errorf(common.DefaultOpResponse)
	default:
		log.Info("")
		log.Error(common.DefaultCloudResponse + "CreateImage")
		log.Info("")
		return CreateImageResponse{}, fmt.Errorf(common.DefaultCloudResponse + "CreateImage")
	}
}
