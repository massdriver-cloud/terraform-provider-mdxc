package client

import (
	"context"
	"errors"
	"terraform-provider-mdxc/internal/cloud/aws"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type MDXCClient interface {
	CreateApplicationIdentity(ctx context.Context, d *schema.ResourceData) diag.Diagnostics
	// ReadAppIdentity(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics
	// UpdateAppIdentity(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics
	DeleteApplicationIdentity(ctx context.Context, d *schema.ResourceData) diag.Diagnostics

	CreateApplicationPermission(ctx context.Context, d *schema.ResourceData) diag.Diagnostics
	// ReadAppIdentity(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics
	// UpdateAppIdentity(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics
	DeleteApplicationPermission(ctx context.Context, d *schema.ResourceData) diag.Diagnostics
}

// type CloudConfig interface {
// 	CreateApplicationIdentity(ctx context.Context, d *schema.ResourceData) diag.Diagnostics
// 	// ReadAppIdentity(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics
// 	// UpdateAppIdentity(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics
// 	DeleteApplicationIdentity(ctx context.Context, d *schema.ResourceData) diag.Diagnostics
// }

// type MDXCClient struct {
// 	CloudConfig CloudConfig
// }

func MDXCClientFactory(ctx context.Context, config map[string]interface{}, cloud string) (MDXCClient, error) {
	switch cloud {
	case "aws":
		return aws.Initialize(ctx, config)
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
