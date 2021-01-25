package config

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ecr"
)

func ECRClientInit(awsRegion string) (awsInitErr error) {
	//TODO: Finish this
	sess, awsSessErr := session.NewSession(&aws.Config{
		Region: aws.String(awsRegion)},
	)
	if awsSessErr != nil {
		return nil, fmt.Errorf("cannot initialize aws session", awsSessErr)
	}

	svc := ecr.New(sess)
	return
}
