package kt

/**
Helper functions to make working with AWS easier.
*/

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/secretsmanager"
)

var (
	secSev *secretsmanager.SecretsManager
)

const (
	AwsSmPrefix  = "aws.sm."
	AWSErrPrefix = "AwsError: "
)

func loadViaAWSSecrets(key string) string {
	if secSev == nil {
		sess := session.Must(session.NewSession())
		secSev = secretsmanager.New(sess)
	}

	input := &secretsmanager.GetSecretValueInput{
		SecretId: aws.String(key),
	}
	result, err := secSev.GetSecretValue(input)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case secretsmanager.ErrCodeResourceNotFoundException:
				return AWSErrPrefix + secretsmanager.ErrCodeResourceNotFoundException
			case secretsmanager.ErrCodeInvalidParameterException:
				return AWSErrPrefix + secretsmanager.ErrCodeInvalidParameterException
			case secretsmanager.ErrCodeInvalidRequestException:
				return AWSErrPrefix + secretsmanager.ErrCodeInvalidRequestException
			case secretsmanager.ErrCodeDecryptionFailure:
				return AWSErrPrefix + secretsmanager.ErrCodeDecryptionFailure
			case secretsmanager.ErrCodeInternalServiceError:
				return AWSErrPrefix + secretsmanager.ErrCodeInternalServiceError
			default:
				return AWSErrPrefix + aerr.Error()
			}
		}
	}
	if result.SecretString != nil {
		return *result.SecretString
	}
	return ""
}
