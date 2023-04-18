package main

import (
	"context"
	"fmt"
	"os"

	"cloud.google.com/go/pubsub"
	"github.com/pkg/errors"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/option"
)

type Keeper interface {
	Send(ctx context.Context, data []byte) *pubsub.PublishResult
	Stop()
}

type GooglePubSub struct {
	client *pubsub.Client
	topic  *pubsub.Topic
}

func NewKeeper(projectId, topicName, jwtPath string,
	publishSetting *pubsub.PublishSettings) (Keeper, error) {
	if projectId == "" || topicName == "" || jwtPath == "" {
		return nil, fmt.Errorf("[err] NewKeeper empty params")
	}

	keyBytes, err := os.ReadFile(jwtPath)
	if err != nil {
		return nil, errors.Wrap(err, "[err] jwt path")
	}

	config, err := google.JWTConfigFromJSON(keyBytes, pubsub.ScopePubSub)
	if err != nil {
		return nil, errors.Wrap(err, "[err] jwt config")
	}

	ctx := context.Background()
	client, err := pubsub.NewClient(ctx, projectId, option.WithTokenSource(config.TokenSource(ctx)))
	if err != nil {
		return nil, errors.Wrap(err, "[err] pubsub client")
	}

	topic := client.Topic(topicName)
	if publishSetting != nil {
		topic.PublishSettings = *publishSetting
	} else {
		topic.PublishSettings = pubsub.DefaultPublishSettings
	}

	pubs := &GooglePubSub{client: client, topic: topic}
	return Keeper(pubs), nil
}

func (gps *GooglePubSub) Send(ctx context.Context, data []byte) *pubsub.PublishResult {
	if len(data) == 0 {
		return nil
	}
	return gps.topic.Publish(ctx, &pubsub.Message{Data: data})
}

func (gps *GooglePubSub) Stop() {
	gps.topic.Stop()
}
