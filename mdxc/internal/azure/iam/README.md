# Azure App Identity POC
This shows off example functions for creating and destroying Applicaiton Identities using the azure go sdk.

## Running the script
Credentials for authing requests to azure inferred from following environment variables:
```bash
export AZURE_TENANT_ID="<active_directory_tenant_id"
export AZURE_CLIENT_ID="<service_principal_appid>"
export AZURE_CLIENT_SECRET="<service_principal_password>"
export AZURE_SUBSCRIPTION_ID="<subscription_id>"
```
you can grab all of this info from an azure service principal artifact in massdriver
and toss them in a `.env` file
```bash
go run main.go
```

this script is basically a go implementaiton of this terraform https://github.com/massdriver-cloud/massdriver-bundles/blob/main/provisioners/terraform/modules/k8s-application-azure/iam.tf
