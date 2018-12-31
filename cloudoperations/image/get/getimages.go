package imagesGet

import (
	"fmt"
	auth "neuron/cloud/aws/interface"
	image "neuron/cloud/aws/operations/image"
	awssess "neuron/cloud/aws/sessions"
	common "neuron/cloudoperations/common"
	log "neuron/logger"
	"strings"
)

type GetImagesResponse struct {
	AwsResponse     []image.ImageResponse `json:"AwsResponse,omitempty"`
	AzureResponse   string                `json:"AzureResponse,omitempty"`
	DefaultResponse string                `json:"Response,omitempty"`
}

// being create_loadbalancer my job is to create required loadbalancer and give back the response who called me
func (img *GetImagesInput) GetImage() (GetImagesResponse, error) {

	switch strings.ToLower(img.Cloud) {
	case "aws":

		creds, crderr := common.GetCredentials(&common.GetCredentialsInput{Profile: img.Profile, Cloud: img.Cloud})
		if crderr != nil {
			return GetImagesResponse{}, crderr
		}
		// I will establish session so that we can carry out the process in cloud
		session_input := awssess.CreateSessionInput{Region: img.Region, KeyId: creds.KeyId, AcessKey: creds.SecretAccess}
		sess := session_input.CreateAwsSession()

		//authorizing to request further
		authinpt := auth.EstablishConnectionInput{Region: img.Region, Resource: "ec2", Session: sess}

		getimage := new(image.GetImageInput)
		getimage.ImageIds = img.ImageIds
		getimage.GetRaw = img.GetRaw
		result, err := getimage.GetImage(authinpt)
		if err != nil {
			return GetImagesResponse{}, err
		}
		return GetImagesResponse{AwsResponse: result}, nil

	case "azure":
		return GetImagesResponse{}, fmt.Errorf(common.DefaultAzResponse)
	case "gcp":
		return GetImagesResponse{}, fmt.Errorf(common.DefaultGcpResponse)
	case "openstack":
		return GetImagesResponse{}, fmt.Errorf(common.DefaultOpResponse)
	default:
		log.Info("")
		log.Error(common.DefaultCloudResponse + "GetImage")
		log.Info("")
		return GetImagesResponse{}, fmt.Errorf(common.DefaultCloudResponse + "GetImage")
	}
}

// being create_loadbalancer my job is to create required loadbalancer and give back the response who called me
func (img *GetImagesInput) GetAllImage() (GetImagesResponse, error) {

	switch strings.ToLower(img.Cloud) {
	case "aws":

		creds, crderr := common.GetCredentials(&common.GetCredentialsInput{Profile: img.Profile, Cloud: img.Cloud})
		if crderr != nil {
			return GetImagesResponse{}, crderr
		}
		// I will establish session so that we can carry out the process in cloud
		session_input := awssess.CreateSessionInput{Region: img.Region, KeyId: creds.KeyId, AcessKey: creds.SecretAccess}
		sess := session_input.CreateAwsSession()

		//authorizing to request further
		authinpt := auth.EstablishConnectionInput{Region: img.Region, Resource: "ec2", Session: sess}

		getimages := new(image.GetImageInput)
		getimages.GetRaw = img.GetRaw
		result, err := getimages.GetAllImage(authinpt)
		if err != nil {
			return GetImagesResponse{}, err
		}
		return GetImagesResponse{AwsResponse: result}, nil

	case "azure":
		return GetImagesResponse{}, fmt.Errorf(common.DefaultAzResponse)
	case "gcp":
		return GetImagesResponse{}, fmt.Errorf(common.DefaultGcpResponse)
	case "openstack":
		return GetImagesResponse{}, fmt.Errorf(common.DefaultOpResponse)
	default:
		log.Info("")
		log.Error(common.DefaultCloudResponse + "GetAllImage")
		log.Info("")
		return GetImagesResponse{}, fmt.Errorf(common.DefaultCloudResponse + "GetAllImage")
	}
}
