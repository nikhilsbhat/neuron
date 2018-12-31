package networkGet

import (
	"fmt"
	auth "neuron/cloud/aws/interface"
	network "neuron/cloud/aws/operations/network"
	awssess "neuron/cloud/aws/sessions"
	common "neuron/cloudoperations/common"
	log "neuron/logger"
	"strings"
)

type GetSubnetsResponse struct {
	AwsResponse     network.NetworkResponse `json:"AwsResponse,omitempty"`
	AzureResponse   string                  `json:"AzureResponse,omitempty"`
	DefaultResponse string                  `json:"DefaultResponse,omitempty"`
}

// being GetSubnets my job is to call appropriate function of operations and give back the response who called me
func (sub GetNetworksInput) GetSubnets() (GetSubnetsResponse, error) {

	switch strings.ToLower(sub.Cloud) {
	case "aws":

		creds, err := common.GetCredentials(&common.GetCredentialsInput{Profile: sub.Profile, Cloud: sub.Cloud})
		if err != nil {
			return GetSubnetsResponse{}, err
		}
		// I will establish session so that we can carry out the process in cloud
		session_input := awssess.CreateSessionInput{Region: sub.Region, KeyId: creds.KeyId, AcessKey: creds.SecretAccess}
		sess := session_input.CreateAwsSession()

		//authorizing to request further
		authinpt := auth.EstablishConnectionInput{Region: sub.Region, Resource: "ec2", Session: sess}

		// I will call getsubnets and get the things done
		networkin := new(network.GetNetworksInput)
		networkin.GetRaw = sub.GetRaw
		if sub.SubnetIds != nil {
			networkin.SubnetIds = sub.SubnetIds
			response, get_sub_err := networkin.GetSubnets(authinpt)
			if get_sub_err != nil {
				return GetSubnetsResponse{}, get_sub_err
			}
			return GetSubnetsResponse{AwsResponse: response}, nil
		} else if sub.VpcIds != nil {
			networkin.VpcIds = sub.VpcIds
			response, get_sub_err := networkin.GetSubnetsFromVpc(authinpt)
			if get_sub_err != nil {
				return GetSubnetsResponse{}, get_sub_err
			}
			return GetSubnetsResponse{AwsResponse: response}, nil
		} else {
			return GetSubnetsResponse{}, fmt.Errorf("You have not passed valid input to get details of server, the input struct looks like empty")
		}

	case "azure":
		return GetSubnetsResponse{}, fmt.Errorf(common.DefaultAzResponse)
	case "gcp":
		return GetSubnetsResponse{}, fmt.Errorf(common.DefaultGcpResponse)
	case "openstack":
		return GetSubnetsResponse{}, fmt.Errorf(common.DefaultOpResponse)
	default:
		log.Info("")
		log.Error("I feel we are lost in getting list of SUBNETS, was unable to find the appropriate cloud :S")
		log.Error("You have entered wrong cloud name else misspelled the name of it, please check it before passing")
		log.Info("")
		return GetSubnetsResponse{}, fmt.Errorf(common.DefaultCloudResponse + "GetSubnets")
	}
}
