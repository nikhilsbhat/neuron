package neuronaws

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/ec2"
	err "github.com/nikhilsbhat/neuron/error"
	"reflect"
	"time"
)

type Filters struct {
	Name  string
	Value []string
}

type CreateNetworkInput struct {
	Cidr            string `json:"Cidr,omitempty"`
	Tenancy         string `json:"Tenancy,omitempty"`
	VpcId           string `json:"VpcId,omitempty"`
	SubId           string `json:"SubId,omitempty"`
	IgwId           string `json:"IgwId,omitempty"`
	Zone            string `json:"Zone,omitempty"`
	Name            string `json:"Name,omitempty"`
	DestinationCidr string `json:"DestinationCidr,omitempty"`
	RouteTableId    string `json:"RouteTableId,omitempty"`
}

type IngressEgressInput struct {
	Port  int64  `json:"Port,omitempty"`
	SecId string `json:"SecId,omitempty"`
}

type DescribeNetworkInput struct {
	SecIds         []string `json:"SecIds,omitempty"`
	IgwIds         []string `json:"IgwId,omitempty"`
	VpcIds         []string `json:"VpcId,omitempty"`
	SubnetIds      []string `json:"SubnetId,omitempty"`
	RouteTableIds  []string `json:"RouteTableIds,omitempty"`
	Filters        Filters  `json:"Filters,omitempty"`
	AssociationsId string   `json:"AssociationsId,omitempty"`
}

func (sess *EstablishedSession) CreateVpc(v *CreateNetworkInput) (*ec2.CreateVpcOutput, error) {

	if sess.Ec2 != nil {
		if (v.Cidr != "") || (v.Tenancy != "") {
			input := &ec2.CreateVpcInput{
				CidrBlock:       aws.String(v.Cidr),
				InstanceTenancy: aws.String(v.Tenancy),
			}
			result, err := (sess.Ec2).CreateVpc(input)

			if err != nil {
				return nil, err
			}
			return result, nil
		}
		return nil, fmt.Errorf(fmt.Sprintf("%v CreateVpc", err.EmptyStructError()))
	}
	return nil, err.InvalidSession()
}

func (sess *EstablishedSession) CreateSubnet(s *CreateNetworkInput) (*ec2.CreateSubnetOutput, error) {

	if sess.Ec2 != nil {
		if (s.Cidr != "") || (s.VpcId != "") || (s.Zone != "") {
			input := &ec2.CreateSubnetInput{
				CidrBlock:        aws.String(s.Cidr),
				VpcId:            aws.String(s.VpcId),
				AvailabilityZone: aws.String(s.Zone),
			}
			result, err := (sess.Ec2).CreateSubnet(input)

			if err != nil {
				return nil, err
			}
			return result, nil
		}
		return nil, fmt.Errorf(fmt.Sprintf("%v CreateSubnet", err.EmptyStructError()))
	}
	return nil, err.InvalidSession()
}

func (sess *EstablishedSession) CreateIgw() (*ec2.CreateInternetGatewayOutput, error) {

	if sess.Ec2 != nil {
		input := &ec2.CreateInternetGatewayInput{}
		result, err := (sess.Ec2).CreateInternetGateway(input)

		if err != nil {
			return nil, err
		}
		return result, nil
	}
	return nil, err.InvalidSession()
}

//never pass arrays to attach the intergateways, it never works
func (sess *EstablishedSession) AttachIgw(a *DescribeNetworkInput) error {

	if sess.Ec2 != nil {
		if (a.IgwIds != nil) || (a.VpcIds != nil) {
			input := &ec2.AttachInternetGatewayInput{
				InternetGatewayId: aws.String(a.IgwIds[0]),
				VpcId:             aws.String(a.VpcIds[0]),
			}
			_, err := (sess.Ec2).AttachInternetGateway(input)

			if err != nil {
				return err
			}
			return nil
		}
		return fmt.Errorf(fmt.Sprintf("%v AttachIgw", err.EmptyStructError()))
	}
	return err.InvalidSession()
}

