package main

import (
	"context"
	"fmt"
	gocardless "github.com/gocardless/gocardless-pro-go"
)

func main() {
	// Use the getSecret function to retrieve the GoCardless access token from AWS Secrets Manager
	secretName := "go_cardless_sandbox_token"
	region := "eu-central-1"
	accessToken := getSecret(secretName, region)

	// Check if accessToken is empty
	if accessToken == "" {
		fmt.Println("Failed to retrieve access token")
		return
	}

	// Initialize the GoCardless client with the retrieved access token
	opts := gocardless.WithEndpoint(gocardless.SandboxEndpoint)
	client, err := gocardless.New(accessToken, opts)
	if err != nil {
		fmt.Println("Error creating GoCardless client:", err)
		return
	}

	// Make your first API request to list customers
	ctx := context.TODO()
	customerListParams := gocardless.CustomerListParams{}
	customers, err := client.Customers.List(ctx, customerListParams)
	if err != nil {
		fmt.Println("Error listing customers:", err)
		return
	}

	fmt.Println(customers)
}
