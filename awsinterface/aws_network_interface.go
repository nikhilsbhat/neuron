package DengineAwsInterface

import (
	"github.com/aws/aws-sdk-go/aws"
	//"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/ec2"
	//log "neuron/logger"
	"strconv"
	"strings"
)

type VpcCreateInput struct {
	Cidr  string
	Name  string
	Type  string
	Ports []string `json:"Ports,omitempty"`
}

type SubnetCreateInput struct {
	Cidr    string
	Name    string
	Zone    string
	VpcData VpcReponse `json:"VpcData,omitempty"`
	VpcId   string     `json:"VpcId,omitempty"`
	Type    string     `json:"Type,omitempty"`
}

type IgwCreateInput struct {
	Name  string
	VpcId string
}

type SecurityCreateInput struct {
	Name  string
	VpcId string
	Ports []string
}

type CreateRouteTableInput struct {
	Name    string
	VpcId   string
	SubId   string
	IgwId   string `json:"IgwId,omitempty"`
	SubType string
}

type DeleteNetworkInput struct {
	VpcId string
}

type GetNetworksInput struct {
	VpcIds []string
}

type NetworkResponse struct {
	Name      string          `json:"Name,omitempty"`
	VpcId     string          `json:"VpcId,omitempty"`
	Subnets   []SubnetReponse `json:"Subnets,omitempty"`
	State     string          `json:"State,omitempty"`
	Igw       string          `json:"Igw,omitempty"`
	IsDefault bool            `json:"IsDefault,omitempty"`
	Region    string          `json:"Region,omitempty"`
}

type VpcReponse struct {
	Name string `json:"Name,omitempty"`
	Id   string `json:"Id,omitempty"`
	Type string `json:"Type,omitempty"`
	Igw  string `json:"Igw,omitempty"`
}

type SubnetReponse struct {
	Name string `json:"Name,omitempty"`
	Id   string `json:"Id,omitempty"`
}

type DeleteNetworkResponse struct {
	Subnets         []string `json:"Subnets,omitempty"`
	Securities      []string `json:"Securities,omitempty"`
	Routetables     []string `json:"Routetables,omitempty"`
	Gateways        []string `json:"Gateways,omitempty"`
	Vpcs            string   `json:"Vpcs,omitempty"`
	DefaultResponse string   `json:"DefaultResponse,omitempty"`
}

/*func init() {

	// Print message stating this interface is being called
	log.Info("")
	log.Info("I WAS INVOKED")
	log.Info("AND I AM AWS INTERFACE")
	log.Info("")

}*/

// being get_availability_zones my job is to list/give back all available availabilityzones in the region specified to me
func GetAvailabilityZones() ([]string, error) {

	availability_zone_input := &ec2.DescribeAvailabilityZonesInput{}
	// I will be returning an array full of zones for the region asked from me
	result, zone_err := Svc.DescribeAvailabilityZones(availability_zone_input)
	if zone_err != nil {
		return nil, zone_err
	} else {
		availabilityzones := result.AvailabilityZones
		zones := make([]string, 0)
		for _, zone := range availabilityzones {
			zones = append(zones, *zone.ZoneName)
		}
		return zones, nil
	}

}

// being create_vpc my job is to create vpc and give back the response who called me
func (vpcin *VpcCreateInput) CreateVpc() (VpcReponse, error) {

	// I am gathering inputs since create vpc needs it
	vpc_create_input := &ec2.CreateVpcInput{
		CidrBlock:       aws.String(vpcin.Cidr),
		InstanceTenancy: aws.String("default"),
	}

	// I will create vpc by collecting values as input from above function
	vpc_result, vpc_err := Svc.CreateVpc(vpc_create_input)

	// handling the error if it throws while vpc is under creation process
	if vpc_err != nil {
		return VpcReponse{}, vpc_err
	} else {

		// I will program wait untill vpc become available
		vpc_status_input := &ec2.DescribeVpcsInput{
			Filters: []*ec2.Filter{
				&ec2.Filter{
					Name:   aws.String("state"),
					Values: aws.StringSlice([]string{"available"}),
				},
			},
		}

		Svc.WaitUntilVpcAvailable(vpc_status_input)

		// I will pass name to create_tags to set a name to the vpc
		tags := Tag{*vpc_result.Vpc.VpcId, "Name", vpcin.Name}
		vpc_tag, tag_err := create_tags(tags)
		if tag_err != nil {
			return VpcReponse{}, tag_err
		}

		// I will make the decision whether we need public network or private, based on the input I recieve
		var igw_id string
		if strings.ToLower(vpcin.Type) == "public" {
			igw := IgwCreateInput{vpcin.Name + "_igw", *vpc_result.Vpc.VpcId}
			igw_response, igw_err := igw.create_igw()
			if igw_err != nil {
				return VpcReponse{}, igw_err
			} else {
				igw_id = igw_response
            }
		} else if strings.ToLower(vpcin.Type) == "private" {
			igw_id = "This is private network"
		} else {
			// I was told nothing to do
		}

		// I will initialize data required to create security group and pass it to respective person to create one
		security := SecurityCreateInput{vpcin.Name + "_Sec", *vpc_result.Vpc.VpcId, vpcin.Ports}
		sec_err := security.create_security()
		if sec_err != nil {
			return VpcReponse{}, sec_err
		}

		return VpcReponse{vpc_tag, *vpc_result.Vpc.VpcId, vpcin.Type, igw_id}, nil
	}

}

