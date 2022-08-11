package app_identity_test

import (
	"context"
	"terraform-provider-mdxc/mdxc/internal/gcp/app_identity"
	"terraform-provider-mdxc/mdxc/internal/massdriver"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
)

// func mockNewIamServiceAccountsService() {
// 	ctx := context.TODO()
// 	return iam.NewService(ctx)
// }

func TestCreate(t *testing.T) {

	appIdentityInput := massdriver.AppIdentityInput{
		Name: aws.String("test"),
	}

	// m := mockNewIamServiceAccountsService()
	m := "wip"
	appIdentityOutput, _ := app_identity.Create(context.TODO(), m, &appIdentityInput)
	got := appIdentityOutput.GcpServiceAccount.Email
	want := "test@PROJECT_ID.iam.gserviceaccount.com"

	if want != got {
		t.Errorf("expect %v, got %v", want, got)
	}
}
