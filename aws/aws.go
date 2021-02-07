package aws

import (
	"fmt"

	log "github.com/sirupsen/logrus"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/ecr"
	"github.com/burizz/ecr-image-replication/config"
)

// GetECRAuthToken - retrieves the authentication from ECR token to be used by docker login
func GetECRAuthToken(awsRegion string) (authToken string, getAuthTokenErr error) {
	ecrClient, ecrClientInitErr := config.ECRClientInit(awsRegion)
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
		return "", getAuthTokenErr
	}
	authTokenValue := result.AuthorizationData[0].AuthorizationToken
	return *authTokenValue, nil
}

// CreateECRRepo - create ECR repository; if it already exists skip
func CreateECRRepo(awsRegion string, ecrRepoName string) (ecrCreateRepo error) {
	ecrClient, ecrClientInitErr := config.ECRClientInit(awsRegion)
	if ecrClientInitErr != nil {
		return ecrClientInitErr
	}

	input := &ecr.CreateRepositoryInput{
		RepositoryName: aws.String(ecrRepoName),
	}

	_, createRepoErr := ecrClient.CreateRepository(input)
	if createRepoErr != nil {
		if aerr, ok := createRepoErr.(awserr.Error); ok {
			switch aerr.Code() {
			case ecr.ErrCodeServerException:
				fmt.Println(ecr.ErrCodeServerException, aerr.Error())
			case ecr.ErrCodeInvalidParameterException:
				fmt.Println(ecr.ErrCodeInvalidParameterException, aerr.Error())
			case ecr.ErrCodeInvalidTagParameterException:
				fmt.Println(ecr.ErrCodeInvalidTagParameterException, aerr.Error())
			case ecr.ErrCodeTooManyTagsException:
				fmt.Println(ecr.ErrCodeTooManyTagsException, aerr.Error())
			case ecr.ErrCodeRepositoryAlreadyExistsException:
				log.Infof("Skip create - ECR repo [%v] already exists", ecrRepoName)
				return nil
			case ecr.ErrCodeLimitExceededException:
				fmt.Println(ecr.ErrCodeLimitExceededException, aerr.Error())
			case ecr.ErrCodeKmsException:
				fmt.Println(ecr.ErrCodeKmsException, aerr.Error())
			default:
				fmt.Println(aerr.Error())
			}
		} else {
			// Print the error, cast err to awserr.Error to get the Code and
			// Message from an error.
			fmt.Println(createRepoErr.Error())
		}
		return createRepoErr
	}

	log.Infof("successfully created ECR repo: %v", ecrRepoName)
	return nil
}
