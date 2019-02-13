package neuronaws

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/elb"
	"github.com/aws/aws-sdk-go/service/elbv2"
	err "neuron/error"
	"strconv"
)

type LoadBalanceCreateInput struct {
	Name              string   `json:"Name,omitempty"`
	VpcId             string   `json:"VpcId,omitempty"`
	Subnets           []string `json:"Subnets,omitempty"`
	AvailabilityZones []string
	SecurityGroups    []string `json:"SecurityGroups,omitempty"`
	Scheme            string   `json:"Scheme,omitempty"`
	Type              string   `json:"Type,omitempty"`
	SslCert           string   `json:"SslCert,omitempty"` //required only if the LB protocol is HTTPS else can be initiazed with dummy value
	SslPolicy         string   `json:"SslPolicy,omitempty"`
	LbPort            int64    `json:"LbPort,omitempty"` //required ex: 8080 or 80 etc
	InstPort          int64    `json:"InstPort,omitempty"`
	Lbproto           string   `json:"Lbproto,omitempty"` //required ex: HTTPS, HTTP
	Instproto         string   `json:"Instproto,omitempty"`
	HttpCode          string   `json:"HttpCode,omitempty"`
	HealthPath        string   `json:"HealthPath,omitempty"`
	IpAddressType     string   `json:"IpAddressType,omitempty"`
	TargetArn         string   `json:"TargetArn,omitempty"`
	LbArn             string   `json:"LbArn,omitempty"`
}

type LoadBalanceResponse struct {
	Name            string                `json:"Name,omitempty""`
	Type            string                `json:"Type,omitempty""`
	LbDns           string                `json:"LbDns,omitempty""`
	LbArn           string                `json:"LbArn,omitempty"`
	TargetArn       string                `json:"TargetArn,omitempty"`
	ListnerArn      string                `json:"ListnerArn,omitempty"`
	Createdon       string                `json:"Createdon,omitempty"`
	VpcId           string                `json:"VpcId,omitempty"`
	Scheme          string                `json:"Scheme,omitempty"`
	DefaultResponse interface{}           `json:"DefaultResponse,omitempty"`
	LbDeleteStatus  string                `json:"LbDeleteStatus,omitempty"`
	ApplicationLb   []LoadBalanceResponse `json:"ApplicationLb,omitempty"`
	ClassicLb       []LoadBalanceResponse `json:"ClassicLb,omitempty"`
}

type DeleteLoadbalancerInput struct {
	LbName      string `json:"LbName,omitempty"`
	LbArn       string `json:"LbArn,omitempty"`
	TargetArn   string `json:"TargetArn,omitempty"`
	ListenerArn string `json:"ListenerArn,omitempty"`
}

type DescribeLoadbalancersInput struct {
	LbNames     []string `json:"LbName,omitempty"`
	LbArns      []string `json:"LbArn,omitempty"`
	TargetArns  []string `json:"TargetArns,omitempty"`
	ListnerArns []string `json:"ListnerArns,omitempty"`
}

func (sess *EstablishedSession) CreateClassicLb(lb LoadBalanceCreateInput) (*elb.CreateLoadBalancerOutput, error) {

	if sess.Elb != nil {
		listners := make([]*elb.Listener, 0)
		switch lb.Lbproto {
		case "HTTP":
			listners = append(listners, &elb.Listener{
				InstancePort:     aws.Int64(lb.InstPort),
				InstanceProtocol: aws.String(lb.Instproto),
				LoadBalancerPort: aws.Int64(lb.LbPort),
				Protocol:         aws.String(lb.Lbproto),
			})
		case "HTTPS":
			listners = append(listners, &elb.Listener{
				InstancePort:     aws.Int64(lb.InstPort),
				InstanceProtocol: aws.String(lb.Instproto),
				LoadBalancerPort: aws.Int64(lb.LbPort),
				Protocol:         aws.String(lb.Lbproto),
				SSLCertificateId: aws.String(lb.SslCert),
			})
		default:
			return nil, fmt.Errorf("You provided unknown loadbalancer protocol, enter a valid protocol")
		}
		input := &elb.CreateLoadBalancerInput{
			Listeners:        listners,
			LoadBalancerName: aws.String(lb.Name),
			Scheme:           aws.String(lb.Scheme),
			SecurityGroups:   aws.StringSlice(lb.SecurityGroups),
			Subnets:          aws.StringSlice(lb.Subnets),
		}
		result, err := (sess.Elb).CreateLoadBalancer(input)
		if err != nil {
			return nil, err
		}
		return result, nil
	}
	return nil, err.InvalidSession()
}

