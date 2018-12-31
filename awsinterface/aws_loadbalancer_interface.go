package DengineAwsInterface

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/elb"
	"github.com/aws/aws-sdk-go/service/elbv2"
	log "neuron/logger"
	"strconv"
	"strings"
	"time"
)

type LoadBalanceCreateInput struct {
	Name       string
	VpcId      string
	Scheme     string
	Type       string
	SslCert    string `json:"SslCert,omitempty"` //required only if the LB protocol is HTTPS else can be initiazed with dummy value
	SslPolicy  string `json:"SslPolicy,omitempty"`
	LbPort     int64  //required ex: 8080 or 80 etc
	InstPort   int64
	Lbproto    string //required ex: HTTPS, HTTP
	Instproto  string
	HttpCode   string `json:"HttpCode,omitempty"`
	HealthPath string `json:"HealthPath,omitempty"`
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

type ApplicationLBResponse struct {
	LbDns string
	LbArn string
}

type DeleteLoadbalancerInput struct {
	LbName string `json:"LbDns,omitempty"`
	LbArn  string `json:"LbArn,omitempty"`
}

func (load *LoadBalanceCreateInput) CreateLoadBalancer() LoadBalanceResponse {

	var (
		LB               LoadBalanceResponse
		AppLbResponse    ApplicationLBResponse
		lb_scheme        string
		subnet_response  []string
		sec_response     string
		target_response  string
		listner_response string
	)
	// collecting subnet details
	subnet_input := &ec2.DescribeSubnetsInput{
		Filters: []*ec2.Filter{
			{
				Name: aws.String("vpc-id"),
				Values: []*string{
					aws.String(load.VpcId),
				}}}}

	subnet_result, _ := Svc.DescribeSubnets(subnet_input)
	for _, subnet := range subnet_result.Subnets {
		subnet_response = append(subnet_response, *subnet.SubnetId)
	}

	// fetching security group so that I can attach it to the loabalancer which I am about to create
	sec_input := &ec2.DescribeSecurityGroupsInput{}

	sec_result, _ := Svc.DescribeSecurityGroups(sec_input)

	for _, sec := range sec_result.SecurityGroups {
		if *sec.VpcId == load.VpcId {
			sec_response = *sec.GroupId
		} else {
			// I am not supposed to anything
		}
	}
	// creating load balancer

	// selecting scheme
	if load.Scheme == "external" {
		lb_scheme = "internet-facing"
	} else if load.Scheme == "internal" {
		lb_scheme = "internal"
	} else {
		// I am not supposed to anything
	}

	instport := int64(load.InstPort)
	lbport := int64(load.LbPort)

	// creating LB according to the input ex: application/classic
	switch strings.ToLower(load.Type) {
	case "classic":

		// creating classic load balancer based on the protocol specified
		switch load.Lbproto {
		case "HTTP":
			lb_create_input := &elb.CreateLoadBalancerInput{
				Listeners: []*elb.Listener{
					{
						InstancePort:     &instport,
						InstanceProtocol: aws.String(load.Instproto),
						LoadBalancerPort: &lbport,
						Protocol:         aws.String(load.Lbproto),
					}},
				LoadBalancerName: aws.String(load.Name),
				Scheme:           aws.String(lb_scheme),
				SecurityGroups: []*string{
					aws.String(sec_response),
				},
				Subnets: aws.StringSlice(subnet_response),
			}
			lb_create_result, _ := Elb.CreateLoadBalancer(lb_create_input)
			lb_dns := lb_create_result.DNSName
			LB = LoadBalanceResponse{Name: load.Name, Type: load.Type, LbDns: *lb_dns}

		case "HTTPS":
			lb_create_input := &elb.CreateLoadBalancerInput{
				Listeners: []*elb.Listener{
					{
						InstancePort:     &instport,
						InstanceProtocol: aws.String(load.Instproto),
						LoadBalancerPort: &lbport,
						Protocol:         aws.String(load.Lbproto),
						SSLCertificateId: aws.String(load.SslCert),
					}},
				LoadBalancerName: aws.String(load.Name),
				Scheme:           aws.String(lb_scheme),
				SecurityGroups: []*string{
					aws.String(sec_response),
				},
				Subnets: aws.StringSlice(subnet_response),
			}
			lb_create_result, _ := Elb.CreateLoadBalancer(lb_create_input)
			lb_dns := lb_create_result.DNSName
			LB = LoadBalanceResponse{Name: load.Name, Type: load.Type, LbDns: *lb_dns}

		default:
			log.Info("")
			log.Error("I feel we are lost while creating protocol specific classic loadbalancer")
			log.Info("")
			LB = LoadBalanceResponse{DefaultResponse: "I feel we are lost while creating protocol specific classic loadbalancer, for more details check log file"}
		}

	case "application":

		// creating load balancer logic
		lb_create_input := &elbv2.CreateLoadBalancerInput{
			Name:           aws.String(load.Name),
			Scheme:         aws.String(lb_scheme),
			Subnets:        aws.StringSlice(subnet_response),
			SecurityGroups: []*string{aws.String(sec_response)},
			IpAddressType:  aws.String("ipv4"),
			Tags: []*elbv2.Tag{
				{
					Key:   aws.String("Name"),
					Value: aws.String(load.Name),
				}},
		}

		lb_create_response, _ := Elb2.CreateLoadBalancer(lb_create_input)

		for _, loadbal := range lb_create_response.LoadBalancers {
			AppLbResponse = ApplicationLBResponse{*loadbal.DNSName, *loadbal.LoadBalancerArn}
		}

		// creating target group
		portint := strconv.FormatInt(load.InstPort, 10)
		target_group_input := &elbv2.CreateTargetGroupInput{
			Name:                       aws.String(load.Name),
			Port:                       aws.Int64(load.LbPort),
			Protocol:                   aws.String(load.Lbproto),
			VpcId:                      aws.String(load.VpcId),
			HealthCheckProtocol:        aws.String(load.Instproto),
			HealthCheckPort:            aws.String(portint),
			HealthCheckPath:            aws.String(load.HealthPath),
			HealthCheckIntervalSeconds: aws.Int64(30),
			HealthCheckTimeoutSeconds:  aws.Int64(5),
			HealthyThresholdCount:      aws.Int64(5),
			UnhealthyThresholdCount:    aws.Int64(2),
			Matcher:                    &elbv2.Matcher{HttpCode: &load.HttpCode},
		}

		target_group_response, _ := Elb2.CreateTargetGroup(target_group_input)
		for _, target := range target_group_response.TargetGroups {
			target_response = *target.TargetGroupArn
		}

		switch load.Lbproto {

		case "HTTP":

			listner_create_input := &elbv2.CreateListenerInput{
				DefaultActions: []*elbv2.Action{
					{
						TargetGroupArn: aws.String(target_response),
						Type:           aws.String("forward"),
					},
				},
				LoadBalancerArn: aws.String(AppLbResponse.LbArn),
				Port:            aws.Int64(80),
				Protocol:        aws.String("HTTP"),
			}

			listner_create_response, _ := Elb2.CreateListener(listner_create_input)
			for _, listners := range listner_create_response.Listeners {
				listner_response = *listners.ListenerArn
			}
			LB = LoadBalanceResponse{Name: load.Name, Type: load.Type, LbDns: AppLbResponse.LbDns, LbArn: AppLbResponse.LbArn, TargetArn: target_response, ListnerArn: listner_response}

		case "HTTPS":

			listner_create_input := &elbv2.CreateListenerInput{
				Certificates: []*elbv2.Certificate{
					{
						CertificateArn: aws.String(load.SslCert),
					},
				},
				DefaultActions: []*elbv2.Action{
					{
						TargetGroupArn: aws.String(target_response),
						Type:           aws.String("forward"),
					},
				},
				LoadBalancerArn: aws.String(AppLbResponse.LbArn),
				Port:            aws.Int64(443),
				Protocol:        aws.String("HTTPS"),
				SslPolicy:       aws.String(load.SslPolicy),
			}

			listner_create_response, _ := Elb2.CreateListener(listner_create_input)
			for _, listner := range listner_create_response.Listeners {
				listner_response = *listner.ListenerArn
			}
			LB = LoadBalanceResponse{Name: load.Name, Type: load.Type, LbDns: AppLbResponse.LbDns, LbArn: AppLbResponse.LbArn, TargetArn: target_response}

		default:
			log.Info("")
			log.Error("I feel we are lost while creating protocol specific application loadbalancer")
			log.Info("")
			LB = LoadBalanceResponse{DefaultResponse: "I feel we are lost while creating protocol specific application loadbalancer"}
		}

	default:

		log.Info("")
		log.Error("I feel we are lost in creating LB :S")
		log.Info("")
		LB = LoadBalanceResponse{DefaultResponse: "I feel we are lost in creating LB :S"}
	}
	return LB
}

func (d *DeleteLoadbalancerInput) DeleteLoadbalancer() LoadBalanceResponse {

	var lb_response LoadBalanceResponse
	if d.LbArn == "" {

		var lb_search string
		find_lb_input := &elb.DescribeLoadBalancersInput{}
		result, _ := Elb.DescribeLoadBalancers(find_lb_input)
		for _, val := range result.LoadBalancerDescriptions {
			if d.LbName == *val.LoadBalancerName {
				lb_search = "Found"
			}
		}

		if lb_search == "Found" {

			input := &elb.DeleteLoadBalancerInput{
				LoadBalancerName: aws.String(d.LbName),
			}
			_, del_err := Elb.DeleteLoadBalancer(input)

			if del_err == nil {
				lb_response = LoadBalanceResponse{LbDeleteStatus: "LoadBalancer deletion is successful"}
			} else {
				lb_response = LoadBalanceResponse{LbDeleteStatus: "LoadBalancer deletion is failure"}
			}

		} else {
			lb_response = LoadBalanceResponse{DefaultResponse: "Could not find the entered loadbalancer, please enter valid/existing loadbalancer name"}
		}

	}
	if d.LbName == "" {

		var lb_search string
		find_lb_input := &elbv2.DescribeLoadBalancersInput{}
		result, _ := Elb2.DescribeLoadBalancers(find_lb_input)
		for _, val := range result.LoadBalancers {
			if d.LbArn == *val.LoadBalancerArn {
				lb_search = "Found"
			}
		}

		if lb_search == "Found" {

			// searching dependency of the load balancers to be deleted. And the deletion process will be carried out along with its dependency.
			find_target := &elbv2.DescribeTargetGroupsInput{LoadBalancerArn: aws.String(d.LbArn)}
			target_result, _ := Elb2.DescribeTargetGroups(find_target)
			for _, tar_val := range target_result.TargetGroups {

				// deletion of only loadbalancer will be carried out by next 6 lines of code.
				del_lb_input := &elbv2.DeleteLoadBalancerInput{LoadBalancerArn: aws.String(d.LbArn)}
				_, err_lb := Elb2.DeleteLoadBalancer(del_lb_input)
				if err_lb != nil {
					lb_response = LoadBalanceResponse{DefaultResponse: err_lb}
				}

				wait_till_comp := &elbv2.DescribeLoadBalancersInput{LoadBalancerArns: aws.StringSlice([]string{d.LbArn})}
				Elb2.WaitUntilLoadBalancersDeleted(wait_till_comp)

				time.Sleep(5 * time.Second)
				// deletion of the target group will be done by below snippet.
				tar_del_input := &elbv2.DeleteTargetGroupInput{TargetGroupArn: aws.String(*tar_val.TargetGroupArn)}
				_, err_tar := Elb2.DeleteTargetGroup(tar_del_input)
				if err_tar != nil {
					lb_response = LoadBalanceResponse{DefaultResponse: err_tar}
				}

				if (err_lb != nil) || (err_tar != nil) {
					lb_response = LoadBalanceResponse{LbDeleteStatus: "LoadBalancer deletion is failure"}
				} else {
					lb_response = LoadBalanceResponse{LbDeleteStatus: "LoadBalancer deletion is successful"}
				}
			}

		} else {
			lb_response = LoadBalanceResponse{DefaultResponse: "Could not find the entered loadbalancer, please enter valid/existing loadbalancer ARN"}
		}

	}
	return lb_response
}

func GetAllLoadbalancer() LoadBalanceResponse {

	lb_chn := make(chan []LoadBalanceResponse, 2)
	go func() {
		lb_chn <- GetAllClassicLb()
	}()
	go func() {
		lb_chn <- GetAllApplicationLb()
	}()

	lb_response := LoadBalanceResponse{ClassicLb: <-lb_chn, ApplicationLb: <-lb_chn}
	close(lb_chn)

	return lb_response
}

func GetAllClassicLb() []LoadBalanceResponse {

	// searching all the classic loadbalancer
	search_lb_input := &elb.DescribeLoadBalancersInput{}
	search_lb_result, search_err := Elb.DescribeLoadBalancers(search_lb_input)

	var lb_list []LoadBalanceResponse
	if search_err != nil {
		lb_list = append(lb_list, LoadBalanceResponse{DefaultResponse: "We encountered errors while searching for classic load balancers in AWS"})
	} else {
		for _, lb := range search_lb_result.LoadBalancerDescriptions {
			lb_list = append(lb_list, LoadBalanceResponse{Name: *lb.LoadBalancerName, LbDns: *lb.DNSName, Createdon: (*lb.CreatedTime).String(), Type: "classic", Scheme: *lb.Scheme, VpcId: *lb.VPCId})
		}
	}
	return lb_list
}

func GetAllApplicationLb() []LoadBalanceResponse {

	// searching all the application loadbalancer
	search_lb_input := &elbv2.DescribeLoadBalancersInput{}
	search_lb_result, search_err := Elb2.DescribeLoadBalancers(search_lb_input)

	var lb_list []LoadBalanceResponse
	if search_err != nil {
		lb_list = append(lb_list, LoadBalanceResponse{DefaultResponse: "We encountered errors while searching for application load balancers in AWS"})
	} else {
		for _, lb := range search_lb_result.LoadBalancers {

			var tar_arn string
			// searching target group for the corresponding
			search_tar_input := &elbv2.DescribeTargetGroupsInput{}
			search_lb_result, _ := Elb2.DescribeTargetGroups(search_tar_input)
			for _, tar := range search_lb_result.TargetGroups {
				if *lb.LoadBalancerArn == *tar.LoadBalancerArns[0] {
					tar_arn = *tar.TargetGroupArn
				}
			}
			lb_list = append(lb_list, LoadBalanceResponse{Name: *lb.LoadBalancerName, LbDns: *lb.DNSName, LbArn: *lb.LoadBalancerArn, Createdon: (*lb.CreatedTime).String(), Type: *lb.Type, Scheme: *lb.Scheme, VpcId: *lb.VpcId, TargetArn: tar_arn})
		}
	}
	return lb_list
}
