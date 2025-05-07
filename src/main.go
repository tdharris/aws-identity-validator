package main

import (
	"context"
	"fmt"
	"os"

	// AWS SDK v1
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sts"

	// AWS SDK v2
	"github.com/aws/aws-sdk-go-v2/config"
	stsv2 "github.com/aws/aws-sdk-go-v2/service/sts"
)

func main() {
	fmt.Println("AWS Identity Validator")
	fmt.Println("===========================")

	// Check IRSA environment
	checkIRSAEnvironment()

	// AWS SDK v1
	fmt.Println("\n[AWS SDK v1]")
	getIdentityV1()

	// AWS SDK v2
	fmt.Println("\n[AWS SDK v2]")
	getIdentityV2()
}

// checkIRSAEnvironment checks and reports on the IRSA configuration in the pod
func checkIRSAEnvironment() {
	fmt.Println("\n[IRSA Environment Check]")

	// Check for the AWS Web Identity token file path
	tokenPath := os.Getenv("AWS_WEB_IDENTITY_TOKEN_FILE")
	if tokenPath != "" {
		fmt.Printf("✓ AWS_WEB_IDENTITY_TOKEN_FILE env var is set: %s\n", tokenPath)
		_, err := os.Stat(tokenPath)
		if err == nil {
			fmt.Println("✓ Token file exists")

			// Get file size
			fileInfo, err := os.Stat(tokenPath)
			if err == nil {
				fmt.Printf("✓ Token file size: %d bytes\n", fileInfo.Size())
			}
		} else {
			fmt.Printf("✗ Token file issue: %v\n", err)
		}
	} else {
		fmt.Println("✗ AWS_WEB_IDENTITY_TOKEN_FILE env var not set - IRSA not configured")
	}

	// Check for role ARN
	roleArn := os.Getenv("AWS_ROLE_ARN")
	if roleArn != "" {
		fmt.Printf("✓ AWS_ROLE_ARN env var is set: %s\n", roleArn)
	} else {
		fmt.Println("✗ AWS_ROLE_ARN env var not set - IRSA not configured")
	}
}

// getIdentityV1 gets and prints AWS identity information using AWS SDK v1
func getIdentityV1() {
	// Create a new AWS session with default credential chain
	sess, err := session.NewSession()
	if err != nil {
		fmt.Printf("Failed to create AWS SDK v1 session: %v\n", err)
		os.Exit(1)
	}

	// Create an STS client
	svc := sts.New(sess)

	// Call GetCallerIdentity
	result, err := svc.GetCallerIdentity(&sts.GetCallerIdentityInput{})
	if err != nil {
		fmt.Printf("Failed to get caller identity with AWS SDK v1: %v\n", err)
		os.Exit(1)
	}

	// Print the results
	fmt.Println("✔  Successfully authenticated with AWS SDK v1")
	fmt.Printf("Account ID: %s\n", *result.Account)
	fmt.Printf("User ID: %s\n", *result.UserId)
	fmt.Printf("ARN: %s\n", *result.Arn)

	// Print more credential information if available
	creds, err := sess.Config.Credentials.Get()
	if err != nil {
		fmt.Printf("Failed to get credential details: %v\n", err)
	} else {
		fmt.Printf("Provider: %s\n", creds.ProviderName)
		// Don't print the actual secret key for security reasons
		fmt.Printf("Access Key ID: %s\n", creds.AccessKeyID)
		if creds.SessionToken != "" {
			fmt.Println("Using temporary credentials (has session token)")
		}
	}
}

// getIdentityV2 gets and prints AWS identity information using AWS SDK v2
func getIdentityV2() {
	ctx := context.Background()

	// Load the SDK's configuration from the default credential chain
	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		fmt.Printf("Failed to load default AWS SDK v2 config: %v\n", err)
		os.Exit(1)
	}

	// Create an STS client
	svc := stsv2.NewFromConfig(cfg)

	// Call GetCallerIdentity
	result, err := svc.GetCallerIdentity(ctx, &stsv2.GetCallerIdentityInput{})
	if err != nil {
		fmt.Printf("Failed to get caller identity with AWS SDK v2: %v\n", err)
		os.Exit(1)
	}

	// Print the results
	fmt.Println("✔  Successfully authenticated with AWS SDK v2")
	fmt.Printf("Account ID: %s\n", *result.Account)
	fmt.Printf("User ID: %s\n", *result.UserId)
	fmt.Printf("ARN: %s\n", *result.Arn)

	// Get credential provider information if available
	credentials, err := cfg.Credentials.Retrieve(ctx)
	if err != nil {
		fmt.Printf("Failed to get credential details: %v\n", err)
	} else {
		fmt.Printf("Provider: %s\n", credentials.Source)
		fmt.Printf("Access Key ID: %s\n", credentials.AccessKeyID)
		if credentials.SessionToken != "" {
			fmt.Println("Using temporary credentials (has session token)")
		}
	}
}
