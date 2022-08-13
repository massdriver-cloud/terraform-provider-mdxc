package gcp

import (
	"context"
	"log"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

type GCPConfig struct {
	credentials string
	project     string
	tokenSource oauth2.TokenSource
}

func Initialize(ctx context.Context, gcpMap map[string]interface{}) (*GCPConfig, error) {
	gcpConfig := GCPConfig{}

	if credentials, ok := gcpMap["credentials"].(string); ok && credentials != "" {
		gcpConfig.credentials = credentials
	}
	if project, ok := gcpMap["project"].(string); ok && project != "" {
		gcpConfig.project = project
	}

	cfg, err := google.JWTConfigFromJSON([]byte(gcpConfig.credentials), "https://www.googleapis.com/auth/cloud-platform")
	if err != nil {
		return nil, err
	}

	gcpConfig.tokenSource = cfg.TokenSource(ctx)
	log.Printf("[debug] GCP Config Created")
	return &gcpConfig, nil
}
