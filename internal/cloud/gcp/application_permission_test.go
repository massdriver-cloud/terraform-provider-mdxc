package gcp_test

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"terraform-provider-mdxc/internal/cloud/gcp"
	"testing"

	"google.golang.org/api/cloudresourcemanager/v1"
	"google.golang.org/api/option"
)

func createMockPermissionClient() (gcp.GCPResourceManagerIface, error) {
	ctx := context.Background()
	apiService := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		resp := &cloudresourcemanager.Policy{
			Bindings: []*cloudresourcemanager.Binding{
				{
					Role: "roles/redis.viewer",
					Members: []string{
						"serviceAccount:test-name-prefix@test-project.iam.gserviceaccount.com",
					},
					Condition: &cloudresourcemanager.Expr{},
				},
			},
		}
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

func TestCreatePermission(t *testing.T) {
	ctx := context.Background()
	config := &gcp.ApplicationPermissionConfig{
		ServiceAccountID: "test-name-prefix@test-project.iam.gserviceaccount.com",
		Role:             "roles/redis.viewer",
		Condition:        "resource.name.startsWith(\"projects/test-project/locations/us-central1/instances/test-instance\")",
		Project:          "test-project",
	}
	client, _ := createMockPermissionClient()
	_ = gcp.CreateApplicationPermission(ctx, config, client)

	permissionID := fmt.Sprintf("%s-%s", config.ServiceAccountID, config.Role)

	compare(t, config.ID, permissionID)
}

func TestReadPermission(t *testing.T) {
	ctx := context.Background()
	config := &gcp.ApplicationPermissionConfig{
		ServiceAccountID: "test-name-prefix@test-project.iam.gserviceaccount.com",
		Role:             "roles/redis.viewer",
		Condition:        "resource.name.startsWith(\"projects/test-project/locations/us-central1/instances/test-instance\")",
		Project:          "test-project",
	}
	client, _ := createMockPermissionClient()
	_ = gcp.ReadApplicationPermission(ctx, config, client)

	// TODO: Add assertions
}
