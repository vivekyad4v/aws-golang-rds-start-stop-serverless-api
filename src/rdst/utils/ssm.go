package utils

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ssm"

	"os"
)

// GetParameterValue from AWS SSM Parameter store
func GetParameterValue(name string) (string, error) {
	sess := session.New()
	svc := ssm.New(sess)

	fullIdentifier := "/" + os.Getenv("ORG_ID") + "/" + os.Getenv("ENVIRON") + "/" + os.Getenv("PROJECT_NAME") + "/" + name

	output, err := svc.GetParameter(
		&ssm.GetParameterInput{
			Name:           aws.String(fullIdentifier),
			WithDecryption: aws.Bool(true),
		},
	)

	if err != nil {
		return "", err
	}

	return aws.StringValue(output.Parameter.Value), nil
}