func (sess *EstablishedSession) CreateSecurityGroup(s *CreateNetworkInput) (*ec2.CreateSecurityGroupOutput, error) {

	if sess.Ec2 != nil {
		if s.VpcId != "" {
			input := &ec2.CreateSecurityGroupInput{
				Description: aws.String("This security group is created by Neuron api"),
				VpcId:       aws.String(s.VpcId),
				GroupName:   aws.String(s.Name),
			}
			result, err := (sess.Ec2).CreateSecurityGroup(input)

			if err != nil {
				return nil, err
			}
			return result, nil

		}
		return nil, fmt.Errorf(fmt.Sprintf("%v CreateSecurityGroup", err.EmptyStructError()))
	}
	return nil, err.InvalidSession()

}

func (sess *EstablishedSession) CreateRouteTable(r *CreateNetworkInput) (*ec2.CreateRouteTableOutput, error) {

	if sess.Ec2 != nil {
		if r.VpcId != "" {
			input := &ec2.CreateRouteTableInput{
				VpcId: aws.String(r.VpcId),
			}
			result, err := (sess.Ec2).CreateRouteTable(input)

			if err != nil {
				return nil, err
			}
			return result, nil
		}
		return nil, fmt.Errorf(fmt.Sprintf("%v CreateRouteTable", err.EmptyStructError()))
	}
	return nil, err.InvalidSession()
}

func (sess *EstablishedSession) WriteRoute(r *CreateNetworkInput) error {

	if sess.Ec2 != nil {
		if r.RouteTableId != "" {
			input := &ec2.CreateRouteInput{
				DestinationCidrBlock: aws.String(r.DestinationCidr),
				GatewayId:            aws.String(r.IgwId),
				RouteTableId:         aws.String(r.RouteTableId),
			}
			_, err := (sess.Ec2).CreateRoute(input)
			if err != nil {
				return err
			}
			return nil
		}
		return fmt.Errorf(fmt.Sprintf("%v WriteRoute", err.EmptyStructError()))
	}
	return err.InvalidSession()

}

func (sess *EstablishedSession) AttachRouteTable(r *CreateNetworkInput) error {

	if sess.Ec2 != nil {
		if r.RouteTableId != "" {
			input := &ec2.AssociateRouteTableInput{
				RouteTableId: aws.String(r.RouteTableId),
				SubnetId:     aws.String(r.SubId),
			}
			_, err := (sess.Ec2).AssociateRouteTable(input)
			if err != nil {
				return err
			}
			return nil
		}
		return fmt.Errorf(fmt.Sprintf("%v AttachRouteTable", err.EmptyStructError()))
	}
	return err.InvalidSession()
}

func (sess *EstablishedSession) DettachRouteTable(r *DescribeNetworkInput) error {

	if sess.Ec2 != nil {
		if r.AssociationsId != "" {
			input := &ec2.DisassociateRouteTableInput{
				AssociationId: aws.String(r.AssociationsId),
			}
			_, err := (sess.Ec2).DisassociateRouteTable(input)
			if err != nil {
				return err
			}
			return nil
		}
		return fmt.Errorf(fmt.Sprintf("%v AttachRouteTable", err.EmptyStructError()))
	}
	return err.InvalidSession()
}

func (sess *EstablishedSession) CreateEgressRule(i *IngressEgressInput) error {

	if sess.Ec2 != nil {
		security_ingress_input := &ec2.AuthorizeSecurityGroupEgressInput{
			//    FromPort   : aws.int64(from_port),
			GroupId: aws.String(i.SecId),
			//    ToPort     : aws.int64(to_port),
			IpPermissions: []*ec2.IpPermission{
				{
					IpProtocol: aws.String("-1"),
					IpRanges: []*ec2.IpRange{
						{
							//CidrIp: aws.String("0.0.0.0/0"),
						},
					},
				},
			},
		}
		_, egress_err := (sess.Ec2).AuthorizeSecurityGroupEgress(security_ingress_input)

		if egress_err != nil {
			return egress_err
		}
		return nil
	}
	return err.InvalidSession()

}

func (sess *EstablishedSession) CreateIngressRule(i *IngressEgressInput) error {

	if sess.Ec2 != nil {
		security_ingress_input := &ec2.AuthorizeSecurityGroupIngressInput{
			FromPort:   aws.Int64(i.Port),
			IpProtocol: aws.String("tcp"),
			GroupId:    aws.String(i.SecId),
			ToPort:     aws.Int64(i.Port),
			CidrIp:     aws.String("0.0.0.0/0"),
		}
		_, ingress_err := (sess.Ec2).AuthorizeSecurityGroupIngress(security_ingress_input)

		if ingress_err != nil {
			return ingress_err
		}
		return nil
	}
	return err.InvalidSession()

}

