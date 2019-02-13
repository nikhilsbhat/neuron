package networkGet

import (
	"fmt"
	auth "github.com/nikhilsbhat/neuron/cloud/aws/interface"
	network "github.com/nikhilsbhat/neuron/cloud/aws/operations/network"
	awssess "github.com/nikhilsbhat/neuron/cloud/aws/sessions"
	common "github.com/nikhilsbhat/neuron/cloudoperations/common"
	support "github.com/nikhilsbhat/neuron/cloudoperations/support"
	//log "neuron/logger"
	"strings"
)

// The struct that will return the filtered/unfiltered responses of variuos clouds.
type GetSubnetsResponse struct {
	// Contains filtered/unfiltered response of AWS.
	AwsResponse network.NetworkResponse `json:"AwsResponse,omitempty"`

	// Contains filtered/unfiltered response of Azure.
	AzureResponse string `json:"AzureResponse,omitempty"`

	// Default response if no inputs or matching the values required.
	DefaultResponse string `json:"DefaultResponse,omitempty"`
}

// Being GetSubnets, job of him is to fetch the details of subnets entered
// and give back the response who called this.
// Below method will take care of fetching details of
// appropriate user and his cloud profile details which was passed while calling it.
func (sub GetNetworksInput) GetSubnets() (GetSubnetsResponse, error) {

	if status := support.DoesCloudSupports(strings.ToLower(sub.Cloud.Name)); status != true {
		return GetSubnetsResponse{}, fmt.Errorf(common.DefaultCloudResponse + "GetSubnets")
	}

	switch strings.ToLower(sub.Cloud.Name) {
	case "aws":

		creds, err := common.GetCredentials(
			&common.GetCredentialsInput{
				Profile: sub.Cloud.Profile,
				Cloud:   sub.Cloud.Name,
			},
		)

		if err != nil {
			return GetSubnetsResponse{}, err
		}
		// I will establish session so that we can carry out the process in cloud
		session_input := awssess.CreateSessionInput{Region: sub.Cloud.Region, KeyId: creds.KeyId, AcessKey: creds.SecretAccess}
		sess := session_input.CreateAwsSession()

		//authorizing to request further
		authinpt := auth.EstablishConnectionInput{Region: sub.Cloud.Region, Resource: "ec2", Session: sess}

		// I will call getsubnets and get the things done
		networkin := new(network.GetNetworksInput)
		networkin.GetRaw = sub.Cloud.GetRaw
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
		return GetSubnetsResponse{}, fmt.Errorf(common.DefaultCloudResponse + "GetSubnets")
	}
}
