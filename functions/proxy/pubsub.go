package proxy

import (
	"context"
	"fmt"
	"sync"

	"cloud.google.com/go/pubsub/v2"
)

var (
	pubsubClientInstance *pubsub.Client
	pubsubClientErr      error
	pubsubOnce           sync.Once
)

func PubSub() (*pubsub.Client, error) {
	pubsubOnce.Do(func() {
		ctx := context.Background()
		pubsubClientInstance, pubsubClientErr = pubsub.NewClient(ctx, projectID)
		if pubsubClientErr != nil {
			pubsubClientErr = fmt.Errorf("failed to create Pub/Sub client: %w", pubsubClientErr)
		}
	})
	return pubsubClientInstance, pubsubClientErr
}