package createLoadbalancer

import (
	"fmt"
	auth "neuron/cloud/aws/interface"
	loadbalance "neuron/cloud/aws/operations/loadbalancer"
	awssess "neuron/cloud/aws/sessions"
	common "neuron/cloudoperations/common"
	support "neuron/cloudoperations/support"
	log "neuron/logger"
	"strings"
)

// The struct that will return the filtered/unfiltered responses of variuos clouds.
type LoadBalanceResponse struct {
	// Contains filtered/unfiltered response of AWS.
	AwsResponse loadbalance.LoadBalanceResponse `json:"AwsResponse,omitempty"`

	// Contains filtered/unfiltered response of Azure.
	AzureResponse string `json:"AzureResponse,omitempty"`

	// Default response if no inputs or matching the values required.
	DefaultResponse string `json:"Response,omitempty"`
}

// Being CreateLoadBalancer, job of him is to create loadbalancer of choice passed
// and give back the response who called this.
// Below method will take care of fetching details of
// appropriate user and his cloud profile details which was passed while calling it.
func (lb *LbCreateInput) CreateLoadBalancer() (LoadBalanceResponse, error) {

	if status := support.DoesCloudSupports(strings.ToLower(lb.Cloud.Name)); status != true {
		return LoadBalanceResponse{}, fmt.Errorf(common.DefaultCloudResponse + "GetNetworks")
	}

	switch strings.ToLower(lb.Cloud.Name) {
	case "aws":

		creds, err := common.GetCredentials(
			&common.GetCredentialsInput{
				Profile: lb.Cloud.Profile,
				Cloud:   lb.Cloud.Name,
			},
		)
		if err != nil {
			return LoadBalanceResponse{}, err
		}
		// I will establish session so that we can carry out the process in cloud
		session_input := awssess.CreateSessionInput{Region: lb.Cloud.Region, KeyId: creds.KeyId, AcessKey: creds.SecretAccess}
		sess := session_input.CreateAwsSession()

		//authorizing to request further
		authinpt := new(auth.EstablishConnectionInput)
		authinpt.Region = lb.Cloud.Region
		authinpt.Session = sess
		switch strings.ToLower(lb.Type) {
		case "classic":
			authinpt.Resource = "elb"
		case "application":
			authinpt.Resource = "elb2"
		}

		lbin := new(loadbalance.LoadBalanceCreateInput)
		lbin.GetRaw = lb.Cloud.GetRaw
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

func New() *LbCreateInput {
	net := &LbCreateInput{}
	return net
}