func (sess *EstablishedSession) CreateApplicationLb(lb LoadBalanceCreateInput) (*elbv2.CreateLoadBalancerOutput, error) {

	if sess.Elb2 != nil {
		input := &elbv2.CreateLoadBalancerInput{
			Name:           aws.String(lb.Name),
			Scheme:         aws.String(lb.Scheme),
			Subnets:        aws.StringSlice(lb.Subnets),
			SecurityGroups: aws.StringSlice(lb.SecurityGroups),
			IpAddressType:  aws.String(lb.IpAddressType),
			Tags: []*elbv2.Tag{
				{
					Key:   aws.String("Name"),
					Value: aws.String(lb.Name),
				}},
		}

		result, err := (sess.Elb2).CreateLoadBalancer(input)
		if err != nil {
			return nil, err
		}
		return result, nil
	}
	return nil, err.InvalidSession()

}

func (sess *EstablishedSession) CreateTargetGroups(lb *LoadBalanceCreateInput) (*elbv2.CreateTargetGroupOutput, error) {

	if sess.Elb2 != nil {
		input := &elbv2.CreateTargetGroupInput{
			Name:                       aws.String(lb.Name),
			Port:                       aws.Int64(lb.LbPort),
			Protocol:                   aws.String(lb.Lbproto),
			VpcId:                      aws.String(lb.VpcId),
			HealthCheckProtocol:        aws.String(lb.Instproto),
			HealthCheckPort:            aws.String(strconv.FormatInt(lb.InstPort, 10)),
			HealthCheckPath:            aws.String(lb.HealthPath),
			HealthCheckIntervalSeconds: aws.Int64(30),
			HealthCheckTimeoutSeconds:  aws.Int64(5),
			HealthyThresholdCount:      aws.Int64(5),
			UnhealthyThresholdCount:    aws.Int64(2),
			Matcher:                    &elbv2.Matcher{HttpCode: &lb.HttpCode},
		}

		result, err := (sess.Elb2).CreateTargetGroup(input)
		if err != nil {
			return nil, err
		}
		return result, nil
	}
	return nil, err.InvalidSession()
}

func (sess *EstablishedSession) CreateApplicationListners(lb *LoadBalanceCreateInput) (*elbv2.CreateListenerOutput, error) {

	if sess.Elb2 != nil {
		var input elbv2.CreateListenerInput
		switch lb.Lbproto {
		case "HTTP":
			input = elbv2.CreateListenerInput{
				DefaultActions: []*elbv2.Action{
					{
						TargetGroupArn: aws.String(lb.TargetArn),
						Type:           aws.String("forward"),
					},
				},
				LoadBalancerArn: aws.String(lb.LbArn),
				Port:            aws.Int64(lb.LbPort),
				Protocol:        aws.String(lb.Lbproto),
			}
		case "HTTPS":
			input = elbv2.CreateListenerInput{
				Certificates: []*elbv2.Certificate{
					{
						CertificateArn: aws.String(lb.SslCert),
					},
				},
				DefaultActions: []*elbv2.Action{
					{
						TargetGroupArn: aws.String(lb.TargetArn),
						Type:           aws.String("forward"),
					},
				},
				LoadBalancerArn: aws.String(lb.LbArn),
				Port:            aws.Int64(lb.LbPort),
				Protocol:        aws.String(lb.Lbproto),
				SslPolicy:       aws.String(lb.SslPolicy),
			}
		default:
			return nil, fmt.Errorf("You provided unknown loadbalancer protocol, enter a valid protocol")
		}

		result, err := (sess.Elb2).CreateListener(&input)
		if err != nil {
			return nil, err
		}
		return result, nil
	}
	return nil, err.InvalidSession()
}

