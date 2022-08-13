package azure

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func (c *AzureConfig) CreateApplicationIdentity(ctx context.Context, d *schema.ResourceData) diag.Diagnostics {
	return nil
}

func (c *AzureConfig) DeleteApplicationIdentity(ctx context.Context, d *schema.ResourceData) diag.Diagnostics {
	return nil
}
