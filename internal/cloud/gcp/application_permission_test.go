package gcp_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"terraform-provider-mdxc/internal/cloud/gcp"
	"testing"

	"google.golang.org/api/cloudresourcemanager/v1"
	"google.golang.org/api/option"
)

// func (c *gcp.GCPConfig) NewResourceManagerService(ctx context.Context) (gcp.GCPResourceManagerIface, error) {
// 	service, err := cloudresourcemanager.NewService(ctx, option.WithTokenSource(c.tokenSource))
// 	if err != nil {
// 		return nil, fmt.Errorf("cloudresourcemanager.NewService: %v", err)
// 	}

// 	return service.Projects, nil
// }

func createMockIamClient() (gcp.GCPResourceManagerIface, error) {
	// config := gcp.GCPConfig{}
	// config.NewResourceManagerService = func(ctx context.Context, tokenSource oauth2.TokenSource) (gcp.GCPResourceManagerIface, error) {
	ctx := context.Background()
	apiService := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		resp := &cloudresourcemanager.Policy{}
		b, err := json.Marshal(resp)
		if err != nil {
			http.Error(w, "unable to marshal request: "+err.Error(), http.StatusBadRequest)
			return
		}
		w.Write(b)
	}))

	service, err := cloudresourcemanager.NewService(ctx, option.WithoutAuthentication(), option.WithEndpoint(apiService.URL))
	if err != nil {
		return nil, err
	}
	return service.Projects, nil
}

func TestCreateServiceAccount(t *testing.T) {
	ctx := context.Background()
	config := &gcp.ApplicationPermissionConfig{
		ID: "test",
	}
	client, _ := createMockIamClient()
	response, _ := gcp.CreateApplicationPermission(ctx, config, client)

	got := response.Email
	want := "test@PROJECT_ID.iam.gserviceaccount.com"

	if want != got {
		t.Errorf("expect %v, got %v", want, got)
	}
}