// being create_subnets my job is to create subnets and give back the response who called me
func (subin *SubnetCreateInput) CreateSubnet() (SubnetReponse, error) {

	var VpcId string
	if subin.VpcData.Id != "" {
		VpcId = subin.VpcData.Id
	} else {
		VpcId = subin.VpcId
	}
	// I am gathering inputs since create subnets needs it
	subnet_create_input := &ec2.CreateSubnetInput{
		CidrBlock:        aws.String(subin.Cidr),
		VpcId:            aws.String(VpcId),
		AvailabilityZone: aws.String(subin.Zone),
	}

	// I will create subnet by collecting values as input from above function
	subnet_result, sub_err := Svc.CreateSubnet(subnet_create_input)
	// handling the error if it throws while subnet is under creation process
	if sub_err != nil {
		return SubnetReponse{}, sub_err
	} else {

		// I will program wait untill subnet become available
		subnet_status_input := &ec2.DescribeSubnetsInput{
			Filters: []*ec2.Filter{
				&ec2.Filter{
					Name:   aws.String("state"),
					Values: aws.StringSlice([]string{"available"}),
				},
			},
		}
		Svc.WaitUntilSubnetAvailable(subnet_status_input)

		// I will pass name to create_tags to set a name to the vpc
		tags := Tag{*subnet_result.Subnet.SubnetId, "Name", subin.Name}
		sub_tag, tag_err := create_tags(tags)
		if tag_err != nil {
			return SubnetReponse{}, tag_err
		}

		var routes CreateRouteTableInput
		if subin.VpcData.Id != "" {
			routes = CreateRouteTableInput{Name: subin.Name + "_route", VpcId: subin.VpcData.Id, SubId: *subnet_result.Subnet.SubnetId, IgwId: subin.VpcData.Igw, SubType: subin.VpcData.Type}
		} else {
			routes = CreateRouteTableInput{Name: subin.Name + "_route", VpcId: subin.VpcId, SubId: *subnet_result.Subnet.SubnetId, SubType: subin.Type}
		}
		route_err := routes.create_route_table()
		if route_err != nil {
			return SubnetReponse{}, nil
		} else {
			return SubnetReponse{sub_tag, *subnet_result.Subnet.SubnetId}, nil
		}
	}
}

func (igw *IgwCreateInput) create_igw() (string, error) {

	// I am here to create internet gateway if incase it is required
	igw_create_input := &ec2.CreateInternetGatewayInput{}
	igw_result, igw_err := Svc.CreateInternetGateway(igw_create_input)
	// handling the error if it throws while IGW is under creation process
	if igw_err != nil {
		return "", igw_err
	}

	// I will pass name to create_tags to set a name to the igw which I created just now
	tags := Tag{*igw_result.InternetGateway.InternetGatewayId, "Name", igw.Name}
	_, tag_err := create_tags(tags)
	if tag_err != nil {
		return "", tag_err
	}

	// I am attaching igw to the vpc which I am told to
	attach_igw_input := &ec2.AttachInternetGatewayInput{
		InternetGatewayId: aws.String(*igw_result.InternetGateway.InternetGatewayId),
		VpcId:             aws.String(igw.VpcId),
	}
	_, att_err := Svc.AttachInternetGateway(attach_igw_input)
	if att_err != nil {
		return "", att_err
	}

	return *igw_result.InternetGateway.InternetGatewayId, nil

}

