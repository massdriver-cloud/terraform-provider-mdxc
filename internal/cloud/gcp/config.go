package gcp

import (
	"context"
	"log"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"

	"github.com/hashicorp/terraform-plugin-framework/types"
)

type GCPProviderConfig struct {
	Credentials types.String `tfsdk:"credentials"`
	Project     types.String `tfsdk:"project"`
}

type GCPConfig struct {
	Provider                  *GCPProviderConfig
	TokenSource               oauth2.TokenSource
	NewIAMService             func(ctx context.Context, tokenSource oauth2.TokenSource) (GCPIamIface, error)
	NewResourceManagerService func(ctx context.Context, tokenSource oauth2.TokenSource) (GCPResourceManagerIface, error)
}

func Initialize(ctx context.Context, providerConfig *GCPProviderConfig) (*GCPConfig, error) {
	gcpConfig := GCPConfig{
		Provider:                  providerConfig,
		NewIAMService:             gcpIAMClientFactory,
		NewResourceManagerService: resourceManagerClientFactory,
	}

	log.Printf("[debug] Creating GCP TokenSource")
	cfg, err := google.JWTConfigFromJSON([]byte(providerConfig.Credentials.Value), "https://www.googleapis.com/auth/cloud-platform")
	if err != nil {
		return nil, err
	}

	gcpConfig.TokenSource = cfg.TokenSource(ctx)
	// gcpConfig.NewIAMService = gcpIAMClientFactory(ctx, gcpConfig.TokenSource)

	log.Printf("[debug] GCP Config Created")
	return &gcpConfig, nil
}
