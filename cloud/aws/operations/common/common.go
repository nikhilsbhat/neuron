package awscommon

import (
	"neuron/cloud/aws/interface"
	"sort"
	"strconv"
)

type Tag struct {
	Resource string
	Name     string
	Value    string
}

type CommonInput struct {
	AvailabilityZone string   `json:"AvailabilityZone,omitempty"`
	SortInput        []string `json:"SortInput,omitempty"`
}

func (r *CommonInput) GetAvailabilityZones(con neuronaws.EstablishConnectionInput) ([]string, error) {

	ec2, sesserr := con.EstablishConnection()
	if sesserr != nil {
		return nil, sesserr
	}

	result, err := ec2.DescribeAllAvailabilityZones(
		&neuronaws.AwsCommonInput{},
	)
	if err != nil {
		return nil, err
	} else {
		availabilityzones := result.AvailabilityZones
		zones := make([]string, 0)
		for _, zone := range availabilityzones {
			zones = append(zones, *zone.ZoneName)
		}
		return zones, nil
	}
}

func (t *Tag) CreateTags(con neuronaws.EstablishConnectionInput) (string, error) {

	ec2, sesserr := con.EstablishConnection()
	if sesserr != nil {
		return "", sesserr
	}

	err := ec2.CreateTags(
		&neuronaws.CreateTagsInput{
			Resource: t.Resource,
			Name:     t.Name,
			Value:    t.Value,
		})
	if err != nil {
		return "", err
	}
	return t.Value, nil
}

func (r *CommonInput) GetRegions(con neuronaws.EstablishConnectionInput) ([]string, error) {

	ec2, sesserr := con.EstablishConnection()
	if sesserr != nil {
		return nil, sesserr
	}

	result, err := ec2.GetRegions()
	if err != nil {
		return nil, err
	}
	regions := make([]string, 0)
	for _, region := range result.Regions {
		regions = append(regions, *region.RegionName)
	}
	return regions, nil
}

func (r *CommonInput) GetRegionFromAvail(con neuronaws.EstablishConnectionInput) (string, error) {

	ec2, sesserr := con.EstablishConnection()
	if sesserr != nil {
		return "", sesserr
	}

	result, err := ec2.DescribeAvailabilityZones(
		&neuronaws.AwsCommonInput{
			AvailabilityZone: r.AvailabilityZone,
		},
	)

	if err != nil {
		return "", err
	}
	return *result.AvailabilityZones[0].RegionName, nil
}

func (r *CommonInput) GetUniqueNumberFromTags() (int, error) {

	// Sort by name, preserving original order
	sort.SliceStable(r.SortInput, func(i, j int) bool { return r.SortInput[i] < r.SortInput[j] })
	if len(r.SortInput) == 0 {
		return 0, nil
	}
	lastchr := r.SortInput[len(r.SortInput)-1]
	uniq, err := strconv.Atoi(string(lastchr[len(lastchr)-1]))
	if err != nil {
		return 0, err
	}
	return (uniq + 1), nil
}