func (sess *EstablishedSession) DescribeClassicLoadbalancer(lb *DescribeLoadbalancersInput) (*elb.DescribeLoadBalancersOutput, error) {

	if sess.Elb != nil {
		if lb.LbNames != nil {
			input := &elb.DescribeLoadBalancersInput{
				LoadBalancerNames: aws.StringSlice(lb.LbNames),
			}
			result, err := (sess.Elb).DescribeLoadBalancers(input)

			if err != nil {
				return nil, err
			}
			return result, nil
		}
		return nil, fmt.Errorf(fmt.Sprintf("%v DescribeClassicLoadbalancer", err.EmptyStructError()))
	}
	return nil, err.InvalidSession()
}

func (sess *EstablishedSession) DescribeAllClassicLoadbalancer(lb *DescribeLoadbalancersInput) (*elb.DescribeLoadBalancersOutput, error) {

	if sess.Elb != nil {
		input := &elb.DescribeLoadBalancersInput{}
		result, err := (sess.Elb).DescribeLoadBalancers(input)

		if err != nil {
			return nil, err
		}
		return result, nil
	}
	return nil, err.InvalidSession()
}

func (sess *EstablishedSession) DescribeApplicationLoadbalancer(lb *DescribeLoadbalancersInput) (*elbv2.DescribeLoadBalancersOutput, error) {

	if sess.Elb2 != nil {
		if lb.LbArns != nil {
			input := &elbv2.DescribeLoadBalancersInput{
				LoadBalancerArns: aws.StringSlice(lb.LbArns),
			}
			result, err := (sess.Elb2).DescribeLoadBalancers(input)

			if err != nil {
				return nil, err
			}
			return result, nil
		}

		if lb.LbNames != nil {
			input := &elbv2.DescribeLoadBalancersInput{
				Names: aws.StringSlice(lb.LbNames),
			}
			result, err := (sess.Elb2).DescribeLoadBalancers(input)

			if err != nil {
				return nil, err
			}
			return result, nil
		}

		return nil, fmt.Errorf(fmt.Sprintf("%v DescribeApplicationLoadbalancer", err.EmptyStructError()))
	}
	return nil, err.InvalidSession()
}

func (sess *EstablishedSession) DescribeAllApplicationLoadbalancer(lb *DescribeLoadbalancersInput) (*elbv2.DescribeLoadBalancersOutput, error) {

	if sess.Elb2 != nil {
		input := &elbv2.DescribeLoadBalancersInput{}
		result, err := (sess.Elb2).DescribeLoadBalancers(input)

		if err != nil {
			return nil, err
		}
		return result, nil
	}
	return nil, err.InvalidSession()
}

func (sess *EstablishedSession) DescribeTargetgroups(lb *DescribeLoadbalancersInput) (*elbv2.DescribeTargetGroupsOutput, error) {
	if sess.Elb2 != nil {
		if lb.TargetArns != nil {
			input := &elbv2.DescribeTargetGroupsInput{
				TargetGroupArns: aws.StringSlice(lb.TargetArns),
			}
			result, err := (sess.Elb2).DescribeTargetGroups(input)

			if err != nil {
				return nil, err
			}
			return result, nil
		}

		if lb.LbArns != nil {
			input := &elbv2.DescribeTargetGroupsInput{
				LoadBalancerArn: aws.String(lb.LbArns[0]),
			}
			result, err := (sess.Elb2).DescribeTargetGroups(input)

			if err != nil {
				return nil, err
			}
			return result, nil
		}
		return nil, fmt.Errorf(fmt.Sprintf("%v DescribeTargetgroups", err.EmptyStructError()))
	}
	return nil, err.InvalidSession()
}

func (sess *EstablishedSession) DescribeAllTargetgroups(lb *DescribeLoadbalancersInput) (*elbv2.DescribeTargetGroupsOutput, error) {

	if sess.Elb2 != nil {
		input := &elbv2.DescribeTargetGroupsInput{}
		result, err := (sess.Elb2).DescribeTargetGroups(input)

		if err != nil {
			return nil, err
		}
		return result, nil
	}
	return nil, err.InvalidSession()
}

