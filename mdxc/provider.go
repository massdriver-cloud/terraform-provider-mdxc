package mdxc

import (
	"context"
	"log"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials/stscreds"
	"github.com/aws/aws-sdk-go-v2/service/sts"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var awsProviderSchema = schema.Schema{
	Type:        schema.TypeList,
	MaxItems:    1,
	Optional:    true,
	Computed:    true,
	Description: "Credentials for AWS Cloud",
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
	Type:        schema.TypeList,
	MaxItems:    1,
	Optional:    true,
	Description: "Credentials for Azure Cloud. See how to authenticate through Service Principal in the [Azure docs](https://registry.terraform.io/providers/hashicorp/azurerm/latest/docs/guides/service_principal_client_secret#creating-a-service-principal)",
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
	Type:        schema.TypeList,
	MaxItems:    1,
	Optional:    true,
	Description: "Credentials for Google Cloud. See how to authenticate through Service Principals in the [Google docs](https://cloud.google.com/compute/docs/authentication)",
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
			"mdxc_test": resourceTest(),
		},
		DataSourcesMap:       map[string]*schema.Resource{},
		ConfigureContextFunc: providerConfigure,
	}
}

type awsConfig struct {
	awsRoleArn string
	externalId string
	region     string
}
type azureConfig struct {
	subscriptionId string
	clientId       string
	clientSecret   string
	tenantId       string
}
type gcpConfig struct {
	credentials string
	project     string
}

func providerConfigure(ctx context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {
	var diags diag.Diagnostics

	var awsCfg *aws.Config
	azureConfig := azureConfig{}
	gcpConfig := gcpConfig{}

	if aws, ok := d.Get("aws").([]interface{}); ok && len(aws) > 0 && aws[0] != nil {
		log.Printf("[debug] Creating AWS client")
		mappedAWSConfig := aws[0].(map[string]interface{})
		awsCfg, diags = initializeAWSConfig(ctx, d, mappedAWSConfig)
		if awsCfg == nil {
			return nil, diags
		}
	}

	if azure, ok := d.Get("azure").([]interface{}); ok && len(azure) > 0 && azure[0] != nil {
		mappedAzureConfig := azure[0].(map[string]interface{})
		if subscriptionId, ok := mappedAzureConfig["subscription_id"].(string); ok && subscriptionId != "" {
			azureConfig.subscriptionId = subscriptionId
		}
		if clientId, ok := mappedAzureConfig["client_id"].(string); ok && clientId != "" {
			azureConfig.clientId = clientId
		}
		if clientSecret, ok := mappedAzureConfig["client_secret"].(string); ok && clientSecret != "" {
			azureConfig.clientSecret = clientSecret
		}
		if tenantId, ok := mappedAzureConfig["tenant_id"].(string); ok && tenantId != "" {
			azureConfig.tenantId = tenantId
		}
	}

	if gcp, ok := d.Get("gcp").([]interface{}); ok && len(gcp) > 0 && gcp[0] != nil {
		mappedGCPConfig := gcp[0].(map[string]interface{})
		if credentials, ok := mappedGCPConfig["credentials"].(string); ok && credentials != "" {
			gcpConfig.credentials = credentials
		}
		if project, ok := mappedGCPConfig["project"].(string); ok && project != "" {
			gcpConfig.project = project
		}
	}

	diags = append(diags, diag.Diagnostic{
		Severity: diag.Warning,
		Summary:  "azure.subscription_id: " + azureConfig.subscriptionId,
	},
		diag.Diagnostic{
			Severity: diag.Warning,
			Summary:  "azure.clientId: " + azureConfig.clientId,
		},
		diag.Diagnostic{
			Severity: diag.Warning,
			Summary:  "azure.client_secret: " + azureConfig.clientSecret,
		},
		diag.Diagnostic{
			Severity: diag.Warning,
			Summary:  "azure.tenant_id: " + azureConfig.tenantId,
		},
	)

	diags = append(diags, diag.Diagnostic{
		Severity: diag.Warning,
		Summary:  "gcp.credentials: " + gcpConfig.credentials,
	},
		diag.Diagnostic{
			Severity: diag.Warning,
			Summary:  "gcp.project: " + gcpConfig.project,
		},
	)

	log.Printf("[debug] Testing AWS Client")

	client := sts.NewFromConfig(*awsCfg)
	foo := sts.GetCallerIdentityInput{}
	log.Printf("[debug] Right before AWS Client")
	out, err := client.GetCallerIdentity(ctx, &foo)
	if err != nil {
		return nil, diag.FromErr(err)
	}
	log.Printf("[debug] Tested AWS Client")

	diags = append(diags, diag.Diagnostic{
		Severity: diag.Warning,
		Summary:  "aws sts get-caller-identity account: " + *out.Account,
	},
		diag.Diagnostic{
			Severity: diag.Warning,
			Summary:  "aws sts get-caller-identity arn: " + *out.Arn,
		},
		diag.Diagnostic{
			Severity: diag.Warning,
			Summary:  "aws sts get-caller-identity UserId: " + *out.UserId,
		},
	)

	return "foo", diags
}

func initializeAWSConfig(ctx context.Context, d *schema.ResourceData, awsMap map[string]interface{}) (*aws.Config, diag.Diagnostics) {
	var diags diag.Diagnostics
	awsConfig := awsConfig{}

	if roleArn, ok := awsMap["role_arn"].(string); ok && roleArn != "" {
		awsConfig.awsRoleArn = roleArn
	}
	if externalId, ok := awsMap["external_id"].(string); ok && externalId != "" {
		awsConfig.externalId = externalId
	}
	if region, ok := awsMap["region"].(string); ok && region != "" {
		awsConfig.region = region
	}

	log.Printf("[debug] Converting AWS values to config")
	cfg, err := config.LoadDefaultConfig(ctx, config.WithRegion(awsConfig.region))
	if err != nil {
		return nil, diag.FromErr(err)
	}
	stsClient := sts.NewFromConfig(cfg)
	provider := stscreds.NewAssumeRoleProvider(stsClient, awsConfig.awsRoleArn, func(o *stscreds.AssumeRoleOptions) {
		o.ExternalID = aws.String(awsConfig.externalId)
	})
	cfg.Credentials = aws.NewCredentialsCache(provider)
	log.Printf("[debug] AWS Config Created")
	return &cfg, diags
}

func resourceTest() *schema.Resource {
	return &schema.Resource{
		Description: "A Massdriver artifact for exporting a connectable type",

		CreateContext: resourceTestCreate,
		ReadContext:   schema.NoopContext,
		UpdateContext: resourceTestUpdate,
		DeleteContext: resourceTestDelete,

		Schema: map[string]*schema.Schema{
			"lol": {
				Description: "A json formatted string containing the artifact.",
				Type:        schema.TypeString,
				Required:    true,
			},
		},
	}
}

func resourceTestCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	lol := d.Get("lol").(string)
	var diags diag.Diagnostics
	d.SetId(time.Now().Format(time.RFC3339))
	d.Set("lol", lol)
	return diags
}

func resourceTestUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	lol := d.Get("lol").(string)
	var diags diag.Diagnostics
	d.Set("lol", lol)
	return diags
}

func resourceTestDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	d.SetId("")
	var diags diag.Diagnostics
	return diags
}
