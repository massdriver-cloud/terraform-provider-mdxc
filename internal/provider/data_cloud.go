package provider

import (
	"context"

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

	diags = resp.State.Set(ctx, &data)
	resp.Diagnostics.Append(diags...)
}
