package aws

import (
	"context"
	"log"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials/stscreds"
	"github.com/aws/aws-sdk-go-v2/service/sts"
)

type AWSConfig struct {
	awsRoleArn string
	externalId string
	region     string
	config     *aws.Config
}

func Initialize(ctx context.Context, awsMap map[string]interface{}) (*AWSConfig, error) {
	awsClient := AWSConfig{}

	if roleArn, ok := awsMap["role_arn"].(string); ok && roleArn != "" {
		awsClient.awsRoleArn = roleArn
	}
	if externalId, ok := awsMap["external_id"].(string); ok && externalId != "" {
		awsClient.externalId = externalId
	}
	if region, ok := awsMap["region"].(string); ok && region != "" {
		awsClient.region = region
	}

	log.Printf("[debug] Converting AWS values to config")
	var loadErr error
	cfg, loadErr := config.LoadDefaultConfig(ctx, config.WithRegion(awsClient.region))
	if loadErr != nil {
		return nil, loadErr
	}
	stsClient := sts.NewFromConfig(cfg)
	provider := stscreds.NewAssumeRoleProvider(stsClient, awsClient.awsRoleArn, func(o *stscreds.AssumeRoleOptions) {
		o.ExternalID = aws.String(awsClient.externalId)
	})
	cfg.Credentials = aws.NewCredentialsCache(provider)
	log.Printf("[debug] AWS Config Created")

	awsClient.config = &cfg

	return &awsClient, nil
}
