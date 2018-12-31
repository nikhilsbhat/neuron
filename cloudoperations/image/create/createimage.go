package imageCreate

import (
	"fmt"
	auth "neuron/cloud/aws/interface"
	image "neuron/cloud/aws/operations/image"
	awssess "neuron/cloud/aws/sessions"
	common "neuron/cloudoperations/common"
	log "neuron/logger"
	"strings"
)

type CreateImageResponse struct {
	AwsResponse     []image.ImageResponse `json:"AwsResponse,omitempty"`
	AzureResponse   string                `json:"AzureResponse,omitempty"`
	DefaultResponse string                `json:"DefaultResponse,omitempty"`
}

// being create_image my job is to capture image/take server backup and give back the response who called me
func (img *CreateImageInput) CreateImage() (CreateImageResponse, error) {

	switch strings.ToLower(img.Cloud) {
	case "aws":

		creds, crderr := common.GetCredentials(&common.GetCredentialsInput{Profile: img.Profile, Cloud: img.Cloud})
		if crderr != nil {
			return CreateImageResponse{}, crderr
		}
		// I will establish session so that we can carry out the process in cloud
		session_input := awssess.CreateSessionInput{Region: img.Region, KeyId: creds.KeyId, AcessKey: creds.SecretAccess}
		sess := session_input.CreateAwsSession()

		//authorizing to request further
		authinpt := auth.EstablishConnectionInput{Region: img.Region, Resource: "ec2", Session: sess}

		response_image := make([]image.ImageResponse, 0)

		for _, id := range img.InstanceIds {
			imgcreate := new(image.ImageCreateInput)
			imgcreate.InstanceId = id
			imgcreate.GetRaw = img.GetRaw
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
