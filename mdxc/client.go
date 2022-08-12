package mdxc

import (
	"context"
	"terraform-provider-mdxc/mdxc/internal/aws/app_identity"
	"terraform-provider-mdxc/mdxc/internal/massdriver"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials/stscreds"
	"github.com/aws/aws-sdk-go-v2/service/sts"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type awsConfig struct {
	awsRoleArn string
	externalId string
	region     string
}

type azureConfig struct {
	subscriptionId string
	clientId       string
	clientSecret   string
	tenantId       string
}

type gcpConfig struct {
	credentials string
	project     string
}

type MdxcClient struct {
	identityService AppIdentity
	cloud string
}

type AppIdentity interface {
	InitializeClient()
	createApplicationIdentity()
	shouldUpdateApplicationIdentity()
	updateApplicationIdentity()
	deleteApplicationIdentity()
	bindPolicy()
	unbindPolicy()
}
/*
Universal state
{
	"role-123": {
		"Idempotency": {
			"step-1": 0,
		}
		"Identity": {}
	}
	"policy-123": {
		"role": ""
		"policy-name": ""
	}
}
*/

func NewMdxcClient(ctx context.Context, d *schema.ResourceData) *MdxcClient {
	c := &MdxcClient{}
	c.determineCloud(ctx, d)
	return c
}

func (c *MdxcClient) determineCloud(ctx context.Context, d *schema.ResourceData) error {
	// Move to AppIdentity implementors as InitializeClient(ctx, d)?
	if aws, ok := d.Get("aws").([]interface{}); ok && len(aws) > 0 && aws[0] != nil {
		mappedAWSConfig := aws[0].(map[string]interface{})
		awsCfg, err := initializeAWSConfig(ctx, mappedAWSConfig)

		if err != nil {
			return err
		}

		c.AWS = *awsCfg
		c.cloud = "AWS"
	}

	return nil
}

func (c MdxcClient) createApplicationIdentity(ctx context.Context, d *schema.Resource) (*schema.Resource, error) {
	client.CloudSDK.createApplicationIdentity()
	return d, nil
}

func (c MdxcClient) deleteApplicationIdentity(ctx context.Context, d *schema.Resource) (*schema.Resource, error) {
	switch c.cloud {
	case "AWS":
		iamClient := app_identity.NewService(c.AWS)
		schemaResource, err := app_identity.Delete(ctx, iamClient, &massdriver.AppIdentityInput{Name: d.})
		// maybeTransform(?)
		// return schemaResource, err
	case "GCP":
		// iamClient := app_identity.NewService(c.GCP)
		// schemaResource, err app_identity.Create(ctx, iamClient, &massdriver.AppIdentityInput{Name: "foo"})
		// maybeTransform(?)
		// return schemaResource, err
	case "Azure":
		// iamClient := app_identity.NewService(c.Azure)
		// schemaResource, err app_identity.Create(ctx, iamClient, &massdriver.AppIdentityInput{Name: "foo"})
		// maybeTransform(?)
		// return schemaResource, err
	}

	return d, nil
}

func (c MdxcClient) attachPolicyToIdentity(ctx context.Context, d *schema.Resource) (*schema.Resource, error) {
	switch c.cloud {
	case "AWS":
		iamClient := app_identity.NewService(c.AWS)
		schemaResource, err := app_identity.Create(ctx, iamClient, &massdriver.AppIdentityInput{Name: d.})
		// maybeTransform(?)
		// return schemaResource, err
	case "GCP":
		// iamClient := app_identity.NewService(c.GCP)
		// schemaResource, err app_identity.Create(ctx, iamClient, &massdriver.AppIdentityInput{Name: "foo"})
		// maybeTransform(?)
		// return schemaResource, err
	case "Azure":
		// iamClient := app_identity.NewService(c.Azure)
		// schemaResource, err app_identity.Create(ctx, iamClient, &massdriver.AppIdentityInput{Name: "foo"})
		// maybeTransform(?)
		// return schemaResource, err
	}

	return d, nil
}

func initializeAWSConfig(ctx context.Context, awsMap map[string]interface{}) (*aws.Config, error) {
	awsConfig := awsConfig{}

	if roleArn, ok := awsMap["role_arn"].(string); ok && roleArn != "" {
		awsConfig.awsRoleArn = roleArn
	}
	if externalId, ok := awsMap["external_id"].(string); ok && externalId != "" {
		awsConfig.externalId = externalId
	}
	if region, ok := awsMap["region"].(string); ok && region != "" {
		awsConfig.region = region
	}

	cfg, err := config.LoadDefaultConfig(ctx, config.WithRegion(awsConfig.region))

	if err != nil {
		return nil, err
	}

	stsClient := sts.NewFromConfig(cfg)
	provider := stscreds.NewAssumeRoleProvider(stsClient, awsConfig.awsRoleArn, func(o *stscreds.AssumeRoleOptions) {
		o.ExternalID = aws.String(awsConfig.externalId)
	})
	cfg.Credentials = aws.NewCredentialsCache(provider)

	if err != nil {
		return nil, err
	}

	return &cfg, nil
}
