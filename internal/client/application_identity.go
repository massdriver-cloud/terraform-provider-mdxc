package client

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func (c *MDXCClient) CreateApplicationIdentity(ctx context.Context, d *schema.ResourceData) diag.Diagnostics {
	return diag.Diagnostics{}
}

func (c *MDXCClient) DeleteApplicationIdentity(ctx context.Context, d *schema.ResourceData) diag.Diagnostics {
	return diag.Diagnostics{}
}