// being create_security my job is to create security group and give back the response who called me
func (sec *SecurityCreateInput) create_security() error {

	security_create_input := &ec2.CreateSecurityGroupInput{
		Description: aws.String("This security group is created automatically by D-Engine api"),
		VpcId:       aws.String(sec.VpcId),
		GroupName:   aws.String(sec.Name),
	}

	security_result, sec_err := Svc.CreateSecurityGroup(security_create_input)

	if sec_err != nil {
		return sec_err
	} else {
		//creating egree and ingres rules for the security group which I created just now
		for _, port := range sec.Ports {
			int_port, _ := strconv.ParseInt(port, 10, 64)
			ingres_err := create_ingress_rule(int_port, *security_result.GroupId)
			if ingres_err != nil {
				return ingres_err
			}else {
				return nil
			}
		}
		egres_err := create_egress_rule(*security_result.GroupId)
		if egres_err != nil {
			return egres_err
		} else {
			return nil
		}
	}

}

//being create_route_table I am responsible for creating route-table and writing routes into it and attaching it to the vpc
func (route *CreateRouteTableInput) create_route_table() error {

	route_table_create_input := &ec2.CreateRouteTableInput{
		VpcId: aws.String(route.VpcId),
	}
	route_table_create_result, route_err := Svc.CreateRouteTable(route_table_create_input)

	if route_err != nil {
		return route_err
	} else {
		// I will pass name to create_tags to set a name to the route table which I created just now
		tags := Tag{*route_table_create_result.RouteTable.RouteTableId, "Name", route.Name}
		_, tag_err := create_tags(tags)
		if tag_err != nil {
			return tag_err
		}

		// Following codes are responsible for writing routes in the route table created earlier.
		if route.IgwId != "" {
			if strings.ToLower(route.SubType) == "public" {
				create_route_input := &ec2.CreateRouteInput{
					DestinationCidrBlock: aws.String("0.0.0.0/0"),
					GatewayId:            aws.String(route.IgwId),
					RouteTableId:         aws.String(*route_table_create_result.RouteTable.RouteTableId),
				}
				_, cr_rou_err := Svc.CreateRoute(create_route_input)
				if cr_rou_err != nil {
					return cr_rou_err
				} else {
					return nil
				}

				assoiate_route_table := &ec2.AssociateRouteTableInput{
					RouteTableId: aws.String(*route_table_create_result.RouteTable.RouteTableId),
					SubnetId:     aws.String(route.SubId),
				}
				_, asso_rou_err := Svc.AssociateRouteTable(assoiate_route_table)
				if asso_rou_err != nil {
					return asso_rou_err
				} else {
					return nil
				}
			} else {
				// since it is private network which has to be created there is no need to perform additional task
				return nil
			}
		} else {

			if strings.ToLower(route.SubType) == "public" {
				igw_search_input := &ec2.DescribeInternetGatewaysInput{}
				igw_search_result, _ := Svc.DescribeInternetGateways(igw_search_input)
				for _, igw := range igw_search_result.InternetGateways {
					if *igw.Attachments[0].VpcId == route.VpcId {
						create_route_input := &ec2.CreateRouteInput{
							DestinationCidrBlock: aws.String("0.0.0.0/0"),
							GatewayId:            aws.String(*igw.InternetGatewayId),
							RouteTableId:         aws.String(*route_table_create_result.RouteTable.RouteTableId),
						}
						_, cr_rou_err := Svc.CreateRoute(create_route_input)
						if cr_rou_err != nil {
							return cr_rou_err
						} else {
							return nil
						}

						assoiate_route_table := &ec2.AssociateRouteTableInput{
							RouteTableId: aws.String(*route_table_create_result.RouteTable.RouteTableId),
							SubnetId:     aws.String(route.SubId),
						}
						_, asso_rou_err := Svc.AssociateRouteTable(assoiate_route_table)
						if asso_rou_err != nil {
							return asso_rou_err
						} else {
							return nil
						}
					}
				}
			} else {
				// Releasing Soon...!!. We are not supporting writing custom routes into route tables as of now.
				return nil
			}
		}
		return nil
	}
}

