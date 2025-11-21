package secrets

import (
	"context"
	"log/slog"
	"os"

	"github.com/joho/godotenv"

	secretmanager "cloud.google.com/go/secretmanager/apiv1"
	"cloud.google.com/go/secretmanager/apiv1/secretmanagerpb"
)

var secretManagerID string

func init() {
	_ = godotenv.Load()
	secretManagerID = os.Getenv("SECRET_MANAGER_ID")
	if secretManagerID == "" {
		panic("SECRET_MANAGER_ID not set in environment")
	}
}

func Secret(ctx context.Context, secret string) ([]byte, error) {
	client, err := secretmanager.NewClient(ctx)
	if err != nil {
		return nil, err
	}
	defer client.Close()

	secretName := "projects/" + secretManagerID + "/secrets/" + secret + "/versions/latest"

	req := &secretmanagerpb.AccessSecretVersionRequest{
		Name: secretName,
	}

	s, err := client.AccessSecretVersion(ctx, req)
	if err != nil {
		slog.Error("Failed to access secret version", "error", err, "secret-name", secretName)
		return nil, err
	}

	data := s.GetPayload().GetData()

	return data, nil
}
