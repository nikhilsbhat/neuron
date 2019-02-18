package imagesget

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

type GetImagesResponse struct {
	AwsResponse     []image.ImageResponse `json:"AwsResponse,omitempty"`
	AzureResponse   string                `json:"AzureResponse,omitempty"`
	DefaultResponse string                `json:"Response,omitempty"`
}

// being create_loadbalancer my job is to create required loadbalancer and give back the response who called me
func (img *GetImagesInput) GetImage() (GetImagesResponse, error) {

	if status := support.DoesCloudSupports(strings.ToLower(img.Cloud.Name)); status != true {
		return GetImagesResponse{}, fmt.Errorf(common.DefaultCloudResponse + "GetImages")
	}

	switch strings.ToLower(img.Cloud.Name) {
	case "aws":

		creds, crderr := common.GetCredentials(&common.GetCredentialsInput{Profile: img.Cloud.Profile, Cloud: img.Cloud.Name})
		if crderr != nil {
			return GetImagesResponse{}, crderr
		}
		// I will establish session so that we can carry out the process in cloud
		session_input := awssess.CreateSessionInput{Region: img.Cloud.Region, KeyId: creds.KeyId, AcessKey: creds.SecretAccess}
		sess := session_input.CreateAwsSession()

		//authorizing to request further
		authinpt := auth.EstablishConnectionInput{Region: img.Cloud.Region, Resource: "ec2", Session: sess}

		getimage := new(image.GetImageInput)
		getimage.ImageIds = img.ImageIds
		getimage.GetRaw = img.Cloud.GetRaw
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

	switch strings.ToLower(img.Cloud.Name) {
	case "aws":

		creds, crderr := common.GetCredentials(&common.GetCredentialsInput{Profile: img.Cloud.Profile, Cloud: img.Cloud.Name})
		if crderr != nil {
			return GetImagesResponse{}, crderr
		}
		// I will establish session so that we can carry out the process in cloud
		session_input := awssess.CreateSessionInput{Region: img.Cloud.Region, KeyId: creds.KeyId, AcessKey: creds.SecretAccess}
		sess := session_input.CreateAwsSession()

		//authorizing to request further
		authinpt := auth.EstablishConnectionInput{Region: img.Cloud.Region, Resource: "ec2", Session: sess}

		getimages := new(image.GetImageInput)
		getimages.GetRaw = img.Cloud.GetRaw
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