func (sess *EstablishedSession) DeleteIgw(i *DescribeNetworkInput) error {

	if sess.Ec2 != nil {
		if i.IgwIds != nil {

			input := &ec2.DeleteInternetGatewayInput{
				InternetGatewayId: aws.String(i.IgwIds[0]),
			}
			_, err := (sess.Ec2).DeleteInternetGateway(input)
			if err != nil {
				return err
			}
			return nil
		}
		return fmt.Errorf(fmt.Sprintf("%v DeleteIgw", err.EmptyStructError()))
	}
	return err.InvalidSession()

}

//never pass arrays to dettach the intergateways, it never works
func (sess *EstablishedSession) DetachIgw(i *DescribeNetworkInput) error {

	if sess.Ec2 != nil {
		if i.IgwIds != nil {
			input := &ec2.DetachInternetGatewayInput{
				InternetGatewayId: aws.String(i.IgwIds[0]),
				VpcId:             aws.String(i.VpcIds[0]),
			}
			_, err := (sess.Ec2).DetachInternetGateway(input)
			if err != nil {
				return err
			}
			return nil
		}
		return fmt.Errorf(fmt.Sprintf("%v DetachIgw", err.EmptyStructError()))
	}
	return err.InvalidSession()

}

func (sess *EstablishedSession) DescribeIgw(d *DescribeNetworkInput) (*ec2.DescribeInternetGatewaysOutput, error) {

	if sess.Ec2 != nil {
		if d.IgwIds != nil {
			input := &ec2.DescribeInternetGatewaysInput{
				InternetGatewayIds: aws.StringSlice(d.IgwIds),
			}
			result, err := (sess.Ec2).DescribeInternetGateways(input)

			if err != nil {
				return nil, err
			}
			return result, nil
		}

		if reflect.DeepEqual(d.Filters, Filters{}) {
			return nil, fmt.Errorf(fmt.Sprintf("%v DescribeIgw. You selected filters for this yet you passed empty", err.EmptyStructError()))
		}
		input := &ec2.DescribeInternetGatewaysInput{
			Filters: []*ec2.Filter{
				{
					Name:   aws.String(d.Filters.Name),
					Values: aws.StringSlice(d.Filters.Value),
				},
			},
		}
		result, err := (sess.Ec2).DescribeInternetGateways(input)

		if err != nil {
			return nil, err
		}
		return result, nil

	}
	return nil, err.InvalidSession()
}

func (sess *EstablishedSession) DescribeAllIgw(d *DescribeNetworkInput) (*ec2.DescribeInternetGatewaysOutput, error) {

	if sess.Ec2 != nil {
		input := &ec2.DescribeInternetGatewaysInput{}
		result, err := (sess.Ec2).DescribeInternetGateways(input)

		if err != nil {
			return nil, err
		}
		return result, nil
	}
	return nil, err.InvalidSession()
}

func (sess *EstablishedSession) DeleteRouteTable(r *DescribeNetworkInput) error {

	if sess.Ec2 != nil {
		if r.RouteTableIds != nil {
			input := &ec2.DeleteRouteTableInput{
				RouteTableId: aws.String(r.RouteTableIds[0]),
			}
			_, err := (sess.Ec2).DeleteRouteTable(input)
			if err != nil {
				return err
			}
			return nil

		}
		return fmt.Errorf(fmt.Sprintf("%v DeleteRouteTable", err.EmptyStructError()))
	}
	return err.InvalidSession()
}

func (sess *EstablishedSession) DescribeRouteTable(d *DescribeNetworkInput) (*ec2.DescribeRouteTablesOutput, error) {

	if sess.Ec2 != nil {
		if d.RouteTableIds != nil {
			input := &ec2.DescribeRouteTablesInput{
				RouteTableIds: aws.StringSlice(d.RouteTableIds),
			}
			result, err := (sess.Ec2).DescribeRouteTables(input)
			if err != nil {
				return nil, err
			}
			return result, nil
		}

		if reflect.DeepEqual(d.Filters, Filters{}) {
			return nil, fmt.Errorf(fmt.Sprintf("%v DescribeRouteTable. You selected filters for this yet you passed empty", err.EmptyStructError()))
		}
		input := &ec2.DescribeRouteTablesInput{
			Filters: []*ec2.Filter{
				{
					Name:   aws.String(d.Filters.Name),
					Values: aws.StringSlice(d.Filters.Value),
				},
			},
		}
		result, err := (sess.Ec2).DescribeRouteTables(input)

		if err != nil {
			return nil, err
		}
		return result, nil
	}
	return nil, err.InvalidSession()
}

