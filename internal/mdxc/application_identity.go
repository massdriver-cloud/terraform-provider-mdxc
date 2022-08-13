package mdxc

import (
	"context"
	"terraform-provider-mdxc/internal/client"
	"terraform-provider-mdxc/internal/verify"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/structure"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

// type AppIdentityInput struct {
// 	Name *string
// }

// type AppIdentityOutput struct {
// 	AwsIamRole        awsTypes.Role
// 	GcpServiceAccount gcpTypes.ServiceAccount
// }

var awsIAMRoleSchema = schema.Schema{
	Type:        schema.TypeList,
	MaxItems:    1,
	Optional:    true,
	ForceNew:    true, // REMOVE ME!!!!!!!!!!!!!!!!!!
	Description: "AWS IAM Role Configuration",
	Elem: &schema.Resource{
		Schema: map[string]*schema.Schema{
			"assume_role_policy": {
				Type:             schema.TypeString,
				ForceNew:         true, // REMOVE ME!!!!!!!!!!!!!!!!!!!
				Required:         true,
				ValidateFunc:     validation.StringIsJSON,
				DiffSuppressFunc: verify.SuppressEquivalentPolicyDiffs,
				StateFunc: func(v interface{}) string {
					json, _ := structure.NormalizeJsonString(v)
					return json
				},
			},
		},
	},
}

func ApplicationIdentityResource() *schema.Resource {
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
			"aws": &awsIAMRoleSchema,
		},
	}
}

func resourceAppIdentityCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return meta.(*client.MDXCClient).CreateApplicationIdentity(ctx, d)
}

func resourceAppIdentityRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return nil
}

// func resourceAppIdentityUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
// 	return nil
// }

func resourceAppIdentityDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return meta.(*client.MDXCClient).DeleteApplicationIdentity(ctx, d)
}
