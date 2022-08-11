package app_identity_test

import (
	"terraform-provider-mdxc/mdxc/internal/gcp/app_identity"
	"testing"
)

func TestCreate(t *testing.T) {
	got := app_identity.Create()
	want := true

	if want != got {
		t.Errorf("expect %v, got %v", want, got)
	}
}