func (sess *EstablishedSession) DeleteSecurityGroup(s *DescribeNetworkInput) error {

	if sess.Ec2 != nil {
		if s.SecIds != nil {
			for _, sec := range s.SecIds {
				input := &ec2.DeleteSecurityGroupInput{
					GroupId: aws.String(sec),
				}
				_, err := (sess.Ec2).DeleteSecurityGroup(input)
				if err != nil {
					return err
				}
				return nil
			}
		}

		return fmt.Errorf(fmt.Sprintf("%v DeleteSecurityGroup", err.EmptyStructError()))
	}
	return err.InvalidSession()

}

func (sess *EstablishedSession) DescribeSecurityGroup(d *DescribeNetworkInput) (*ec2.DescribeSecurityGroupsOutput, error) {

	if sess.Ec2 != nil {
		if d.SecIds != nil {
			input := &ec2.DescribeSecurityGroupsInput{
				GroupIds: aws.StringSlice(d.SecIds),
			}
			result, err := (sess.Ec2).DescribeSecurityGroups(input)
			if err != nil {
				return nil, err
			}
			return result, nil

		}

		if reflect.DeepEqual(d.Filters, Filters{}) {
			return nil, fmt.Errorf(fmt.Sprintf("%v DescribeSecurityGroup. You selected filters for this yet you passed empty", err.EmptyStructError()))
		}
		input := &ec2.DescribeSecurityGroupsInput{
			Filters: []*ec2.Filter{
				{
					Name:   aws.String(d.Filters.Name),
					Values: aws.StringSlice(d.Filters.Value),
				},
			},
		}
		result, err := (sess.Ec2).DescribeSecurityGroups(input)
		if err != nil {
			return nil, err
		}
		return result, nil
	}
	return nil, err.InvalidSession()

}

func (sess *EstablishedSession) DescribeAllSecurityGroup(d *DescribeNetworkInput) (*ec2.DescribeSecurityGroupsOutput, error) {

	if sess.Ec2 != nil {
		input := &ec2.DescribeSecurityGroupsInput{}
		result, err := (sess.Ec2).DescribeSecurityGroups(input)
		if err != nil {
			return nil, err
		}
		return result, nil
	}
	return nil, err.InvalidSession()
}

func (sess *EstablishedSession) DescribeAllSubnet(d *DescribeNetworkInput) (*ec2.DescribeSubnetsOutput, error) {

	if sess.Ec2 != nil {
		input := &ec2.DescribeSubnetsInput{}
		result, err := (sess.Ec2).DescribeSubnets(input)
		if err != nil {
			return nil, err
		}
		return result, nil
	}
	return nil, err.InvalidSession()
}

func (sess *EstablishedSession) DescribeSubnet(d *DescribeNetworkInput) (*ec2.DescribeSubnetsOutput, error) {

	if sess.Ec2 != nil {
		if d.SubnetIds != nil {
			input := &ec2.DescribeSubnetsInput{
				SubnetIds: aws.StringSlice(d.SubnetIds),
			}
			result, err := (sess.Ec2).DescribeSubnets(input)
			if err != nil {
				return nil, err
			}
			return result, nil
		}

		if reflect.DeepEqual(d.Filters, Filters{}) {
			return nil, fmt.Errorf(fmt.Sprintf("%v DescribeSubnet. You selected filters for this yet you passed empty", err.EmptyStructError()))
		}
		input := &ec2.DescribeSubnetsInput{
			Filters: []*ec2.Filter{
				{
					Name:   aws.String(d.Filters.Name),
					Values: aws.StringSlice(d.Filters.Value),
				},
			},
		}
		result, err := (sess.Ec2).DescribeSubnets(input)
		if err != nil {
			return nil, err
		}
		return result, nil
	}
	return nil, err.InvalidSession()

}

