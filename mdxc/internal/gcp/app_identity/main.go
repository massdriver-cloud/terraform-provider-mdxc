// Package app_identity implements the massdriver.AppIdentity for GCP
package app_identity

// SA API Enablement
// resource "google_service_account" "application" {
//   account_id = module.k8s_application.params.md_metadata.name_prefix
// }

import (
	iam "google.golang.org/api/iam/v1"
)

// TODO: API Enablement
func Create() bool {
	foo := iam.CreateServiceAccountRequest{AccountId: "foo"}
	_ = foo
	return true
}
