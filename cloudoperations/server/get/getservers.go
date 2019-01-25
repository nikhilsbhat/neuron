package getServers

import (
	"fmt"
	auth "neuron/cloud/aws/interface"
	awscommon "neuron/cloud/aws/operations/common"
	server "neuron/cloud/aws/operations/server"
	awssess "neuron/cloud/aws/sessions"
	common "neuron/cloudoperations/common"
	db "neuron/database"
	log "neuron/logger"
	"strings"
	"sync"
)

type GetServerResponse struct {
	AwsResponse     []server.ServerResponse `json:"AwsResponse,omitempty"`
	AzureResponse   string                  `json:"AzureResponse,omitempty"`
	DefaultResponse string                  `json:"Response,omitempty"`
}

// being getserversdetails my job is to call appropriate function to get serversdetails and give back the response who called me.
func (serv *GetServersInput) GetServersDetails() (GetServerResponse, error) {

	switch strings.ToLower(serv.Cloud) {
	case "aws":

		creds, err := common.GetCredentials(&common.GetCredentialsInput{Profile: serv.Profile, Cloud: serv.Cloud})
		if err != nil {
			return GetServerResponse{}, err
		}
		// I will establish session so that we can carry out the process in cloud
		session_input := awssess.CreateSessionInput{Region: serv.Region, KeyId: creds.KeyId, AcessKey: creds.SecretAccess}
		sess := session_input.CreateAwsSession()

		//authorizing to request further
		authinpt := auth.EstablishConnectionInput{Region: serv.Region, Resource: "ec2", Session: sess}
		// I will call CreateServer of interface and get the things done

		if serv.InstanceIds != nil {
			serverin := server.DescribeInstanceInput{InstanceIds: serv.InstanceIds, GetRaw: serv.GetRaw}
			server_response, serverr := serverin.GetServersDetails(authinpt)
			if serverr != nil {
				return GetServerResponse{}, serverr
			}
			return GetServerResponse{AwsResponse: server_response}, nil
		} else if serv.SubnetIds != nil {
			serverin := server.DescribeInstanceInput{SubnetIds: serv.SubnetIds, GetRaw: serv.GetRaw}
			server_response, serverr := serverin.GetServersFromSubnet(authinpt)
			if serverr != nil {
				return GetServerResponse{}, serverr
			}
			return GetServerResponse{AwsResponse: server_response}, nil
		} else if serv.VpcIds != nil {
			serverin := server.DescribeInstanceInput{VpcIds: serv.VpcIds, GetRaw: serv.GetRaw}
			server_response, serverr := serverin.GetServersFromNetwork(authinpt)
			if serverr != nil {
				return GetServerResponse{}, serverr
			}
			return GetServerResponse{AwsResponse: server_response}, nil
		} else {
			serverin := server.DescribeInstanceInput{GetRaw: serv.GetRaw}
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

func (serv *GetServersInput) GetAllServers() ([]GetServerResponse, error) {

	//fetchinig credentials from loged-in user to establish the connection with appropriate cloud.
	creds, err := common.GetCredentials(&common.GetCredentialsInput{Profile: serv.Profile, Cloud: serv.Cloud})
	if err != nil {
		return nil, err
	}

	switch strings.ToLower(serv.Cloud) {
	case "aws":

		// I will establish session so that we can carry out the process in cloud
		session_input := awssess.CreateSessionInput{Region: serv.Region, KeyId: creds.KeyId, AcessKey: creds.SecretAccess}
		sess := session_input.CreateAwsSession()

		//authorize
		authinpt := auth.EstablishConnectionInput{Region: serv.Region, Resource: "ec2", Session: sess}

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
		log.Info("")
		log.Error("I feel we are lost in getting details of all the server :S, guess you have entered wrong cloud")
		log.Info("")
		return nil, fmt.Errorf(common.DefaultCloudResponse + "GetAllServers")
	}
}

func (serv *GetServersInput) getservers(regions []string, reg chan []server.ServerResponse, creds db.CloudProfiles) {

	switch strings.ToLower(serv.Cloud) {
	case "aws":
		// I will establish session so that we can carry out the process in cloud
		session_input := awssess.CreateSessionInput{Region: serv.Region, KeyId: creds.KeyId, AcessKey: creds.SecretAccess}
		sess := session_input.CreateAwsSession()
		var wg sync.WaitGroup
		wg.Add(len(regions))
		for _, region := range regions {
			go func(region string) {
				defer wg.Done()

				//authorize
				authinpt := auth.EstablishConnectionInput{Region: region, Resource: "ec2", Session: sess}
				serverin := server.DescribeInstanceInput{GetRaw: serv.GetRaw}
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
