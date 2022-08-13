package client

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func (c *MDXCClient) CreateApplicationIdentity(ctx context.Context, d *schema.ResourceData) diag.Diagnostics {
	cloudClient, _ := c.CloudConfig.NewService()
	return c.CloudConfig.CreateApplicationIdentity(ctx, d, cloudClient)
}

func (c *MDXCClient) DeleteApplicationIdentity(ctx context.Context, d *schema.ResourceData) diag.Diagnostics {
	cloudClient, _ := c.CloudConfig.NewService()
	return c.CloudConfig.DeleteApplicationIdentity(ctx, d, cloudClient)
}
