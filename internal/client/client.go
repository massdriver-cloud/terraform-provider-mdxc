package client

import (
	"context"
	"errors"
	"terraform-provider-mdxc/internal/cloud/aws"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// type MDXCClient interface {
// 	CreateAppIdentity(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics
// 	// ReadAppIdentity(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics
// 	// UpdateAppIdentity(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics
// 	DeleteAppIdentity(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics
// }

type CloudConfig interface {
	NewService() (interface{}, error)

	CreateApplicationIdentity(ctx context.Context, d *schema.ResourceData, cloudClient interface{}) diag.Diagnostics
	// ReadAppIdentity(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics
	// UpdateAppIdentity(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics
	DeleteApplicationIdentity(ctx context.Context, d *schema.ResourceData, cloudClient interface{}) diag.Diagnostics
}

type MDXCClient struct {
	CloudConfig CloudConfig
}

func MDXCClientFactory(ctx context.Context, config map[string]interface{}, cloud string) (*MDXCClient, error) {
	mdxcClient := MDXCClient{}

	switch cloud {
	case "aws":
		var awsErr error
		mdxcClient.CloudConfig, awsErr = aws.Initialize(ctx, config)
		return &mdxcClient, awsErr
		// case "azure":
		// 	var azureErr error
		// 	mdxcClient.CloudConfig, azureErr = azure.Initialize(ctx, config)
		// 	return &mdxcClient, azureErr
		// case "gcp":
		// 	var gcpErr error
		// 	mdxcClient.CloudConfig, gcpErr = gcp.Initialize(ctx, config)
		// 	return &mdxcClient, gcpErr
	}
	return nil, errors.New("cloud not specified")
}
