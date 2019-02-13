package getServers

import (
	"fmt"
	auth "github.com/nikhilsbhat/neuron/cloud/aws/interface"
	awscommon "github.com/nikhilsbhat/neuron/cloud/aws/operations/common"
	server "github.com/nikhilsbhat/neuron/cloud/aws/operations/server"
	awssess "github.com/nikhilsbhat/neuron/cloud/aws/sessions"
	common "github.com/nikhilsbhat/neuron/cloudoperations/common"
	support "github.com/nikhilsbhat/neuron/cloudoperations/support"
	db "github.com/nikhilsbhat/neuron/database"
	log "github.com/nikhilsbhat/neuron/logger"
	"strings"
	"sync"
)

// The struct that will return the filtered/unfiltered responses of variuos clouds.
type GetServerResponse struct {
	// Contains filtered/unfiltered response of AWS.
	AwsResponse []server.ServerResponse `json:"AwsResponse,omitempty"`

	// Contains filtered/unfiltered response of Azure.
	AzureResponse string `json:"AzureResponse,omitempty"`

	// Default response if no inputs or matching the values required.
	DefaultResponse string `json:"Response,omitempty"`
}

// Being GetServersDetails, job of him is to fetch the details of servers with the instructions passed to him
// and give back the response who called this.
// Below method will take care of fetching details of
// appropriate user and his cloud profile details which was passed while calling it.
func (serv *GetServersInput) GetServersDetails() (GetServerResponse, error) {

	if status := support.DoesCloudSupports(strings.ToLower(serv.Cloud.Name)); status != true {
		return GetServerResponse{}, fmt.Errorf(common.DefaultCloudResponse + "GetNetworks")
	}

	switch strings.ToLower(serv.Cloud.Name) {
	case "aws":

		creds, crederr := common.GetCredentials(
			&common.GetCredentialsInput{
				Profile: serv.Cloud.Profile,
				Cloud:   serv.Cloud.Name,
			},
		)
		if crederr != nil {
			return GetServerResponse{}, crederr
		}
		// I will establish session so that we can carry out the process in cloud
		session_input := awssess.CreateSessionInput{Region: serv.Cloud.Region, KeyId: creds.KeyId, AcessKey: creds.SecretAccess}
		sess := session_input.CreateAwsSession()

		//authorizing to request further
		authinpt := auth.EstablishConnectionInput{Region: serv.Cloud.Region, Resource: "ec2", Session: sess}
		// I will call CreateServer of interface and get the things done

		if serv.InstanceIds != nil {
			serverin := server.DescribeInstanceInput{}
			serverin.InstanceIds = serv.InstanceIds
			serverin.GetRaw = serv.Cloud.GetRaw
			server_response, serverr := serverin.GetServersDetails(authinpt)
			if serverr != nil {
				return GetServerResponse{}, serverr
			}
			return GetServerResponse{AwsResponse: server_response}, nil
		} else if serv.SubnetIds != nil {
			serverin := server.DescribeInstanceInput{}
			serverin.SubnetIds = serv.SubnetIds
			serverin.GetRaw = serv.Cloud.GetRaw
			server_response, serverr := serverin.GetServersFromSubnet(authinpt)
			if serverr != nil {
				return GetServerResponse{}, serverr
			}
			return GetServerResponse{AwsResponse: server_response}, nil
		} else if serv.VpcIds != nil {
			serverin := server.DescribeInstanceInput{}
			serverin.VpcIds = serv.VpcIds
			serverin.GetRaw = serv.Cloud.GetRaw
			server_response, serverr := serverin.GetServersFromNetwork(authinpt)
			if serverr != nil {
				return GetServerResponse{}, serverr
			}
			return GetServerResponse{AwsResponse: server_response}, nil
		} else {
			serverin := server.DescribeInstanceInput{GetRaw: serv.Cloud.GetRaw}
			server_response, serverr := serverin.GetAllServers(authinpt)
			if serverr != nil {
				return GetServerResponse{}, serverr
			}
			return GetServerResponse{AwsResponse: server_response}, nil
		}
		return GetServerResponse{}, fmt.Errorf("You have not passed valid input to get details of server, the input struct looks like empty")

	case "azure":
		return GetServerResponse{DefaultResponse: common.DefaultAzResponse}, nil
	case "gcp":
		return GetServerResponse{DefaultResponse: common.DefaultGcpResponse}, nil
	case "openstack":
		return GetServerResponse{DefaultResponse: common.DefaultOpResponse}, nil
	default:
		log.Info("")
		log.Error("I feel we are lost in getting details of all the server :S, guess you have entered wrong cloud")
		log.Info("")
		return GetServerResponse{}, fmt.Errorf(common.DefaultCloudResponse + "GetServers")
	}
}

