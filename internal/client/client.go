package client

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type MDXCClient interface {
	CreateAppIdentity(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics
	// ReadAppIdentity(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics
	// UpdateAppIdentity(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics
	DeleteAppIdentity(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics
}
