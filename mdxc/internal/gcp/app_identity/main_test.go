package app_identity_test

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"terraform-provider-mdxc/mdxc/internal/gcp/app_identity"
	"terraform-provider-mdxc/mdxc/internal/massdriver"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"google.golang.org/api/iam/v1"
	"google.golang.org/api/option"
)

// Testing using HTTP Service mock: https://github.com/googleapis/google-api-go-client/blob/main/testing.md

func TestCreateServiceAccount(t *testing.T) {
	appIdentityInput := massdriver.AppIdentityInput{
		Name: aws.String("test"),
	}

	ctx := context.Background()
	mock := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		createReq := iam.CreateServiceAccountRequest{}

		body, _ := io.ReadAll(r.Body)
		json.Unmarshal([]byte(body), &createReq)

		email := fmt.Sprintf("%s@PROJECT_ID.iam.gserviceaccount.com", createReq.AccountId)
		resp := &iam.ServiceAccount{Email: email}

		b, err := json.Marshal(resp)
		if err != nil {
			http.Error(w, "unable to marshal request: "+err.Error(), http.StatusBadRequest)
			return
		}
		w.Write(b)
	}))

	defer mock.Close()
	iamSvc, err := iam.NewService(ctx, option.WithoutAuthentication(), option.WithEndpoint(mock.URL))
	if err != nil {
		t.Fatalf("unable to create client: %v", err)
	}

	mockSvc := iamSvc.Projects.ServiceAccounts

	serviceAccount, _ := app_identity.CreateServiceAccount(context.TODO(), mockSvc, &appIdentityInput)
	got := serviceAccount.Email
	want := "test@PROJECT_ID.iam.gserviceaccount.com"

	if want != got {
		t.Errorf("expect %v, got %v", want, got)
	}
}

func TestBindServiceAccountUserRole(t *testing.T) {
	api := "WHICH_API_IS_USED_TO_BIND_A_ROLE_TO_MEMBERS???"
	svcAcct := iam.ServiceAccount{Email: "test@foo123.iam.gserviceaccount.com"}
	binding := app_identity.BindServiceAccountUserRole(context.TODO(), api, &svcAcct)
	got := binding.Role
	want := "roles/iam.serviceAccountUser"

	if want != got {
		t.Errorf("expect %v, got %v", want, got)
	}
}