// Being GetAllServers, job of him is to fetch the details of all servers across the cloud
// and give back the response who called this.
// Below method will take care of fetching details of
// appropriate user and his cloud profile details which was passed while calling it.
func (serv *GetServersInput) GetAllServers() ([]GetServerResponse, error) {

	if status := support.DoesCloudSupports(strings.ToLower(serv.Cloud.Name)); status != true {
		return nil, fmt.Errorf(common.DefaultCloudResponse + "GetNetworks")
	}

	//fetchinig credentials from loged-in user to establish the connection with appropriate cloud.
	creds, err := common.GetCredentials(
		&common.GetCredentialsInput{
			Profile: serv.Cloud.Profile,
			Cloud:   serv.Cloud.Name,
		},
	)
	if err != nil {
		return nil, err
	}

	switch strings.ToLower(serv.Cloud.Name) {
	case "aws":

		// I will establish session so that we can carry out the process in cloud
		session_input := awssess.CreateSessionInput{Region: serv.Cloud.Region, KeyId: creds.KeyId, AcessKey: creds.SecretAccess}
		sess := session_input.CreateAwsSession()

		//authorize
		authinpt := auth.EstablishConnectionInput{Region: serv.Cloud.Region, Resource: "ec2", Session: sess}

		// Fetching list of regions to get details  of server across the account
		regionin := awscommon.CommonInput{}
		regions, regerr := regionin.GetRegions(authinpt)
		if regerr != nil {
			return nil, regerr
		}

		reg := make(chan []server.ServerResponse, len(regions.Regions))
		serv.getservers(regions.Regions, reg, creds)

		server_response := make([]GetServerResponse, 0)
		for region_detail := range reg {
			if len(region_detail) != 0 {
				server_response = append(server_response, GetServerResponse{AwsResponse: region_detail})
			}
		}
		return server_response, nil

	case "azure":
		return nil, fmt.Errorf(common.DefaultAzResponse)
	case "gcp":
		return nil, fmt.Errorf(common.DefaultGcpResponse)
	case "openstack":
		return nil, fmt.Errorf(common.DefaultOpResponse)
	default:
		return nil, fmt.Errorf(common.DefaultCloudResponse + "GetAllServers")
	}
}

// this will be called by getallservers, he is the one who gets the details of all the servers,
// and send over a channel.
func (serv *GetServersInput) getservers(regions []string, reg chan []server.ServerResponse, creds db.CloudProfiles) {

	switch strings.ToLower(serv.Cloud.Name) {
	case "aws":
		// I will establish session so that we can carry out the process in cloud
		session_input := awssess.CreateSessionInput{Region: serv.Cloud.Region, KeyId: creds.KeyId, AcessKey: creds.SecretAccess}
		sess := session_input.CreateAwsSession()
		var wg sync.WaitGroup
		wg.Add(len(regions))
		for _, region := range regions {
			go func(region string) {
				defer wg.Done()

				//authorize
				authinpt := auth.EstablishConnectionInput{Region: region, Resource: "ec2", Session: sess}
				serverin := server.DescribeInstanceInput{GetRaw: serv.Cloud.GetRaw}
				server_response, _ := serverin.GetAllServers(authinpt)
				reg <- server_response
			}(region)
		}
		wg.Wait()
		close(reg)

	case "azure":
	case "gcp":
	case "openstack":
	default:
	}
}

func New() *GetServersInput {
	net := &GetServersInput{}
	return net
}
