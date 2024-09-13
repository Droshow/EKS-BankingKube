package main

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/secretsmanager"
)

func getSecret(secretName, region string) string {
	// Create a new AWS session. If there's an error, the program will panic.
	sess := session.Must(session.NewSession())
	// Create a new Secrets Manager client using the AWS session and a config object with the specified region.
	svc := secretsmanager.New(sess, aws.NewConfig().WithRegion(region))
	// Create a new GetSecretValueInput struct. This struct is used as input to the GetSecretValue function.
	input := &secretsmanager.GetSecretValueInput{
		SecretId: aws.String(secretName),
	}
	// Call the GetSecretValue function with the input struct. This function retrieves the secret value from Secrets Manager.
	result, err := svc.GetSecretValue(input)
	// Check if there was an error retrieving the secret. If there was an error, print the error and return an empty string.
	if err != nil {
		fmt.Println(err)
		return ""
	}

	return *result.SecretString
}
