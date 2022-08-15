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
	Provider    *GCPProviderConfig
	tokenSource oauth2.TokenSource
}

func Initialize(ctx context.Context, providerConfig *GCPProviderConfig) (*GCPConfig, error) {
	gcpConfig := GCPConfig{}

	gcpConfig.Provider = providerConfig

	cfg, err := google.JWTConfigFromJSON([]byte(providerConfig.Credentials.Value), "https://www.googleapis.com/auth/cloud-platform")
	if err != nil {
		return nil, err
	}

	gcpConfig.tokenSource = cfg.TokenSource(ctx)
	log.Printf("[debug] GCP Config Created")
	return &gcpConfig, nil
}
