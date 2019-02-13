package neuronaws

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/elb"
	"github.com/aws/aws-sdk-go/service/elbv2"
	"strings"
)

type EstablishedSession struct {
	//This will hold the session for all ec2 resource
	Ec2 *ec2.EC2 `json:"Ec2,omitempty"`
	//This will hold the session for all elb(loadbalancer) resource
	Elb *elb.ELB `json:"Elb,omitempty"`
	//This will hold the session for all elb2(loadbalancer version2) resource
	Elb2 *elbv2.ELBV2 `json:"Elb2,omitempty"`
}

type EstablishConnectionInput struct {
	Region   string
	Resource string
	Session  *session.Session
}

func (con *EstablishConnectionInput) EstablishConnection() (EstablishedSession, error) {

	sesscopy := (con.Session).Copy(&aws.Config{Region: aws.String(con.Region)})

	switch strings.ToLower(con.Resource) {
	case "ec2":
		return EstablishedSession{Ec2: ec2.New(sesscopy)}, nil
	case "elb":
		return EstablishedSession{Ec2: ec2.New(sesscopy), Elb: elb.New(sesscopy)}, nil
	case "elb2":
		return EstablishedSession{Ec2: ec2.New(sesscopy), Elb2: elbv2.New(sesscopy)}, nil
	case "elb12":
		return EstablishedSession{Ec2: ec2.New(sesscopy), Elb: elb.New(sesscopy), Elb2: elbv2.New(sesscopy)}, nil
	default:
		return EstablishedSession{}, fmt.Errorf("Session not established..!!. Unknown resource type, either we don't support this resource or entered resource does not exists")
	}
}
