package snap

import (
	"context"
	"os"

	"cloud.google.com/go/firestore"
)

var firestoreDB string

func init() {
	firestoreDB = os.Getenv("FIRESTORE_DB")
	if firestoreDB == "" {
		panic("FIRESTORE_DB not set in environment")
	}
}

func Firestore(ctx context.Context) (*firestore.Client, error) {
	client, err := firestore.NewClientWithDatabase(ctx, projectID, firestoreDB)
	if err != nil {
		return nil, err
	}
	return client, nil
}