func (sess *EstablishedSession) DeleteSubnet(d *DescribeNetworkInput) error {

	if sess.Ec2 != nil {
		if d.SubnetIds != nil {
			input := &ec2.DeleteSubnetInput{
				SubnetId: aws.String(d.SubnetIds[0]),
			}
			_, err := (sess.Ec2).DeleteSubnet(input)
			if err != nil {
				return err
			}
			return nil
		}
		return fmt.Errorf(fmt.Sprintf("%v DeleteSubnet", err.EmptyStructError()))
	}
	return err.InvalidSession()
}

func (sess *EstablishedSession) DescribeAllVpc(d *DescribeNetworkInput) (*ec2.DescribeVpcsOutput, error) {

	if sess.Ec2 != nil {
		input := &ec2.DescribeVpcsInput{}
		result, err := (sess.Ec2).DescribeVpcs(input)
		if err != nil {
			return nil, err
		}
		return result, nil
	}
	return nil, err.InvalidSession()
}

func (sess *EstablishedSession) DescribeVpc(d *DescribeNetworkInput) (*ec2.DescribeVpcsOutput, error) {

	if sess.Ec2 != nil {
		if d.VpcIds != nil {
			input := &ec2.DescribeVpcsInput{
				VpcIds: aws.StringSlice(d.VpcIds),
			}
			result, err := (sess.Ec2).DescribeVpcs(input)
			if err != nil {
				return nil, err
			}
			return result, nil
		}

		if reflect.DeepEqual(d.Filters, Filters{}) {
			return nil, fmt.Errorf(fmt.Sprintf("%v DescribeVpc. You selected filters for this yet you passed empty", err.EmptyStructError()))
		}
		input := &ec2.DescribeVpcsInput{
			Filters: []*ec2.Filter{
				{
					Name:   aws.String(d.Filters.Name),
					Values: aws.StringSlice(d.Filters.Value),
				},
			},
		}
		result, err := (sess.Ec2).DescribeVpcs(input)
		if err != nil {
			return nil, err
		}
		return result, nil
	}
	return nil, err.InvalidSession()
}

func (sess *EstablishedSession) DeleteVpc(d *DescribeNetworkInput) error {

	if sess.Ec2 != nil {
		if d.VpcIds != nil {
			for _, vpc := range d.VpcIds {
				input := &ec2.DeleteVpcInput{
					VpcId: aws.String(vpc),
				}
				_, err := (sess.Ec2).DeleteVpc(input)
				if err != nil {
					return err
				}
				return nil
			}
		}
		return fmt.Errorf(fmt.Sprintf("%v DeleteVpc", err.EmptyStructError()))
	}
	return err.InvalidSession()
}

func (sess *EstablishedSession) WaitTillVpcAvailable(d *DescribeNetworkInput) error {

	if sess.Ec2 != nil {

		if reflect.DeepEqual(d.Filters, Filters{}) {
			return fmt.Errorf(fmt.Sprintf("%v WaitTillVpcAvailable", err.EmptyStructError()))
		}
		input := &ec2.DescribeVpcsInput{
			Filters: []*ec2.Filter{
				&ec2.Filter{
					Name:   aws.String(d.Filters.Name),
					Values: aws.StringSlice(d.Filters.Value),
				},
			},
		}
		err := (sess.Ec2).WaitUntilVpcAvailable(input)
		if err != nil {
			return err
		}
		return nil
	}
	return err.InvalidSession()
}

func (sess *EstablishedSession) WaitTillSubnetAvailable(d *DescribeNetworkInput) error {

	if sess.Ec2 != nil {
		if reflect.DeepEqual(d.Filters, Filters{}) {
			return fmt.Errorf(fmt.Sprintf("%v WaitTillSubnetAvailable", err.EmptyStructError()))
		}
		input := &ec2.DescribeSubnetsInput{
			Filters: []*ec2.Filter{
				&ec2.Filter{
					Name:   aws.String(d.Filters.Name),
					Values: aws.StringSlice(d.Filters.Value),
				},
			},
		}
		err := (sess.Ec2).WaitUntilSubnetAvailable(input)
		if err != nil {
			return err
		}
		return nil
	}
	return err.InvalidSession()
}

