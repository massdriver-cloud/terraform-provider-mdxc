package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// type AppIdentityInput struct {
// 	Name *string
// }

// type AppIdentityOutput struct {
// 	AwsIamRole        awsTypes.Role
// 	GcpServiceAccount gcpTypes.ServiceAccount
// }

// Ensure provider defined types fully satisfy framework interfaces
var _ provider.ResourceType = ApplicationPermissionType{}
var _ resource.Resource = ApplicationPermission{}
var _ resource.ResourceWithImportState = ApplicationPermission{}

type ApplicationPermissionType struct{}

var awsApplicationPermissionSchema = tfsdk.Attribute{
	Optional:    true,
	Description: "AWS IAM Role Configuration",
	Attributes: tfsdk.SingleNestedAttributes(map[string]tfsdk.Attribute{

		"policy_arn": {
			Type:        types.StringType,
			Required:    true,
			Description: "AWS IAM Policy ARN to associate with the application identity (AWS IAM role)",
		},
	}),
}

func (t ApplicationPermissionType) GetSchema(ctx context.Context) (tfsdk.Schema, diag.Diagnostics) {
	return tfsdk.Schema{
		MarkdownDescription: "A cross-cloud application identity resource (AWS IAM Role, GCP Service Account, Azure Application)",

		Attributes: map[string]tfsdk.Attribute{
			"application_identity_id": {
				Type:        types.StringType,
				Description: "The ID of the Application Identity resource",
				Required:    true,
			},
			"aws": awsApplicationPermissionSchema,
		},
	}, nil
}

func (t ApplicationPermissionType) NewResource(ctx context.Context, in provider.Provider) (resource.Resource, diag.Diagnostics) {
	return ApplicationPermission{}, diag.Diagnostics{}
}

type ApplicationPermissionData struct {
	Name types.String `tfsdk:"name"`
	Id   types.String `tfsdk:"id"`
}

type ApplicationPermission struct {
	provider MDXCProvider
}

// func (r exampleResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
// 	return meta.(client.MDXCClient).CreateApplicationPermission(ctx, d)
// }

// func resourceApplicationPermissionRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
// 	return nil
// }

// // func resourceAppIdentityUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
// // 	return nil
// // }

// func resourceApplicationPermissionDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
// 	return meta.(client.MDXCClient).DeleteApplicationPermission(ctx, d)
// }

func (r ApplicationPermission) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data ApplicationPermissionData

	diags := req.Config.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	// If applicable, this is a great opportunity to initialize any necessary
	// provider client data and make a call using it.
	// example, err := d.provider.client.CreateExample(...)
	// if err != nil {
	//     resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create example, got error: %s", err))
	//     return
	// }

	// For the purposes of this example code, hardcoding a response value to
	// save into the Terraform state.
	data.Id = types.String{Value: "example-id"}

	// Write logs using the tflog package
	// Documentation: https://terraform.io/plugin/log
	tflog.Trace(ctx, "created a resource")

	diags = resp.State.Set(ctx, &data)
	resp.Diagnostics.Append(diags...)
}

func (r ApplicationPermission) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data ApplicationPermissionData

	diags := req.State.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	// If applicable, this is a great opportunity to initialize any necessary
	// provider client data and make a call using it.
	// example, err := d.provider.client.ReadExample(...)
	// if err != nil {
	//     resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read example, got error: %s", err))
	//     return
	// }

	diags = resp.State.Set(ctx, &data)
	resp.Diagnostics.Append(diags...)
}

func (r ApplicationPermission) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data ApplicationPermissionData

	diags := req.Plan.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	// If applicable, this is a great opportunity to initialize any necessary
	// provider client data and make a call using it.
	// example, err := d.provider.client.UpdateExample(...)
	// if err != nil {
	//     resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to update example, got error: %s", err))
	//     return
	// }

	diags = resp.State.Set(ctx, &data)
	resp.Diagnostics.Append(diags...)
}

func (r ApplicationPermission) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data ApplicationPermissionData

	diags := req.State.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	// If applicable, this is a great opportunity to initialize any necessary
	// provider client data and make a call using it.
	// example, err := d.provider.client.DeleteExample(...)
	// if err != nil {
	//     resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to delete example, got error: %s", err))
	//     return
	// }
}

func (r ApplicationPermission) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
