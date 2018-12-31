package DengineAwsInterface

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/elb"
	"github.com/aws/aws-sdk-go/service/elbv2"
	log "neuron/logger"
	"strings"
)

var (
	Svc  *ec2.EC2
	Elb  *elb.ELB
	Elb2 *elbv2.ELBV2
)

type EstablishConnectionInput struct {
	Region   string
	Resource string
}

type Tag struct {
	Resource string
	Name     string
	Value    string
}

func (con *EstablishConnectionInput) EstablishConnection() {

	sess := session.Must(session.NewSession(&aws.Config{
		Region: aws.String(con.Region),
	}))

	switch strings.ToLower(con.Resource) {
	case "ec2":
		svc := ec2.New(sess)
		Svc = svc
	case "elb":
		Svc = ec2.New(sess)
		Elb = elb.New(sess)
	case "elb2":
		Svc = ec2.New(sess)
		Elb2 = elbv2.New(sess)
	default:
		log.Info("")
		log.Error("I feel we are lost in creating session :S")
		log.Error("Either the keys which you used to establish connection is not valid else you have not chosen right resource")
		log.Info("")
	}

}

// being create_tags my job is to give names to the resouces depending on who called me
func create_tags(t Tag) (string, error) {

	// My task is to create tags for the resources from where I am called
	tag_name := &ec2.CreateTagsInput{
		Resources: []*string{
			aws.String(t.Resource),
		},
		Tags: []*ec2.Tag{
			{
				Key:   aws.String(t.Name),
				Value: aws.String(t.Value),
			},
		},
	}
	_, tag_err := Svc.CreateTags(tag_name)
	if tag_err != nil {
		return "", tag_err
	} else {
		return t.Value, nil
	}
}

//To get all the regions in aws
func GetRegions() ([]string, error) {

	region_input := &ec2.DescribeRegionsInput{}

	result, region_err := Svc.DescribeRegions(region_input)

	if region_err != nil {
		return nil, region_err
	} else {
		var regions []string
		for _, region := range result.Regions {
			regions = append(regions, *region.RegionName)
		}
		return regions, nil
	}
}
