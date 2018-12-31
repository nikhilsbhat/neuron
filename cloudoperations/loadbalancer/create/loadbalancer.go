package createLoadbalancer

import (
	"fmt"
	auth "neuron/cloud/aws/interface"
	loadbalance "neuron/cloud/aws/operations/loadbalancer"
	awssess "neuron/cloud/aws/sessions"
	common "neuron/cloudoperations/common"
	log "neuron/logger"
	"strings"
)

type LoadBalanceResponse struct {
	AwsResponse     loadbalance.LoadBalanceResponse `json:"AwsResponse,omitempty"`
	AzureResponse   string                          `json:"AzureResponse,omitempty"`
	DefaultResponse string                          `json:"Response,omitempty"`
}

// being create_loadbalancer my job is to create required loadbalancer and give back the response who called me
func (lb *LbCreateInput) CreateLoadBalancer() (LoadBalanceResponse,error) {

	switch strings.ToLower(lb.Cloud) {
	case "aws":

		creds, err := common.GetCredentials(&common.GetCredentialsInput{Profile: lb.Profile, Cloud: lb.Cloud})
		if err != nil {
			return LoadBalanceResponse{}, err
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

		lbin := new(loadbalance.LoadBalanceCreateInput)
		lbin.GetRaw = lb.GetRaw
		lbin.Name = lb.Name
		lbin.VpcId = lb.VpcId
		lbin.SubnetIds = lb.SubnetIds
		lbin.AvailabilityZones = lb.AvailabilityZones
		lbin.SecurityGroupIds = lb.SecurityGroupIds
		lbin.Scheme = lb.Scheme
		lbin.Type = lb.Type
		lbin.LbPort = lb.LbPort
		lbin.InstPort = lb.InstPort
		lbin.Lbproto = lb.Lbproto
		lbin.Instproto = lb.Instproto
		lbin.HttpCode = lb.HttpCode
		lbin.HealthPath = lb.HealthPath
		lbin.SslCert = lb.SslCert
		lbin.SslPolicy = lb.SslPolicy
		lbin.IpAddressType = lb.IpAddressType
		response, lberr := lbin.CreateLoadBalancer(*authinpt)
		if lberr != nil {
			return LoadBalanceResponse{}, lberr
		}
		return LoadBalanceResponse{AwsResponse: response}, nil

	case "azure":
		return LoadBalanceResponse{}, fmt.Errorf(common.DefaultAzResponse)
	case "gcp":
		return LoadBalanceResponse{}, fmt.Errorf(common.DefaultGcpResponse)
	case "openstack":
		return LoadBalanceResponse{}, fmt.Errorf(common.DefaultOpResponse)
	default:
		log.Info("")
		log.Error(common.DefaultCloudResponse + "CreateLoadBalancer")
		log.Info("")
		return LoadBalanceResponse{}, fmt.Errorf(common.DefaultCloudResponse + "CreateLoadBalancer")
	}
}
