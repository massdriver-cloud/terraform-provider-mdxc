package provider

import (
	"context"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure provider defined types fully satisfy framework interfaces
var _ provider.DataSourceType = DataSourceCloudType{}
var _ datasource.DataSource = DataSourceCloud{}

type DataSourceCloudType struct{}

func (t DataSourceCloudType) GetSchema(ctx context.Context) (tfsdk.Schema, diag.Diagnostics) {
	return tfsdk.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "A utility data source to get the provider's currently configured cloud",

		Attributes: map[string]tfsdk.Attribute{
			"cloud": {
				MarkdownDescription: "The currently configured cloud (aws, gcp, azure, etc)",
				Computed:            true,
				Type:                types.StringType,
			},
			"id": {
				MarkdownDescription: "The ID of the underlying cloud scope. For AWS, it will be the account ID, for GCP it will be the project ID, for Azure it will be the subcription ID.",
				Computed:            true,
				Type:                types.StringType,
			},
		},
	}, nil
}

func (t DataSourceCloudType) NewDataSource(ctx context.Context, in provider.Provider) (datasource.DataSource, diag.Diagnostics) {
	return DataSourceCloud{
		provider: *(in.(*MDXCProvider)),
	}, diag.Diagnostics{}
}

type DataSourceCloudData struct {
	Cloud types.String `tfsdk:"cloud"`
	ID    types.String `tfsdk:"id"`
}

type DataSourceCloud struct {
	provider MDXCProvider
}

func (d DataSourceCloud) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data DataSourceCloudData

	diags := req.Config.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	data.Cloud = types.String{Value: d.provider.Client.Cloud}

	switch data.Cloud.Value {
	case "aws":
		accountId := strings.Split(d.provider.Client.AWSConfig.Provider.AwsRoleArn.Value, ":")[4]
		data.ID = types.String{Value: accountId}
	case "gcp":
		data.ID = types.String{Value: d.provider.Client.GCPConfig.Provider.Project.Value}
	case "azure":
		data.ID = types.String{Value: d.provider.Client.AzureConfig.Provider.SubscriptionID.Value}
	default:
		resp.Diagnostics.Append(diag.NewErrorDiagnostic("Unrecognized cloud", "Failed to recognize cloud when extracting the ID: "+data.Cloud.Value))
		return
	}

	diags = resp.State.Set(ctx, &data)
	resp.Diagnostics.Append(diags...)
}
