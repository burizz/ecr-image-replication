package aws

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/ecr"
	"github.com/burizz/ecr-image-replication/config"
)

// GetECRAuthToken - retrieves the authentication from ECR token to be used by docker login
func GetECRAuthToken() (authToken string, getAuthTokenErr error) {
	ecrClient, ecrClientInitErr := config.ECRClientInit("us-east-1")
	if ecrClientInitErr != nil {
		return "", ecrClientInitErr
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
		return "", nil
	}
	authTokenValue := result.AuthorizationData[0].AuthorizationToken
	return *authTokenValue, nil
}
