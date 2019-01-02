package awsimage

import (
	"fmt"
	aws "neuron/cloud/aws/interface"
	"strings"
)

// This particular function is tailored to find the Id's of the images, of whom's name is matched with the keyword entered.
func (i *GetImageInput) SearchImage(con aws.EstablishConnectionInput) (ImageResponse, error) {

	ec2, seserr := con.EstablishConnection()
	if seserr != nil {
		return ImageResponse{}, seserr
	}

	result, deserr := ec2.DescribeAllImages(
		&aws.DescribeComputeInput{},
	)

	if deserr != nil {
		return ImageResponse{}, deserr
	}

	image_id := make([]string, 0)
	for _, image := range result.Images {
		if strings.Contains(*image.Name, i.Kind) {
			image_id = append(image_id, *image.ImageId)
		}
	}

	if image_id != nil {
		return ImageResponse{ImageIds: image_id}, nil
	}
	return ImageResponse{}, fmt.Errorf("We were unable to find the image with the keyword you entered")
}

//This function will check if the entered image exists in account for that particular region or not.
func (i *GetImageInput) IsImageAvailable(con aws.EstablishConnectionInput) (bool, error) {

	ec2, seserr := con.EstablishConnection()
	if seserr != nil {
		return false, seserr
	}

	// desribing image to check if image exists
	image_result, image_err := ec2.DescribeImages(
		&aws.DescribeComputeInput{
			ImageIds: i.ImageIds,
		},
	)

	if image_err != nil {
		return false, image_err
	}

	switch images := len(image_result.Images); images {
	case 1:
		return true, nil
	default:
		if images > 1 {
			return false, fmt.Errorf("Oops...!!. found multiple images, something is not right as it has to be")
		}
		return false, fmt.Errorf("Oops...!!. Could find the images you entered, hence not proceedig further.")
	}
}
