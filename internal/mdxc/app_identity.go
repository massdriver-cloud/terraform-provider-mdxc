package mdxc

import (
	"context"
	"terraform-provider-mdxc/internal/client"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// type AppIdentityInput struct {
// 	Name *string
// }

// type AppIdentityOutput struct {
// 	AwsIamRole        awsTypes.Role
// 	GcpServiceAccount gcpTypes.ServiceAccount
// }

func AppIdentityResource() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceAppIdentityCreate,
		ReadContext:   resourceAppIdentityRead,
		//UpdateContext: schema.NoopContext,
		DeleteContext: resourceAppIdentityDelete,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Description: "The name of the IAM entity in the respective cloud (AWS IAM Role, GCP Service Account, Azure Application)",
				Required:    true,
				ForceNew:    true,
			},
		},
	}
}

func resourceAppIdentityCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return meta.(client.MDXCClient).CreateAppIdentity(ctx, d, meta)
}

func resourceAppIdentityRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return nil
}

// func resourceAppIdentityUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
// 	return nil
// }

func resourceAppIdentityDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return meta.(client.MDXCClient).DeleteAppIdentity(ctx, d, meta)
}
