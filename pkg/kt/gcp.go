package kt

/**
Helper functions to make working with GCP easier.
*/

import (
	"context"
	"fmt"
	"log"

	secretmanager "cloud.google.com/go/secretmanager/apiv1"
	secretmanagerpb "google.golang.org/genproto/googleapis/cloud/secretmanager/v1"
)

var (
	gSMClient  *secretmanager.Client
	gcpProject = ""
)

const (
	GCPSmPrefix  = "gcp.sm."
	GCPErrPrefix = "GCPError: "
)

func loadViaGCPSecrets(key string) string {
	ctx := context.Background()

	if gSMClient == nil {
		client, err := secretmanager.NewClient(ctx)
		if err != nil {
			log.Printf("failed to create secretmanager client: %v", err)
			return GCPErrPrefix + fmt.Sprintf("failed to create secretmanager client: %v", err)
		}
		gSMClient = client
	}

	// We need this project explicitly set for us also.
	if gcpProject == "" {
		gcpProject = LookupEnvString("GOOGLE_CLOUD_PROJECT", "")
		if gcpProject == "" {
			log.Printf("Missing env var GOOGLE_CLOUD_PROJECT")
			return fmt.Sprintf("Missing env var GOOGLE_CLOUD_PROJECT")
		}
	}

	// Build the request.
	key = fmt.Sprintf("projects/%s/secrets/%s/versions/latest", gcpProject, key)
	req := &secretmanagerpb.AccessSecretVersionRequest{
		Name: key,
	}

	// Call the API.
	result, err := gSMClient.AccessSecretVersion(ctx, req)
	if err != nil {
		log.Printf("failed to get gcp secret: %v", err)
		return fmt.Sprintf("failed to get gcp secret: %v", err)
	}

	return string(result.Payload.Data)
}
