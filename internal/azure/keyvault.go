package azure

import (
	"context"
	"fmt"
	"hrapplication/internal/utils"

	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/security/keyvault/azsecrets"
)

func GetCertAndKey(env map[string]string) (string, error) {
	cred, err := azidentity.NewEnvironmentCredential(nil)
	if err != nil {
		return "", fmt.Errorf("could not create credential: %v", err)
	}

	client, err := azsecrets.NewClient(env["AZURE_KEYVAULT_URL"], cred, nil)
	if err != nil {
		return "", fmt.Errorf("could not create secrets client: %v", err)
	}

	pemResp, err := client.GetSecret(context.Background(), env["AZURE_CERT_NAME"], "", nil)
	if err != nil {
		return "", fmt.Errorf("could not fetch PEM file: %v", err)
	}

	return utils.SaveToTempFile([]byte(*pemResp.Value), "pem-*.pem")
}

func GetSecret(env map[string]string, secretName string) (string, error) {
	cred, err := azidentity.NewEnvironmentCredential(nil)
	if err != nil {
		return "", fmt.Errorf("could not create credential: %v", err)
	}

	client, err := azsecrets.NewClient(env["AZURE_KEYVAULT_URL"], cred, nil)
	if err != nil {
		return "", fmt.Errorf("could not create secrets client: %v", err)
	}

	secretResp, err := client.GetSecret(context.Background(), secretName, "", nil)
	if err != nil {
		return "", fmt.Errorf("error retrieving secret '%s': %v", secretName, err)
	}

	return *secretResp.Value, nil
}
