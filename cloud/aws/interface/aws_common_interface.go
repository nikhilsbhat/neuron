package neuronaws

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
)

type AwsCommonInput struct {
	AvailabilityZone string `json:"AvailabilityZone,omitempty"`
}

type CreateTagsInput struct {
	Resource string
	Name     string
	Value    string
}

func (sess *EstablishedSession) DescribeAllAvailabilityZones(a *AwsCommonInput) (*ec2.DescribeAvailabilityZonesOutput, error) {

	if sess.Ec2 != nil {
		input := &ec2.DescribeAvailabilityZonesInput{}
		result, err := (sess.Ec2).DescribeAvailabilityZones(input)

		if err != nil {
			return nil, err
		}
		return result, nil
	}
	return nil, fmt.Errorf("Did not get session to perform action, cannot proceed further")

}

func (sess *EstablishedSession) DescribeAvailabilityZones(a *AwsCommonInput) (*ec2.DescribeAvailabilityZonesOutput, error) {

	if sess.Ec2 != nil {
		if a.AvailabilityZone != "" {
			input := &ec2.DescribeAvailabilityZonesInput{
				ZoneNames: aws.StringSlice([]string{a.AvailabilityZone}),
			}
			result, err := (sess.Ec2).DescribeAvailabilityZones(input)

			if err != nil {
				return nil, err
			}
			return result, nil
		}
		return nil, fmt.Errorf("Did not get session to perform action, cannot proceed further")
	}
	return nil, fmt.Errorf("Did not get session to perform action, cannot proceed further")

}

func (sess *EstablishedSession) CreateTags(t *CreateTagsInput) error {

	if sess.Ec2 != nil {
		input := &ec2.CreateTagsInput{
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
		_, err := (sess.Ec2).CreateTags(input)
		if err != nil {
			return err
		}
		return nil
	}
	return fmt.Errorf("Did not get session to perform action, cannot proceed further")
}

func (sess *EstablishedSession) GetRegions() (*ec2.DescribeRegionsOutput, error) {

	if sess.Ec2 != nil {
		input := &ec2.DescribeRegionsInput{}
		result, err := (sess.Ec2).DescribeRegions(input)

		if err != nil {
			return nil, err
		}
		return result, nil
	}
	return nil, fmt.Errorf("Did not get session to perform action, cannot proceed further")
}
