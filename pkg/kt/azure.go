package kt

/**
Helper functions to make working with Azure easier.
*/

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/keyvault/azsecrets"
)

var (
	azureClient *azsecrets.Client
)

const (
	AzureKVPrefix  = "azure.kv."
	AzureErrPrefix = "AzureError: "
	azureTimeout   = 10 * time.Second
)

func loadViaAzureKeyVault(key string) string {
	if azureClient == nil {
		keyVaultURL := os.Getenv("KT_AZURE_KEY_VAULT_URL") // Optional, default to name + common url.
		if keyVaultURL == "" {
			keyVaultName := os.Getenv("KT_AZURE_KEY_VAULT_NAME") // Must be set.
			if keyVaultName == "" {
				log.Printf("ENV Var 'KT_AZURE_KEY_VAULT_NAME' or 'KT_AZURE_KEY_VAULT_URL' must be set")
				return "ENV Var KT_AZURE_KEY_VAULT_NAME or KT_AZURE_KEY_VAULT_URL must be set"
			}
			keyVaultURL = fmt.Sprintf("https://%s.vault.azure.net/", keyVaultName)
		}

		cred, err := azidentity.NewDefaultAzureCredential(nil)
		if err != nil {
			log.Printf("failed to obtain credential: %v", err)
			return AzureErrPrefix + fmt.Sprintf("failed to obtain credential: %v", err)
		}

		client, err := azsecrets.NewClient(keyVaultURL, cred, nil)
		if err != nil {
			log.Printf("failed to connect to client: %v", err)
			return AzureErrPrefix + fmt.Sprintf("failed to connect to client: %v", err)
		}
		azureClient = client
	}

	ctx, cancel := context.WithTimeout(context.Background(), azureTimeout)
	defer cancel()
	resp, err := azureClient.GetSecret(ctx, key, nil)
	if err != nil {
		log.Printf("failed to get azure secret: %v", err)
		return fmt.Sprintf("failed to get azure secret: %v", err)
	}
	return *resp.Value
}