func (sess *EstablishedSession) WaitUntilSubnetDeleted(d *DescribeNetworkInput) (bool, error) {

	if sess.Ec2 != nil {
		if d.SubnetIds != nil {
			input := &ec2.DescribeSubnetsInput{
				SubnetIds: aws.StringSlice(d.SubnetIds),
			}

			response, deserr := (sess.Ec2).DescribeSubnets(input)
			if response.Subnets != nil {
				start := time.Now()
				for len(response.Subnets) > 0 {
					response, deserr = (sess.Ec2).DescribeSubnets(input)
					if deserr != nil {
						switch deserr.(awserr.Error).Code() {
						case "InvalidSubnetID.NotFound":
							return true, nil
						default:
							return false, deserr
						}
					}
					if time.Since(start) > time.Duration(10*time.Second) {
						return false, fmt.Errorf("Time Out .Oops...!! it took annoyingly more than anticipated time while waiting for subnet to get deleted. Guess I was not called after delete subnet function")
					}
				}
			}
			if deserr != nil {
				switch deserr.(awserr.Error).Code() {
				case "InvalidSubnetID.NotFound":
					return true, nil
				default:
					return false, deserr
				}
			}
			return false, fmt.Errorf("Error occured while waiting for the subnet deletion")
		}
		return false, fmt.Errorf(fmt.Sprintf("%v WaitUntilSubnetDeleted", err.EmptyStructError()))
	}
	return false, err.InvalidSession()
}

func (sess *EstablishedSession) WaitUntilRoutTableDeleted(d *DescribeNetworkInput) (bool, error) {

	if sess.Ec2 != nil {
		if d.RouteTableIds != nil {
			input := &ec2.DescribeRouteTablesInput{
				RouteTableIds: aws.StringSlice(d.RouteTableIds),
			}

			response, deserr := (sess.Ec2).DescribeRouteTables(input)
			if response.RouteTables != nil {
				start := time.Now()
				for len(response.RouteTables) > 0 {
					response, deserr = (sess.Ec2).DescribeRouteTables(input)
					if deserr != nil {
						switch deserr.(awserr.Error).Code() {
						case "InvalidRouteTableID.NotFound":
							return true, nil
						default:
							return false, deserr
						}
					}
					if time.Since(start) > time.Duration(10*time.Second) {
						return false, fmt.Errorf("Time Out .Oops...!! it took annoyingly more than anticipated time while waiting for routetable to get deleted. Guess I was not called after delete routetable function")
					}
				}
			}
			if deserr != nil {
				switch deserr.(awserr.Error).Code() {
				case "InvalidRouteTableID.NotFound":
					return true, nil
				default:
					return false, deserr
				}
			}
			return false, fmt.Errorf("Error occured while waiting for the routetable deletion")
		}
		return false, fmt.Errorf(fmt.Sprintf("%v WaitUntilRoutTableDeleted", err.EmptyStructError()))
	}
	return false, err.InvalidSession()
}

func (sess *EstablishedSession) WaitUntilIgwDeleted(d *DescribeNetworkInput) (bool, error) {

	if sess.Ec2 != nil {
		if d.IgwIds != nil {
			input := &ec2.DescribeInternetGatewaysInput{
				InternetGatewayIds: aws.StringSlice(d.IgwIds),
			}

			response, deserr := (sess.Ec2).DescribeInternetGateways(input)
			if response.InternetGateways != nil {
				start := time.Now()
				for len(response.InternetGateways) > 0 {
					response, deserr = (sess.Ec2).DescribeInternetGateways(input)
					if deserr != nil {
						switch deserr.(awserr.Error).Code() {
						case "InvalidInternetGatewayID.NotFound":
							return true, nil
						default:
							return false, deserr
						}
					}
					if time.Since(start) > time.Duration(10*time.Second) {
						return false, fmt.Errorf("Time Out .Oops...!! it took annoyingly more than anticipated time while waiting for igw to get deleted. Guess I was not called after delete igw function")
					}
				}
			}
			if deserr != nil {
				switch deserr.(awserr.Error).Code() {
				case "InvalidInternetGatewayID.NotFound":
					return true, nil
				default:
					return false, deserr
				}
			}
			return false, fmt.Errorf("Error occured while waiting for the InternetGateways deletion")
		}
		return false, fmt.Errorf(fmt.Sprintf("%v WaitUntilIgwDeleted", err.EmptyStructError()))
	}
	return false, err.InvalidSession()
}
