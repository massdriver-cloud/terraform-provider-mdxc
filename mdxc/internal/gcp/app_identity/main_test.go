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

func TestCreateServiceAccount(t *testing.T) {
	appIdentityInput := massdriver.AppIdentityInput{
		Name: aws.String("test"),
	}

	ctx := context.Background()
	iamSvc := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
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

	defer iamSvc.Close()
	svc, err := iam.NewService(ctx, option.WithoutAuthentication(), option.WithEndpoint(iamSvc.URL))
	if err != nil {
		t.Fatalf("unable to create client: %v", err)
	}

	m := svc.Projects.ServiceAccounts

	serviceAccount, _ := app_identity.CreateServiceAccount(context.TODO(), m, &appIdentityInput)
	got := serviceAccount.Email
	want := "test@PROJECT_ID.iam.gserviceaccount.com"

	if want != got {
		t.Errorf("expect %v, got %v", want, got)
	}
}
