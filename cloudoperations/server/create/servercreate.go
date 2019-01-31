package serverCreate

import (
	"fmt"
	auth "neuron/cloud/aws/interface"
	server "neuron/cloud/aws/operations/server"
	awssess "neuron/cloud/aws/sessions"
	common "neuron/cloudoperations/common"
	support "neuron/cloudoperations/support"
	//log "neuron/logger"
	"reflect"
	"strings"
)

// The struct that will return the filtered/unfiltered responses of variuos clouds.
type ServerCreateResponse struct {
	// Contains filtered/unfiltered response of AWS.
	AwsResponse []server.ServerResponse `json:"AwsResponse,omitempty"`

	// Contains filtered/unfiltered response of Azure.
	AzureResponse string `json:"AzureResponse,omitempty"`

	// Default response if no inputs or matching the values required.
	DefaultResponse interface{} `json:"DefaultResponse,omitempty"`
}

// Being CreateServer, job of him is to create the server with the requirement passed to him
// and give back the response who called this.
// Below method will take care of fetching details of
// appropriate user and his cloud profile details which was passed while calling it.
func (serv ServerCreateInput) CreateServer() (ServerCreateResponse, error) {

	if status := support.DoesCloudSupports(strings.ToLower(serv.Cloud)); status != true {
		return ServerCreateResponse{}, fmt.Errorf(common.DefaultCloudResponse + "GetNetworks")
	}
	switch strings.ToLower(serv.Cloud) {
	case "aws":

		creds, crederr := common.GetCredentials(
			&common.GetCredentialsInput{
				Profile: serv.Profile,
				Cloud:   serv.Cloud,
			},
		)
		if crederr != nil {
			return ServerCreateResponse{}, crederr
		}
		// I will establish session so that we can carry out the process in cloud
		session_input := awssess.CreateSessionInput{Region: serv.Region, KeyId: creds.KeyId, AcessKey: creds.SecretAccess}
		sess := session_input.CreateAwsSession()

		//authorizing to request further
		authinpt := auth.EstablishConnectionInput{Region: serv.Region, Resource: "ec2", Session: sess}

		// I will call CreateServer of interface and get the things done

		serverin := server.CreateServerInput{}
		serverin.InstanceName = serv.InstanceName
		serverin.ImageId = serv.ImageId
		serverin.InstanceType = serv.Flavor
		serverin.KeyName = serv.KeyName
		serverin.MaxCount = serv.Count
		serverin.SubnetId = serv.SubnetId
		serverin.UserData = serv.UserData
		serverin.AssignPubIp = serv.AssignPubIp
		serverin.GetRaw = serv.GetRaw
		response, err := serverin.CreateServer(authinpt)
		if err != nil {
			return ServerCreateResponse{}, err
		}
		return ServerCreateResponse{AwsResponse: response}, nil

	case "azure":
		return ServerCreateResponse{}, fmt.Errorf(common.DefaultAzResponse)
	case "gcp":
		return ServerCreateResponse{}, fmt.Errorf(common.DefaultGcpResponse)
	case "openstack":
		return ServerCreateResponse{}, fmt.Errorf(common.DefaultOpResponse)
	default:

		return ServerCreateResponse{}, fmt.Errorf(common.DefaultCloudResponse + "CreateServer")
	}
}

func CreateServerMock() (ServerCreateResponse, error) {

	input := new(ServerCreateInput)
	defaultval := reflect.Indirect(reflect.ValueOf(input))

	defaults := make(map[string]interface{})
	for i := 0; i < defaultval.NumField(); i++ {
		defaults[defaultval.Type().Field(i).Name] = defaultval.Type().Field(i).Type
	}

	return ServerCreateResponse{DefaultResponse: defaults}, nil
}

func New() *ServerCreateInput {
	net := &ServerCreateInput{}
	return net
}
