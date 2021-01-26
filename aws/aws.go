package aws

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/ecr"
	"github.com/burizz/ecr-image-replication/config"
)

func LoginToECR() (ecrLoginErr error) {
	ecrClient, ecrInitErr := config.ECRClientInit("us-east-1")
	if ecrInitErr != nil {
		return ecrInitErr
	}

	input := &ecr.GetAuthorizationTokenInput{}

	result, getAuthTokenErr := ecrClient.GetAuthorizationToken(input)
	if getAuthTokenErr != nil {
		if aerr, ok := getAuthTokenErr.(awserr.Error); ok {
			switch aerr.Code() {
			case ecr.ErrCodeServerException:
				fmt.Println(ecr.ErrCodeServerException, aerr.Error())
			case ecr.ErrCodeInvalidParameterException:
				fmt.Println(ecr.ErrCodeInvalidParameterException, aerr.Error())
			default:
				fmt.Println(aerr.Error())
			}
		} else {
			// Print the error, cast err to awserr.Error to get the Code and
			// Message from an error.
			fmt.Println(getAuthTokenErr.Error())
		}
		return nil
	}
	fmt.Println(result)
	return nil
}
