package gcp

import (
	"context"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

type GcpClient struct {
	credentials string
	project     string
	tokenSource oauth2.TokenSource
}

func Initialize(ctx context.Context, d *schema.ResourceData, gcpMap map[string]interface{}) (*GcpClient, diag.Diagnostics) {
	var diags diag.Diagnostics
	gcpConfig := GcpClient{}

	if credentials, ok := gcpMap["credentials"].(string); ok && credentials != "" {
		gcpConfig.credentials = credentials
	}
	if project, ok := gcpMap["project"].(string); ok && project != "" {
		gcpConfig.project = project
	}

	cfg, err := google.JWTConfigFromJSON([]byte(gcpConfig.credentials), "https://www.googleapis.com/auth/cloud-platform")
	if err != nil {
		return nil, diag.FromErr(err)
	}

	gcpConfig.tokenSource = cfg.TokenSource(ctx)
	log.Printf("[debug] GCP Config Created")
	return &gcpConfig, diags
}
