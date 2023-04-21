package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"testing"
	"time"
	"unsafe"

	"cloud.google.com/go/pubsub"

	"github.com/fluent/fluent-bit-go/output"
	"github.com/stretchr/testify/assert"
)

type testOutput struct {
	inc int
}

func (o *testOutput) Register(ctx unsafe.Pointer, name string, desc string) int {
	return output.FLBPluginRegister(ctx, name, desc)
}

func (o *testOutput) GetConfigKey(ctx unsafe.Pointer, key string) string {
	if key == "Project" {
		return os.Getenv("PROJECT_ID")
	}
	if key == "Topic" {
		return os.Getenv("TOPIC_NAME")
	}
	if key == "JwtPath" {
		return os.Getenv("JWT_PATH")
	}
	if key == "Debug" {
		return "true"
	}
	if key == "Timeout" {
		return "10000"
	}
	if key == "ByteThreshold" {
		return "1000000"
	}
	if key == "CountThreshold" {
		return "100"
	}
	if key == "DelayThreshold" {
		return "100"
	}
	return ""
}

func (o *testOutput) NewDecoder(data unsafe.Pointer, length int) *output.FLBDecoder {
	return nil
}

func (o *testOutput) GetRecord(dec *output.FLBDecoder) (ret int, ts interface{}, rec map[interface{}]interface{}) {
	if o.inc == 0 {
		o.inc++
		return 0, output.FLBTime{Time: time.Now()}, map[interface{}]interface{}{
			"testvalue1": []byte("record1"),
			"testvalue2": []byte("record2"),
		}
	}
	return -1, nil, nil
}

func TestFLBPluginInit(t *testing.T) {
	assert := assert.New(t)
	wrapper = OutputWrapper(&testOutput{})
	if os.Getenv("PROJECT_ID") == "" || os.Getenv("TOPIC_NAME") == "" ||
		os.Getenv("JWT_PATH") == "" {
		assert.Equal(output.FLB_ERROR, FLBPluginInit(nil))
	} else {
		assert.Equal(output.FLB_OK, FLBPluginInit(nil))
	}
}

func TestFLBPluginFlush(t *testing.T) {
	assert := assert.New(t)
	wrapper = OutputWrapper(&testOutput{})
	if os.Getenv("PROJECT_ID") == "" || os.Getenv("TOPIC_NAME") == "" ||
		os.Getenv("JWT_PATH") == "" {
		return
	}
	ok := FLBPluginFlush(nil, 0, nil)
	assert.Equal(output.FLB_OK, ok)

	projectId := os.Getenv("PROJECT_ID")
	topicName := os.Getenv("TOPIC_NAME")
	jwtPath := os.Getenv("JWT_PATH")
	if projectId == "" || topicName == "" || jwtPath == "" {
		return
	}
	keeper, err := NewKeeper(projectId, topicName, jwtPath, nil, nil)
	assert.NoError(err)
	sub := keeper.(*GooglePubSub).client.Subscription(topicName)
	go func() {
		err = sub.Receive(context.Background(), func(ctx context.Context, m *pubsub.Message) {
			log.Printf("Got message: %s", m.Data)
			m.Ack()
		})
	}()
	time.Sleep(5 * time.Second)
}

func TestInterfaceToBytes(t *testing.T) {
	assert := assert.New(t)

	now := time.Now()
	tests := map[string]struct {
		input  interface{}
		output []byte
	}{
		"float": {
			input:  float64(10.0),
			output: []byte(fmt.Sprintf("%f", float64(10.0))),
		},
		"[]byte": {
			input:  []byte(string("hello")),
			output: []byte(string("hello")),
		},
		"int": {
			input:  int(20),
			output: []byte(string("20")),
		},
		"string": {
			input:  "hello",
			output: []byte(string("hello")),
		},
		"time": {
			input:  now,
			output: []byte(now.Format(time.RFC3339)),
		},
		"bool": {
			input:  true,
			output: []byte("true"),
		},
		"etc": {
			input:  map[string]string{"hello": "world"},
			output: []byte(fmt.Sprintf("%v", map[string]string{"hello": "world"})),
		},
	}

	for _, t := range tests {
		output := interfaceToBytes(t.input)
		assert.Equal(t.output, output)
	}
}
