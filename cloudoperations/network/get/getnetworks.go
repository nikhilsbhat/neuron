package networkGet

import (
	"fmt"
	auth "neuron/cloud/aws/interface"
	awscommon "neuron/cloud/aws/operations/common"
	network "neuron/cloud/aws/operations/network"
	awssess "neuron/cloud/aws/sessions"
	common "neuron/cloudoperations/common"
	support "neuron/cloudoperations/support"
	//log "neuron/logger"
	"strings"
)

// The struct that will return the filtered/unfiltered responses of variuos clouds.
type GetNetworksResponse struct {
	// Contains filtered/unfiltered response of AWS.
	AwsResponse []network.NetworkResponse `json:"AwsResponse,omitempty"`

	// Contains filtered/unfiltered response of Azure.
	AzureResponse string `json:"AzureResponse,omitempty"`

	// Default response if no inputs or matching the values required.
	DefaultResponse string `json:"DefaultResponse,omitempty"`
}

// Being GetNetwork, job of him is to fetch the details of networks entered
// and give back the response who called this.
// Below method will take care of fetching details of
// appropriate user and his cloud profile details which was passed while calling it.
func (net *GetNetworksInput) GetNetworks() (GetNetworksResponse, error) {

	if status := support.DoesCloudSupports(strings.ToLower(net.Cloud.Name)); status != true {
		return GetNetworksResponse{}, fmt.Errorf(common.DefaultCloudResponse + "GetNetworks")
	}

	switch strings.ToLower(net.Cloud.Name) {
	case "aws":

		creds, err := common.GetCredentials(
			&common.GetCredentialsInput{
				Profile: net.Cloud.Profile,
				Cloud:   net.Cloud.Name,
			},
		)

		if err != nil {
			return GetNetworksResponse{}, err
		}
		// I will establish session so that we can carry out the process in cloud
		session_input := awssess.CreateSessionInput{Region: net.Cloud.Region, KeyId: creds.KeyId, AcessKey: creds.SecretAccess}
		sess := session_input.CreateAwsSession()

		//authorizing to request further
		authinpt := auth.EstablishConnectionInput{Region: net.Cloud.Region, Resource: "ec2", Session: sess}

		// Fetching all the networks accross cloud aws
		networkin := network.GetNetworksInput{}
		networkin.VpcIds = net.VpcIds
		networkin.GetRaw = net.Cloud.GetRaw
		response, net_err := networkin.GetNetwork(authinpt)
		if net_err != nil {
			return GetNetworksResponse{}, net_err
		}
		return GetNetworksResponse{AwsResponse: response}, nil

	case "azure":
		return GetNetworksResponse{}, fmt.Errorf(common.DefaultAzResponse)
	case "gcp":
		return GetNetworksResponse{}, fmt.Errorf(common.DefaultGcpResponse)
	case "openstack":
		return GetNetworksResponse{}, fmt.Errorf(common.DefaultOpResponse)
	default:
		return GetNetworksResponse{}, fmt.Errorf(common.DefaultCloudResponse + "GetNetworks")
	}
}

// Being GetAllNetworks, job of him is to fetch the details of all networks entered and
// give back the response who called this.
// Below method will take care of fetching details of
// appropriate user and his cloud profile details which was passed while calling it.
func (net GetNetworksInput) GetAllNetworks() ([]GetNetworksResponse, error) {

	if status := support.DoesCloudSupports(strings.ToLower(net.Cloud.Name)); status != true {
		return nil, fmt.Errorf(common.DefaultCloudResponse + "GetAllNetworks")
	}

	switch strings.ToLower(net.Cloud.Name) {
	case "aws":

		creds, err := common.GetCredentials(&common.GetCredentialsInput{Profile: net.Cloud.Profile, Cloud: net.Cloud.Name})
		if err != nil {
			return nil, err
		}
		// I will establish session so that we can carry out the process in cloud
		session_input := awssess.CreateSessionInput{Region: net.Cloud.Region, KeyId: creds.KeyId, AcessKey: creds.SecretAccess}
		sess := session_input.CreateAwsSession()

		//authorizing to request further
		authinpt := auth.EstablishConnectionInput{Region: net.Cloud.Region, Resource: "ec2", Session: sess}

		// I will call GetAllNetworks of interface and get the things done
		// Fetching all the regions from the cloud aws
		regionin := awscommon.CommonInput{}
		regions, regerr := regionin.GetRegions(authinpt)
		if regerr != nil {
			return nil, regerr
		}
		// Fetching all the networks accross all the regions of cloud aws
		/*reg := make(chan []DengineAwsInterface.NetworkResponse, len(get_region_response.Regions))

		getnetworkdetails_input := GetAllNetworksInput{net.Cloud, net.Region}
		getnetworkdetails_input.getnetworkdetails(get_region_response.Regions, reg)
		for region_detail := range reg {
				get_all_network_response = append(get_all_network_response, GetAllNetworksResponse{AwsResponse: region_detail})
		}*/
		network_response := make([]GetNetworksResponse, 0)
		for _, region := range regions.Regions {
			//authorizing to request further
			authinpt := auth.EstablishConnectionInput{Region: region, Resource: "ec2", Session: sess}

			networkin := network.GetNetworksInput{GetRaw: net.Cloud.GetRaw}
			response, net_err := networkin.GetAllNetworks(authinpt)
			if net_err != nil {
				return nil, net_err
			}
			network_response = append(network_response, GetNetworksResponse{AwsResponse: response})
		}
		return network_response, nil

	case "azure":
		return nil, fmt.Errorf(common.DefaultAzResponse)
	case "gcp":
		return nil, fmt.Errorf(common.DefaultGcpResponse)
	case "openstack":
		return nil, fmt.Errorf(common.DefaultOpResponse)
	default:
		return nil, fmt.Errorf(common.DefaultCloudResponse + "GetAllNetworks")
	}
}

func New() *GetNetworksInput {
	net := &GetNetworksInput{}
	return net
}
