package mdxc

import (
	"context"
	"errors"
	"terraform-provider-mdxc/internal/cloud/aws"
	"terraform-provider-mdxc/internal/cloud/azure"
	"terraform-provider-mdxc/internal/cloud/gcp"
)

// type MDXCClient interface {
// 	CreateApplicationIdentity(ctx context.Context, d *schema.ResourceData) diag.Diagnostics
// 	// ReadAppIdentity(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics
// 	// UpdateAppIdentity(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics
// 	DeleteApplicationIdentity(ctx context.Context, d *schema.ResourceData) diag.Diagnostics

// 	// CreateApplicationPermission(ctx context.Context, d *schema.ResourceData) diag.Diagnostics
// 	// // ReadAppIdentity(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics
// 	// // UpdateAppIdentity(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics
// 	// DeleteApplicationPermission(ctx context.Context, d *schema.ResourceData) diag.Diagnostics
// }

// type CloudConfig interface {
// 	CreateApplicationIdentity(ctx context.Context, d *schema.ResourceData) diag.Diagnostics
// 	// ReadAppIdentity(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics
// 	// UpdateAppIdentity(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics
// 	DeleteApplicationIdentity(ctx context.Context, d *schema.ResourceData) diag.Diagnostics
// }

// type MDXCClient struct {
// 	CloudConfig CloudConfig
// }

type MDXCProviderConfig struct {
	AWS   *aws.AWSProviderConfig     `tfsdk:"aws"`
	Azure *azure.AzureProviderConfig `tfsdk:"azure"`
	GCP   *gcp.GCPProviderConfig     `tfsdk:"gcp"`
}

type MDXCClient struct {
	Cloud       string
	AWSConfig   *aws.AWSConfig
	AzureConfig *azure.AzureConfig
	GCPConfig   *gcp.GCPConfig
}

func MDXCClientFactory(ctx context.Context, config *MDXCProviderConfig) (*MDXCClient, error) {
	client := MDXCClient{}
	var err error
	if config.AWS != nil {
		client.Cloud = "aws"
		client.AWSConfig, err = aws.Initialize(ctx, config.AWS)
		if err != nil {
			return nil, err
		}
		return &client, nil
	}
	if config.Azure != nil {
		client.Cloud = "azure"
		client.AzureConfig, err = azure.Initialize(ctx, config.Azure)
		if err != nil {
			return nil, err
		}
		return &client, nil
	}
	if config.GCP != nil {
		client.Cloud = "gcp"
		client.GCPConfig, err = gcp.Initialize(ctx, config.GCP)
		if err != nil {
			return nil, err
		}
		return &client, nil
	}

	return nil, errors.New("at least one of 'aws', 'azure' or 'gcp' must be set")
}
