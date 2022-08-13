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

var awsPermissionSchema = schema.Schema{
	Type:        schema.TypeList,
	MaxItems:    1,
	Optional:    true,
	ForceNew:    true, // REMOVE ME!!!!!!!!!!!!!!!!!!
	Description: "AWS IAM Role Configuration",
	Elem: &schema.Resource{
		Schema: map[string]*schema.Schema{
			"policy_arn": {
				Type:        schema.TypeString,
				ForceNew:    true, // REMOVE ME!!!!!!!!!!!!!!!!!!!
				Required:    true,
				Description: "AWS IAM Policy ARN to associate with the application identity (AWS IAM role)",
			},
		},
	},
}

func ApplicationPermissionResource() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceApplicationPermissionCreate,
		ReadContext:   resourceApplicationPermissionRead,
		//UpdateContext: schema.NoopContext,
		DeleteContext: resourceApplicationPermissionDelete,

		Schema: map[string]*schema.Schema{
			"application_identity_id": {
				Type:        schema.TypeString,
				Description: "The ID of the Application Identity resource",
				Required:    true,
				ForceNew:    true,
			},
			"aws": &awsPermissionSchema,
		},
	}
}

func resourceApplicationPermissionCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return meta.(client.MDXCClient).CreateApplicationPermission(ctx, d)
}

func resourceApplicationPermissionRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return nil
}

// func resourceAppIdentityUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
// 	return nil
// }

func resourceApplicationPermissionDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return meta.(client.MDXCClient).DeleteApplicationPermission(ctx, d)
}