// I will be called by create_security to attach egress rules as required
func create_egress_rule(id string) error {

	security_ingress_input := &ec2.AuthorizeSecurityGroupEgressInput{
		//    FromPort   : aws.int64(from_port),
		IpProtocol: aws.String("-1"),
		GroupId:    aws.String(id),
		//    ToPort     : aws.int64(to_port),
		CidrIp: aws.String("0.0.0.0/0"),
	}

	_, egress_err := Svc.AuthorizeSecurityGroupEgress(security_ingress_input)
	if egress_err != nil {
		return egress_err
	} else {
		return nil
	}
}

// I will be called by create_security to attach ingress rules as required
func create_ingress_rule(i int64, id string) error {

	port := int64(i)
	security_ingress_input := &ec2.AuthorizeSecurityGroupIngressInput{
		FromPort:   &port,
		IpProtocol: aws.String("tcp"),
		GroupId:    aws.String(id),
		ToPort:     &port,
		CidrIp:     aws.String("0.0.0.0/0"),
	}

	_, ingress_err := Svc.AuthorizeSecurityGroupIngress(security_ingress_input)
	if ingress_err != nil {
		return ingress_err
	} else {
		return nil
	}

}

// Deletion function are below
func (d *DeleteNetworkInput) DeleteNetwork() (DeleteNetworkResponse, error) {

	var vpc_search string
	search_vpc_input := &ec2.DescribeVpcsInput{}
	vpc_search_result, vpc_search_err := Svc.DescribeVpcs(search_vpc_input)

	if vpc_search_err != nil {
		return DeleteNetworkResponse{}, vpc_search_err
	} else {
		for _, val := range vpc_search_result.Vpcs {
			if d.VpcId == *val.VpcId {
				vpc_search = "Found"
			}
		}
	}

	if vpc_search == "Found" {

		//collecting all subnet details of the particular vpc, so that the data of these can be erased
		fetch_subnet_input := &ec2.DescribeSubnetsInput{}
		fetch_subnet_response, sub_err := Svc.DescribeSubnets(fetch_subnet_input)
		var subnets []string
		var subnetsdeletestatus []string
		if sub_err != nil {
			return DeleteNetworkResponse{}, nil
		} else {
			for _, subnet := range fetch_subnet_response.Subnets {
				if *subnet.VpcId == d.VpcId {
					subnets = append(subnets, *subnet.SubnetId)
				}
			}
			//deleting Subnets of the vpc mentioned to achieve the final work
			for _, subnet := range subnets {
				subnet_input := &ec2.DeleteSubnetInput{
					SubnetId: aws.String(subnet),
				}
				_, del_sub_err := Svc.DeleteSubnet(subnet_input)
				if del_sub_err != nil {
					subnetsdeletestatus = append(subnetsdeletestatus, "Subnets deletion is failure")
				} else {
					subnetsdeletestatus = append(subnetsdeletestatus, "Subnets deletion is successful")
				}
			}

			//fetching info about the security group associated with the vpc
			fetch_sec_input := &ec2.DescribeSecurityGroupsInput{}
			fetch_sec_response, sec_err := Svc.DescribeSecurityGroups(fetch_sec_input)
			var securities []string
			var securitiesdeletestatus []string
			if sec_err != nil {
				return DeleteNetworkResponse{}, nil
			} else {
				for _, security := range fetch_sec_response.SecurityGroups {
					if (*security.VpcId == d.VpcId) && (*security.GroupName != "default") {
						securities = append(securities, *security.GroupId)
					}
				}
				//deleting security groups of the vpc mentioned to achieve the final work
				for _, security := range securities {
					sec_input := &ec2.DeleteSecurityGroupInput{
						GroupId: aws.String(security),
					}
					_, del_sec_err := Svc.DeleteSecurityGroup(sec_input)
					if del_sec_err != nil {
						securitiesdeletestatus = append(securitiesdeletestatus, "SecurityGroup deletion is failure")
					} else {
						securitiesdeletestatus = append(securitiesdeletestatus, "SecurityGroup deletion is successful")
					}
				}

				//fetching info about the route tables associated with the vpc
				fetch_routetable_input := &ec2.DescribeRouteTablesInput{}
				fetch_routetable_response, _ := Svc.DescribeRouteTables(fetch_routetable_input)
				var routetables []string
				var routedeletestatus []string
				for _, routetable := range fetch_routetable_response.RouteTables {
					if *routetable.VpcId == d.VpcId {
						routetables = append(routetables, *routetable.RouteTableId)
					}
				}
				//deleting routetables of the vpc mentioned to achieve the final work
				for _, table := range routetables {
					route_delete_input := &ec2.DeleteRouteTableInput{
						RouteTableId: aws.String(table),
					}
					_, route_err := Svc.DeleteRouteTable(route_delete_input)
					if route_err != nil {
						routedeletestatus = append(routedeletestatus, "RouteTable deletion is failure")
					} else {
						routedeletestatus = append(routedeletestatus, "RouteTable deletion is successful")
					}
				}

				//fetching info about the IGW's associated with the vpc
				fetch_igw_input := &ec2.DescribeInternetGatewaysInput{}
				fetch_igw_result, _ := Svc.DescribeInternetGateways(fetch_igw_input)
				var gateways []string
				var igwdeletestatus []string
				for _, igws := range fetch_igw_result.InternetGateways {
					/*var igwvpc string
					for _, igw := range igws.Attachments {
						igwvpc = *igw.VpcId
					}*/
					if *igws.Attachments[0].VpcId == d.VpcId {
						gateways = append(gateways, *igws.InternetGatewayId)
					}
				}
				//detaching IGW's from vpc so that it will be able to remove from vpc
				for _, gateway := range gateways {
					inter_delete_input := &ec2.DetachInternetGatewayInput{
						InternetGatewayId: aws.String(gateway),
						VpcId:             aws.String(d.VpcId),
					}
					_, dettach_err := Svc.DetachInternetGateway(inter_delete_input)
					if dettach_err != nil {
						return DeleteNetworkResponse{}, dettach_err
					}
				}
				//deleting IGW's of the vpc mentioned to achieve the final work
				for _, gateway := range gateways {
					inter_input := &ec2.DeleteInternetGatewayInput{
						InternetGatewayId: aws.String(gateway),
					}
					_, igw_err := Svc.DeleteInternetGateway(inter_input)
					if igw_err != nil {
						igwdeletestatus = append(igwdeletestatus, "IGW deletion is failure")
					} else {
						igwdeletestatus = append(igwdeletestatus, "IGW deletion is successful")
					}
				}

				//deleting the vpc as per the inputs
				var vpcdeletestatus string
				vpc_delete_input := &ec2.DeleteVpcInput{
					VpcId: aws.String(d.VpcId),
				}
				_, vpc_err := Svc.DeleteVpc(vpc_delete_input)
				if vpc_err != nil {
					vpcdeletestatus = "VPC deletion is failure"
				} else {
					vpcdeletestatus = "VPC deletion is successful"
				}

				return DeleteNetworkResponse{Subnets: subnetsdeletestatus, Securities: securitiesdeletestatus, Routetables: routedeletestatus, Gateways: igwdeletestatus, Vpcs: vpcdeletestatus}, nil
			}
		}

	} else {

		return DeleteNetworkResponse{DefaultResponse: "Could not find the entered VPC, please enter valid/existing VPC id"}, nil

	}

}

