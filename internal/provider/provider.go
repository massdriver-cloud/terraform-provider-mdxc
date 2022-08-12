package provider

import (
	"context"
	"errors"
	"log"
	"terraform-provider-mdxc/internal/cloud/aws"
	"terraform-provider-mdxc/internal/cloud/azure"
	"terraform-provider-mdxc/internal/cloud/gcp"
	"terraform-provider-mdxc/internal/mdxc"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var awsProviderSchema = schema.Schema{
	Type:          schema.TypeList,
	MaxItems:      1,
	Optional:      true,
	ConflictsWith: []string{"azure", "gcp"},
	Description:   "Credentials for AWS Cloud",
	Elem: &schema.Resource{
		Schema: map[string]*schema.Schema{
			"role_arn": {
				Optional:     true,
				Description:  "ARN of AWS Role to assume",
				Type:         schema.TypeString,
				Sensitive:    true,
				RequiredWith: []string{"aws.0.external_id", "aws.0.region"},
			},
			"external_id": {
				Optional:     true,
				Description:  "A unique identifier that might be required when you assume a role in another account.",
				Type:         schema.TypeString,
				Sensitive:    true,
				RequiredWith: []string{"aws.0.role_arn", "aws.0.region"},
			},
			"region": {
				Optional:     true,
				Description:  "The region where AWS operations will take place.",
				Type:         schema.TypeString,
				RequiredWith: []string{"aws.0.role_arn", "aws.0.external_id"},
			},
		},
	},
}

var azureProviderSchema = schema.Schema{
	Type:          schema.TypeList,
	MaxItems:      1,
	Optional:      true,
	ConflictsWith: []string{"aws", "gcp"},
	Description:   "Credentials for Azure Cloud. See how to authenticate through Service Principal in the [Azure docs](https://registry.terraform.io/providers/hashicorp/azurerm/latest/docs/guides/service_principal_client_secret#creating-a-service-principal)",
	Elem: &schema.Resource{
		Schema: map[string]*schema.Schema{
			"subscription_id": {
				Optional:     true,
				Description:  "Azure Subscription ID",
				Type:         schema.TypeString,
				Sensitive:    true,
				RequiredWith: []string{"azure.0.client_id", "azure.0.client_secret", "azure.0.tenant_id"},
			},
			"client_id": {
				Optional:     true,
				Description:  "Azure Client ID",
				Type:         schema.TypeString,
				RequiredWith: []string{"azure.0.subscription_id", "azure.0.client_secret", "azure.0.tenant_id"},
			},
			"client_secret": {
				Optional:     true,
				Description:  "Azure Client Secret",
				Type:         schema.TypeString,
				Sensitive:    true,
				RequiredWith: []string{"azure.0.subscription_id", "azure.0.client_id", "azure.0.tenant_id"},
			},
			"tenant_id": {
				Optional:     true,
				Description:  "Azure Tenant ID",
				Type:         schema.TypeString,
				Sensitive:    true,
				RequiredWith: []string{"azure.0.subscription_id", "azure.0.client_id", "azure.0.client_secret"},
			},
		},
	},
}
var gcpProviderSchema = schema.Schema{
	Type:          schema.TypeList,
	MaxItems:      1,
	Optional:      true,
	ConflictsWith: []string{"aws", "azure"},
	Description:   "Credentials for Google Cloud. See how to authenticate through Service Principals in the [Google docs](https://cloud.google.com/compute/docs/authentication)",
	Elem: &schema.Resource{
		Schema: map[string]*schema.Schema{
			"credentials": {
				Optional:     true,
				Description:  "Either the path to or the contents of a service account key file in JSON format.",
				Type:         schema.TypeString,
				Sensitive:    true,
				RequiredWith: []string{"gcp.0.project"},
			},
			"project": {
				Optional:     true,
				Description:  "The GCP project to manage resources in.",
				Type:         schema.TypeString,
				RequiredWith: []string{"gcp.0.credentials"},
			},
		},
	},
}

// Provider -
func Provider() *schema.Provider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"aws":   &awsProviderSchema,
			"azure": &azureProviderSchema,
			"gcp":   &gcpProviderSchema,
		},
		ResourcesMap: map[string]*schema.Resource{
			"mdxc_app_identity": mdxc.AppIdentityResource(),
		},
		DataSourcesMap:       map[string]*schema.Resource{},
		ConfigureContextFunc: providerConfigure,
	}
}

func providerConfigure(ctx context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {

	if awsBlock, ok := d.Get("aws").([]interface{}); ok && len(awsBlock) > 0 && awsBlock[0] != nil {
		log.Printf("[debug] Creating AWS client")
		mappedAWSConfig := awsBlock[0].(map[string]interface{})
		return aws.Initialize(ctx, d, mappedAWSConfig)
	}

	if azureBlock, ok := d.Get("azure").([]interface{}); ok && len(azureBlock) > 0 && azureBlock[0] != nil {
		log.Printf("[debug] Creating Azure client")
		mappedAzureConfig := azureBlock[0].(map[string]interface{})
		return azure.Initialize(ctx, d, mappedAzureConfig)
	}

	if gcpBlock, ok := d.Get("gcp").([]interface{}); ok && len(gcpBlock) > 0 && gcpBlock[0] != nil {
		log.Printf("[debug] Creating GCP client")
		mappedGCPConfig := gcpBlock[0].(map[string]interface{})
		return gcp.Initialize(ctx, d, mappedGCPConfig)
	}

	return nil, diag.FromErr(errors.New("At least one of 'aws', 'azure' or 'gcp' must be set"))
}
