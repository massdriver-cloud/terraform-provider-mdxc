package azure

import (
	"context"
	"fmt"
	"strings"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/msi/armmsi"
)

type ManagedIdentityClient interface {
	CreateOrUpdate(ctx context.Context, resourceGroupName string, resourceName string, parameters armmsi.Identity, options *armmsi.UserAssignedIdentitiesClientCreateOrUpdateOptions) (armmsi.UserAssignedIdentitiesClientCreateOrUpdateResponse, error)
	Get(ctx context.Context, resourceGroupName string, resourceName string, options *armmsi.UserAssignedIdentitiesClientGetOptions) (armmsi.UserAssignedIdentitiesClientGetResponse, error)
	Delete(ctx context.Context, resourceGroupName string, resourceName string, options *armmsi.UserAssignedIdentitiesClientDeleteOptions) (armmsi.UserAssignedIdentitiesClientDeleteResponse, error)
}

type FederatedIdentityCredentialClient interface {
	CreateOrUpdate(ctx context.Context, resourceGroupName string, resourceName string, federatedIdentityCredentialResourceName string, parameters armmsi.FederatedIdentityCredential, options *armmsi.FederatedIdentityCredentialsClientCreateOrUpdateOptions) (armmsi.FederatedIdentityCredentialsClientCreateOrUpdateResponse, error)
}

type ApplicationIdentityConfig struct {
	// READ-ONLY; Fully qualified resource ID for the resource.
	// /subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/{resourceProviderNamespace}/{resourceType}/{resourceName}
	ID                           string
	Name                         string
	Location                     string
	ResourceGroupName            string
	KubernetesNamspace           string
	KubernetesServiceAccountName string
	KubernetesOIDCURL            string
	ClientID                     string
	TenantID                     string
	ResourceID                   string
}

func newManagedIdentityClientFactory(ctx context.Context, config *AzureProviderConfig) (ManagedIdentityClient, error) {
	cred, err := azidentity.NewClientSecretCredential(config.TenantID.Value, config.ClientID.Value, config.ClientSecret.Value, nil)
	if err != nil {
		return nil, err
	}

	client, errClient := armmsi.NewUserAssignedIdentitiesClient(config.SubscriptionID.Value, cred, nil)
	if errClient != nil {
		return nil, errClient
	}

	return client, nil
}

func newFederatedIdentityCredentialClientFactory(ctx context.Context, config *AzureProviderConfig) (FederatedIdentityCredentialClient, error) {
	cred, err := azidentity.NewClientSecretCredential(config.TenantID.Value, config.ClientID.Value, config.ClientSecret.Value, nil)
	if err != nil {
		return nil, err
	}

	client, errClient := armmsi.NewFederatedIdentityCredentialsClient(config.SubscriptionID.Value, cred, nil)
	if errClient != nil {
		return nil, errClient
	}

	return client, nil
}

func CreateApplicationIdentity(ctx context.Context, config *ApplicationIdentityConfig, client ManagedIdentityClient, fedClient FederatedIdentityCredentialClient) error {
	identity, errCreate := client.CreateOrUpdate(ctx,
		config.ResourceGroupName,
		config.Name,
		armmsi.Identity{
			Location: &config.Location,
			Tags: map[string]*string{
				"managed-by": to.Ptr("massdriver"),
			},
		},
		nil,
	)
	if errCreate != nil {
		return errCreate
	}

	resourceID := *identity.ID
	id := strings.Replace(resourceID, "resourcegroup", "resourceGroup", -1)
	config.ID = *identity.Properties.PrincipalID
	config.ClientID = *identity.Properties.ClientID
	config.TenantID = *identity.Properties.TenantID
	config.ResourceID = id

	if config.KubernetesNamspace != "" {
		if errAddRole := addWorkloadIdentityRole(ctx, config, fedClient); errAddRole != nil {
			return errAddRole
		}
	}
	return nil
}

func ReadApplicationIdentity(ctx context.Context, config *ApplicationIdentityConfig, client ManagedIdentityClient, fedClient FederatedIdentityCredentialClient) error {
	identity, err := client.Get(ctx, config.ResourceGroupName, config.Name, nil)
	if err != nil {
		return err
	}

	id := strings.Replace(*identity.ID, "resourcegroup", "resourceGroup", -1)
	config.ID = *identity.Properties.PrincipalID
	config.ClientID = *identity.Properties.ClientID
	config.TenantID = *identity.Properties.TenantID
	config.ResourceID = id

	return nil
}

func UpdateApplicationIdentity(ctx context.Context, config *ApplicationIdentityConfig, client ManagedIdentityClient, fedClient FederatedIdentityCredentialClient) error {
	return nil
}

func DeleteApplicationIdentity(ctx context.Context, config *ApplicationIdentityConfig, client ManagedIdentityClient, fedClient FederatedIdentityCredentialClient) error {
	_, err := client.Delete(ctx, config.ResourceGroupName, config.Name, nil)
	if err != nil {
		return err
	}

	config.ID = ""

	return nil
}

func addWorkloadIdentityRole(ctx context.Context, config *ApplicationIdentityConfig, client FederatedIdentityCredentialClient) error {
	_, err := client.CreateOrUpdate(ctx,
		config.ResourceGroupName,
		config.Name,
		// This is the Name of the Federated Identity Credential,
		// we can use the same name as the resource group, they are guranteed
		/// to be unique to this application
		config.Name,
		armmsi.FederatedIdentityCredential{
			Properties: &armmsi.FederatedIdentityCredentialProperties{
				Audiences: []*string{
					to.Ptr("api://AzureADTokenExchange"),
				},
				Issuer: &config.KubernetesOIDCURL,
				// k8s service account
				Subject: to.Ptr(fmt.Sprintf("system:serviceaccount:%s:%s", config.KubernetesNamspace, config.Name)),
			},
		},
		nil)
	if err != nil {
		return err
	}

	return nil
}