// I will get you back the list of subnets
func GetSubnets() ([]string, error) {

	subnet_input := &ec2.DescribeSubnetsInput{}
	subnet_response, sub_err := Svc.DescribeSubnets(subnet_input)
	if sub_err != nil {
		return nil, sub_err
	} else {
		var subnets []string
		for _, subnet := range subnet_response.Subnets {
			subnets = append(subnets, *subnet.SubnetId)
		}
		return subnets, nil
	}
}

func (net *GetNetworksInput) GetNetwork() ([]NetworkResponse, error) {

	var network_response []NetworkResponse
	find_vpc_input := &ec2.DescribeVpcsInput{VpcIds: aws.StringSlice(net.VpcIds)}
	find_vpc_result, vpc_err := Svc.DescribeVpcs(find_vpc_input)

	if vpc_err != nil {
		return nil, vpc_err
	} else {
		for _, vpc := range find_vpc_result.Vpcs {

			// searching subnets
			search_subnet_input := &ec2.DescribeSubnetsInput{}
			search_subnet_response, sub_err := Svc.DescribeSubnets(search_subnet_input)
			if sub_err != nil {
				return nil, sub_err
			} else {
				var subnets_list []SubnetReponse
				for _, subnet := range search_subnet_response.Subnets {
					if *subnet.VpcId == *vpc.VpcId {
						if subnet.Tags != nil {
							subnets_list = append(subnets_list, SubnetReponse{*subnet.Tags[0].Value, *subnet.SubnetId})
						} else {
							subnets_list = append(subnets_list, SubnetReponse{Id: *subnet.SubnetId})
						}
					}
				}
				// searching igws
				igw_search_input := &ec2.DescribeInternetGatewaysInput{}
				igw_search_result, igw_err := Svc.DescribeInternetGateways(igw_search_input)
				var igws string
				if igw_err != nil {
					return nil, igw_err
				} else {
					for _, igw := range igw_search_result.InternetGateways {
						if *igw.Attachments[0].VpcId == *vpc.VpcId {
							igws = *igw.InternetGatewayId
						}
					}
				}
				if vpc.Tags != nil {
					if igws == "" {
						network_response = append(network_response, NetworkResponse{Name: *vpc.Tags[0].Value, VpcId: *vpc.VpcId, Subnets: subnets_list, State: *vpc.State, IsDefault: *vpc.IsDefault})
						return network_response, nil
					} else {
						network_response = append(network_response, NetworkResponse{Name: *vpc.Tags[0].Value, VpcId: *vpc.VpcId, Subnets: subnets_list, State: *vpc.State, Igw: igws, IsDefault: *vpc.IsDefault})
						return network_response, nil
					}

				} else {
					if igws == "" {
						network_response = append(network_response, NetworkResponse{VpcId: *vpc.VpcId, Subnets: subnets_list, State: *vpc.State, IsDefault: *vpc.IsDefault})
					} else {
						network_response = append(network_response, NetworkResponse{VpcId: *vpc.VpcId, Subnets: subnets_list, State: *vpc.State, Igw: igws, IsDefault: *vpc.IsDefault})
					}
				}
			}
		}
		return network_response, nil
	}
}

