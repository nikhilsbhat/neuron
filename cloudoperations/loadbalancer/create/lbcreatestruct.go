package createLoadbalancer

type LbCreateInput struct {
	Name              string
	VpcId             string
	SubnetIds         []string `json:"SubnetIds,omitempty"`
	AvailabilityZones []string `json:"AvailabilityZones,omitempty"`
	SecurityGroupIds  []string `json:"SecurityGroupIds,omitempty"`
	Scheme            string
	Type              string //required only if the LB protocol is HTTPS else can be initiazed with dummy value
	SslCert           string `json:"SslCert,omitempty"`
	SslPolicy         string `json:"SslPolicy,omitempty"`
	LbPort            int64  //required ex: 8080 or 80 etc
	InstPort          int64
	Lbproto           string //required ex: HTTPS, HTTP
	Instproto         string
	HttpCode          string `json:"HttpCode,omitempty"`
	HealthPath        string `json:"HealthPath,omitempty"`
	IpAddressType     string `json:"IpAddressType,omitempty"`
	Cloud             string
	Region            string
	Profile           string
	GetRaw            bool
}

//Nothing much from this file. This file contains only the structs for loadbalance/create
