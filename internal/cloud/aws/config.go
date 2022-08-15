package aws

import (
	"context"
	"log"

	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials/stscreds"
	"github.com/aws/aws-sdk-go-v2/service/sts"
)

type AWSProviderConfig struct {
	AwsRoleArn types.String `tfsdk:"role_arn"`
	ExternalId types.String `tfsdk:"external_id"`
	Region     types.String `tfsdk:"region"`
}

type AWSConfig struct {
	provider *AWSProviderConfig
	config   *aws.Config
}

func Initialize(ctx context.Context, providerConfig *AWSProviderConfig) (*AWSConfig, error) {
	awsClient := AWSConfig{}

	log.Printf("[debug] Converting AWS values to config")
	var loadErr error
	cfg, loadErr := config.LoadDefaultConfig(ctx, config.WithRegion(providerConfig.Region.Value))
	if loadErr != nil {
		return nil, loadErr
	}
	stsClient := sts.NewFromConfig(cfg)
	provider := stscreds.NewAssumeRoleProvider(stsClient, providerConfig.AwsRoleArn.Value, func(o *stscreds.AssumeRoleOptions) {
		o.ExternalID = aws.String(providerConfig.ExternalId.Value)
	})
	cfg.Credentials = aws.NewCredentialsCache(provider)
	log.Printf("[debug] AWS Config Created")

	awsClient.config = &cfg
	awsClient.provider = providerConfig

	return &awsClient, nil
}