func GetAllNetworks(region string) ([]NetworkResponse, error) {

	var network_response []NetworkResponse
	find_vpc_input := &ec2.DescribeVpcsInput{}
	find_vpc_result, vpc_err := Svc.DescribeVpcs(find_vpc_input)
	if vpc_err != nil {
		return nil, vpc_err
	} else {
		for _, vpc := range find_vpc_result.Vpcs {

			// searching subnets
			search_subnet_input := &ec2.DescribeSubnetsInput{}
			search_subnet_response, sub_err := Svc.DescribeSubnets(search_subnet_input)
			var subnets_list []SubnetReponse
			if sub_err != nil {
				return nil, sub_err
			} else {
				for _, subnet := range search_subnet_response.Subnets {
					if *subnet.VpcId == *vpc.VpcId {
						if subnet.Tags != nil {
							subnets_list = append(subnets_list, SubnetReponse{Name: *subnet.Tags[0].Value, Id: *subnet.SubnetId})
						} else {
							subnets_list = append(subnets_list, SubnetReponse{Id: *subnet.SubnetId})
						}
					}
				}
				// searching igws
				igw_search_input := &ec2.DescribeInternetGatewaysInput{}
				igw_search_result, igw_err := Svc.DescribeInternetGateways(igw_search_input)
				var igws string
				if igw_err != nil {
					return nil, igw_err
				} else {
					for _, igw := range igw_search_result.InternetGateways {
						if *igw.Attachments[0].VpcId == *vpc.VpcId {
							igws = *igw.InternetGatewayId
						}
					}
					if vpc.Tags != nil {
						if igws == "" {
							network_response = append(network_response, NetworkResponse{Name: *vpc.Tags[0].Value, VpcId: *vpc.VpcId, Subnets: subnets_list, State: *vpc.State, IsDefault: *vpc.IsDefault, Region: region})
							return network_response, nil
						} else {
							network_response = append(network_response, NetworkResponse{Name: *vpc.Tags[0].Value, VpcId: *vpc.VpcId, Subnets: subnets_list, State: *vpc.State, Igw: igws, IsDefault: *vpc.IsDefault, Region: region})
							return network_response, nil
						}

					} else {
						if igws == "" {
							network_response = append(network_response, NetworkResponse{VpcId: *vpc.VpcId, Subnets: subnets_list, State: *vpc.State, IsDefault: *vpc.IsDefault, Region: region})
						} else {
							network_response = append(network_response, NetworkResponse{VpcId: *vpc.VpcId, Subnets: subnets_list, State: *vpc.State, Igw: igws, IsDefault: *vpc.IsDefault, Region: region})
						}
					}
				}
			}
		}
		return network_response, nil
	}
}