func (sess *EstablishedSession) DescribeListners(lb *DescribeLoadbalancersInput) (*elbv2.DescribeListenersOutput, error) {

	if sess.Elb2 != nil {
		if lb.ListnerArns != nil {
			input := &elbv2.DescribeListenersInput{
				ListenerArns: aws.StringSlice(lb.ListnerArns),
			}
			result, err := (sess.Elb2).DescribeListeners(input)

			if err != nil {
				return nil, err
			}
			return result, nil
		}

		if lb.LbArns != nil {
			input := &elbv2.DescribeListenersInput{
				LoadBalancerArn: aws.String(lb.LbArns[0]),
			}
			result, err := (sess.Elb2).DescribeListeners(input)

			if err != nil {
				return nil, err
			}
			return result, nil
		}
		return nil, fmt.Errorf(fmt.Sprintf("%v DescribeListners", err.EmptyStructError()))
	}
	return nil, err.InvalidSession()

}

func (sess *EstablishedSession) DescribeAllListners(lb *DescribeLoadbalancersInput) (*elbv2.DescribeListenersOutput, error) {

	if sess.Elb2 != nil {
		input := &elbv2.DescribeListenersInput{}
		result, err := (sess.Elb2).DescribeListeners(input)

		if err != nil {
			return nil, err
		}
		return result, nil
	}
	return nil, err.InvalidSession()

}

func (sess *EstablishedSession) DeleteClassicLoadbalancer(lb *DeleteLoadbalancerInput) error {

	if sess.Elb != nil {
		input := &elb.DeleteLoadBalancerInput{
			LoadBalancerName: aws.String(lb.LbName),
		}
		_, err := (sess.Elb).DeleteLoadBalancer(input)

		if err != nil {
			return err
		}
		return nil
	}
	return err.InvalidSession()

}

func (sess *EstablishedSession) DeleteAppLoadbalancer(lb *DeleteLoadbalancerInput) error {

	if sess.Elb2 != nil {
		input := &elbv2.DeleteLoadBalancerInput{
			LoadBalancerArn: aws.String(lb.LbArn),
		}
		_, err := (sess.Elb2).DeleteLoadBalancer(input)

		if err != nil {
			return err
		}
		return nil
	}
	return err.InvalidSession()
}

func (sess *EstablishedSession) DeleteTargetGroup(lb *DeleteLoadbalancerInput) error {

	if sess.Elb2 != nil {
		if lb.TargetArn != "" {
			input := &elbv2.DeleteTargetGroupInput{
				TargetGroupArn: aws.String(lb.TargetArn),
			}
			_, err := (sess.Elb2).DeleteTargetGroup(input)

			if err != nil {
				return err
			}
			return nil
		}
		return fmt.Errorf(fmt.Sprintf("%v DeleteTargetGroup", err.EmptyStructError()))
	}
	return err.InvalidSession()
}

func (sess *EstablishedSession) DeleteAppListeners(lb *DeleteLoadbalancerInput) error {

	if sess.Elb2 != nil {
		if lb.ListenerArn != "" {
			input := &elbv2.DeleteListenerInput{
				ListenerArn: aws.String(lb.ListenerArn),
			}
			_, err := (sess.Elb2).DeleteListener(input)

			if err != nil {
				return err
			}
			return nil
		}
		return fmt.Errorf(fmt.Sprintf("%v DeleteAppListeners", err.EmptyStructError()))
	}
	return err.InvalidSession()
}

func (sess *EstablishedSession) WaitTillLbDeletionSuccessfull(lb *DescribeLoadbalancersInput) error {

	if sess.Elb2 != nil {
		input := &elbv2.DescribeLoadBalancersInput{
			LoadBalancerArns: aws.StringSlice(lb.LbArns),
		}
		err := (sess.Elb2).WaitUntilLoadBalancersDeleted(input)

		if err != nil {
			return err
		}
		return nil
	}
	return err.InvalidSession()
}
