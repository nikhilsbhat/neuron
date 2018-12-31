package getloadbalancer

import (
	"fmt"
	auth "neuron/cloud/aws/interface"
	loadbalance "neuron/cloud/aws/operations/loadbalancer"
	awssess "neuron/cloud/aws/sessions"
	common "neuron/cloudoperations/common"
	log "neuron/logger"
	"strings"
)

type GetLoadbalancerResponse struct {
	AwsResponse     []loadbalance.LoadBalanceResponse `json:"AwsResponse,omitempty"`
	AzureResponse   string                            `json:"AzureResponse,omitempty"`
	DefaultResponse string                            `json:"Response,omitempty"`
}

// being get_all_loadbalancers my job is to gather info on all the loadbalancer and give back the response who called me
func (lb *GetLoadbalancerInput) GetLoadbalancers() (GetLoadbalancerResponse, error) {

	switch strings.ToLower(lb.Cloud) {
	case "aws":

		creds, err := common.GetCredentials(&common.GetCredentialsInput{Profile: lb.Profile, Cloud: lb.Cloud})
		if err != nil {
			return GetLoadbalancerResponse{}, err
		}
		// I will establish session so that we can carry out the process in cloud
		session_input := awssess.CreateSessionInput{Region: lb.Region, KeyId: creds.KeyId, AcessKey: creds.SecretAccess}
		sess := session_input.CreateAwsSession()

		//authorizing to request further
		authinpt := new(auth.EstablishConnectionInput)
		authinpt.Region = lb.Region
		authinpt.Session = sess
		switch strings.ToLower(lb.Type) {
		case "classic":
			authinpt.Resource = "elb"
		case "application":
			authinpt.Resource = "elb2"
		}

		lbin := new(loadbalance.GetLoadbalancerInput)
		lbin.GetRaw = lb.GetRaw
		lbin.LbNames = lb.LbNames
		lbin.LbArns = lb.LbArns
		lbin.Type = lb.Type
		response, lberr := lbin.Getloadbalancers(*authinpt)
		if lberr != nil {
			return GetLoadbalancerResponse{}, lberr
		}
		return GetLoadbalancerResponse{AwsResponse: response}, nil

	case "azure":
		return GetLoadbalancerResponse{}, fmt.Errorf(common.DefaultAzResponse)
	case "gcp":
		return GetLoadbalancerResponse{}, fmt.Errorf(common.DefaultGcpResponse)
	case "openstack":
		return GetLoadbalancerResponse{}, fmt.Errorf(common.DefaultOpResponse)
	default:
		log.Info("")
		log.Error(common.DefaultCloudResponse + "GetLoadbalancers")
		log.Info("")
		return GetLoadbalancerResponse{}, fmt.Errorf(common.DefaultCloudResponse + "GetLoadbalancers")
	}
}

func (lb *GetLoadbalancerInput) GetAllLoadbalancer() (GetLoadbalancerResponse, error) {

	switch strings.ToLower(lb.Cloud) {
	case "aws":

		creds, err := common.GetCredentials(&common.GetCredentialsInput{Profile: lb.Profile, Cloud: lb.Cloud})
		if err != nil {
			return GetLoadbalancerResponse{}, err
		}
		// I will establish session so that we can carry out the process in cloud
		session_input := awssess.CreateSessionInput{Region: lb.Region, KeyId: creds.KeyId, AcessKey: creds.SecretAccess}
		sess := session_input.CreateAwsSession()

		//authorizing to request further
		authinpt := new(auth.EstablishConnectionInput)
		authinpt.Region = lb.Region
		authinpt.Session = sess
		authinpt.Resource = "elb12"
		lbin := new(loadbalance.GetLoadbalancerInput)
		lbin.GetRaw = lb.GetRaw

		switch strings.ToLower(lb.Type) {
		case "classic":
			response, lberr := lbin.GetAllClassicLb(*authinpt)
			if lberr != nil {
				return GetLoadbalancerResponse{}, lberr
			}
			return GetLoadbalancerResponse{AwsResponse: response}, nil
		case "application":
			response, lberr := lbin.GetAllApplicationLb(*authinpt)
			if lberr != nil {
				return GetLoadbalancerResponse{}, lberr
			}
			return GetLoadbalancerResponse{AwsResponse: response}, nil
		case "":
			response, lberr := lbin.GetAllLoadbalancer(*authinpt)
			if lberr != nil {
				return GetLoadbalancerResponse{}, lberr
			}
			return GetLoadbalancerResponse{AwsResponse: response}, nil
		default:
			return GetLoadbalancerResponse{}, fmt.Errorf("The loadbalancer type you entered is unknown to me")
		}

	case "azure":
		return GetLoadbalancerResponse{}, fmt.Errorf(common.DefaultAzResponse)
	case "gcp":
		return GetLoadbalancerResponse{}, fmt.Errorf(common.DefaultGcpResponse)
	case "openstack":
		return GetLoadbalancerResponse{}, fmt.Errorf(common.DefaultOpResponse)
	default:
		log.Info("")
		log.Error(common.DefaultCloudResponse + "GetLoadbalancers")
		log.Info("")
		return GetLoadbalancerResponse{}, fmt.Errorf(common.DefaultCloudResponse + "GetLoadbalancers")
	}
}
