package gcp_test

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"terraform-provider-mdxc/internal/cloud/gcp"
	"testing"

	"google.golang.org/api/iam/v1"
	"google.golang.org/api/option"
)

func createMockIamClient() (gcp.GCPIamIface, error) {
	ctx := context.Background()
	apiService := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		project := strings.Split(r.URL.String(), "/")[3]

		resp := &iam.ServiceAccount{
			Email:       fmt.Sprintf("test-name-prefix@%s.iam.gserviceaccount.com", project),
			DisplayName: "test-name-prefix",
			ProjectId:   project,
		}

		b, err := json.Marshal(resp)
		if err != nil {
			http.Error(w, "unable to marshal request: "+err.Error(), http.StatusBadRequest)
			return
		}
		w.Write(b)
	}))

	service, err := iam.NewService(ctx, option.WithoutAuthentication(), option.WithEndpoint(apiService.URL))
	if err != nil {
		return nil, err
	}

	return service.Projects.ServiceAccounts, nil
}

func TestCreateIdentity(t *testing.T) {
	ctx := context.Background()
	config := &gcp.ApplicationIdentityConfig{
		Name:    "test-name-prefix",
		Project: "test-project",
	}
	client, _ := createMockIamClient()
	_ = gcp.CreateApplicationIdentity(ctx, config, client)

	compare(t, config.ID, "test-name-prefix@test-project.iam.gserviceaccount.com")
	compare(t, config.Name, "test-name-prefix")
	compare(t, config.Project, "test-project")
}

func TestReadIdentity(t *testing.T) {
	ctx := context.Background()
	config := &gcp.ApplicationIdentityConfig{
		Name:    "test-name-prefix",
		Project: "test-project",
	}
	client, _ := createMockIamClient()
	_ = gcp.ReadApplicationIdentity(ctx, config, client)

	compare(t, config.ID, "test-name-prefix@test-project.iam.gserviceaccount.com")
	compare(t, config.Name, "test-name-prefix")
	compare(t, config.Project, "test-project")
}

func compare(t *testing.T, got string, want string) {
	if want != got {
		t.Errorf("expect %v, got %v", want, got)
	}
}
