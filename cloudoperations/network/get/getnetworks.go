package networkGet

import (
	"fmt"
	auth "neuron/cloud/aws/interface"
	awscommon "neuron/cloud/aws/operations/common"
	network "neuron/cloud/aws/operations/network"
	awssess "neuron/cloud/aws/sessions"
	common "neuron/cloudoperations/common"
	log "neuron/logger"
	"strings"
)

type GetNetworksResponse struct {
	AwsResponse     []network.NetworkResponse `json:"AwsResponse,omitempty"`
	AzureResponse   string                    `json:"AzureResponse,omitempty"`
	DefaultResponse string                    `json:"DefaultResponse,omitempty"`
}

// being GetNetwork my job is to call appropriate function of operations and give back the response who called me
func (net *GetNetworksInput) GetNetworks() (GetNetworksResponse, error) {

	switch strings.ToLower(net.Cloud) {
	case "aws":

		creds, err := common.GetCredentials(&common.GetCredentialsInput{Profile: net.Profile, Cloud: net.Cloud})
		if err != nil {
			return GetNetworksResponse{}, err
		}
		// I will establish session so that we can carry out the process in cloud
		session_input := awssess.CreateSessionInput{Region: net.Region, KeyId: creds.KeyId, AcessKey: creds.SecretAccess}
		sess := session_input.CreateAwsSession()

		//authorizing to request further
		authinpt := auth.EstablishConnectionInput{Region: net.Region, Resource: "ec2", Session: sess}

		// Fetching all the networks accross cloud aws
		networkin := network.GetNetworksInput{}
		networkin.VpcIds = net.VpcIds
		networkin.GetRaw = net.GetRaw
		response, net_err := networkin.GetNetwork(authinpt)
		if net_err != nil {
			return GetNetworksResponse{}, nil
		}
		return GetNetworksResponse{AwsResponse: response}, nil

	case "azure":
		return GetNetworksResponse{}, fmt.Errorf(common.DefaultAzResponse)
	case "gcp":
		return GetNetworksResponse{}, fmt.Errorf(common.DefaultGcpResponse)
	case "openstack":
		return GetNetworksResponse{}, fmt.Errorf(common.DefaultOpResponse)
	default:
		log.Info("")
		log.Error("I feel we are lost in getting details of networks :S")
		log.Info("")
		return GetNetworksResponse{}, fmt.Errorf(common.DefaultCloudResponse + "GetNetworks")
	}
}

// being GetAllNetworks my job is to call appropriate function of operations and give back the response who called me
func (net GetNetworksInput) GetAllNetworks() ([]GetNetworksResponse, error) {

	switch strings.ToLower(net.Cloud) {
	case "aws":

		creds, err := common.GetCredentials(&common.GetCredentialsInput{Profile: net.Profile, Cloud: net.Cloud})
		if err != nil {
			return nil, err
		}
		// I will establish session so that we can carry out the process in cloud
		session_input := awssess.CreateSessionInput{Region: net.Region, KeyId: creds.KeyId, AcessKey: creds.SecretAccess}
		sess := session_input.CreateAwsSession()

		//authorizing to request further
		authinpt := auth.EstablishConnectionInput{Region: net.Region, Resource: "ec2", Session: sess}

		// I will call GetAllNetworks of interface and get the things done
		// Fetching all the regions from the cloud aws
		regionin := awscommon.CommonInput{}
		regions, regerr := regionin.GetRegions(authinpt)
		if regerr != nil {
			return nil, regerr
		}
		// Fetching all the servers accross all the regions of cloud aws
		/*reg := make(chan []DengineAwsInterface.NetworkResponse, len(get_region_response.Regions))

		getnetworkdetails_input := GetAllNetworksInput{net.Cloud, net.Region}
		getnetworkdetails_input.getnetworkdetails(get_region_response.Regions, reg)
		for region_detail := range reg {
				get_all_network_response = append(get_all_network_response, GetAllNetworksResponse{AwsResponse: region_detail})
		}*/
		network_response := make([]GetNetworksResponse, 0)
		for _, region := range regions {
			//authorizing to request further
			authinpt := auth.EstablishConnectionInput{Region: region, Resource: "ec2", Session: sess}

			networkin := network.GetNetworksInput{GetRaw: net.GetRaw}
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
		log.Info("")
		log.Error("I feel we are lost in getting details of all the networks :S")
		log.Info("")
		return nil, fmt.Errorf(common.DefaultCloudResponse + "GetAllNetworks")
	}
}
