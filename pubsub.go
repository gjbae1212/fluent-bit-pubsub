package main

import (
	"context"
	"fmt"
	"os"

	"cloud.google.com/go/pubsub"
	"github.com/linkedin/goavro"
	"github.com/pkg/errors"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/option"
)

type CodecRecord = map[string]interface{}
type FBRecord = map[interface{}]interface{}

type Keeper interface {
	Send(ctx context.Context, data []byte) *pubsub.PublishResult
	Stop()
	InterfaceMapToByte(record FBRecord) ([]byte, error)
}

type GooglePubSub struct {
	client *pubsub.Client
	topic  *pubsub.Topic
	codec  func(record CodecRecord) ([]byte, error)
}

func NewKeeper(projectId, topicName, jwtPath string,
	publishSetting *pubsub.PublishSettings,
	schemaConfig *pubsub.SchemaConfig,
) (Keeper, error) {
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

	cfg, err := topic.Config(ctx)
	if err != nil {
		return nil, fmt.Errorf("topic.Config err: %v", err)
	}
	encoding := cfg.SchemaSettings.Encoding

	var codec func(record CodecRecord) ([]byte, error)
	var pubs *GooglePubSub
	switch schemaType {
	case pubsub.SchemaAvro:
		avroCodec, err := goavro.NewCodec(schemaConfig.Definition)
		if err != nil {
			return nil, fmt.Errorf("goavro.NewCodec err: %v", err)
		}

		switch encoding {
		case pubsub.EncodingBinary:
			codec = func(record CodecRecord) ([]byte, error) { return avroCodec.BinaryFromNative(nil, record) }
		case pubsub.EncodingJSON:
			codec = func(record CodecRecord) ([]byte, error) { return avroCodec.TextualFromNative(nil, record) }
		default:
			return nil, fmt.Errorf("invalid encoding: %v", encoding)
		}
		pubs = &GooglePubSub{client, topic, codec}
	// case pubsub.SchemaProtocolBuffer: [TODO]
	//    switch encoding {
	//    case pubsub.EncodingBinary:
	//        codec = func(record CodecRecord) ([]byte, error) {return ####}
	//    case pubsub.EncodingJSON:
	//        codec = func(record CodecRecord) ([]byte, error) {return ####}
	//    default;
	//        return nil, fmt.Errorf("invalid encoding: %v", encoding)
	//    pubs = &GooglePubSub{client, topic, codec}
	// }
	default:
		pubs = &GooglePubSub{client, topic, nil}
	}
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

func (gps *GooglePubSub) InterfaceMapToByte(fbr FBRecord) ([]byte, error) {
	cr := make(CodecRecord)
	for k, v := range fbr {
		strKey := fmt.Sprintf("%v", k)
		cv := convert(v)
		//d := reflect.ValueOf(cv).Kind()
		//fmt.Printf("[debug] {%v: %v} (%v)\n", k, cv, d)
		cr[strKey] = cv
	}
	return gps.codec(cr)
}

func convert(v interface{}) interface{} {
	switch d := v.(type) {
	case []byte:
		return string(d)
	default:
		return d
	}
}
