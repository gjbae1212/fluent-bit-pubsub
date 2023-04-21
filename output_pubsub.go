package main

import (
	"C"
	"fmt"
	"os"
	"strconv"
	"time"
	"unsafe"

	"cloud.google.com/go/pubsub"

	"context"

	"github.com/fluent/fluent-bit-go/output"
)

var (
	plugin   Keeper
	hostname string
	wrapper  = OutputWrapper(&Output{})

	timeout        = pubsub.DefaultPublishSettings.Timeout
	delayThreshold = pubsub.DefaultPublishSettings.DelayThreshold
	countThreshold = pubsub.DefaultPublishSettings.CountThreshold
	byteThreshold  = pubsub.DefaultPublishSettings.ByteThreshold
	debug          = false

	schemaType       pubsub.SchemaType
	schemaDefinition string
)

type Output struct{}

type OutputWrapper interface {
	Register(ctx unsafe.Pointer, name string, desc string) int
	GetConfigKey(ctx unsafe.Pointer, key string) string
	NewDecoder(data unsafe.Pointer, length int) *output.FLBDecoder
	GetRecord(dec *output.FLBDecoder) (ret int, ts interface{}, rec map[interface{}]interface{})
}

func (o *Output) Register(ctx unsafe.Pointer, name string, desc string) int {
	return output.FLBPluginRegister(ctx, name, desc)
}

func (o *Output) GetConfigKey(ctx unsafe.Pointer, key string) string {
	return output.FLBPluginConfigKey(ctx, key)
}

func (o *Output) NewDecoder(data unsafe.Pointer, length int) *output.FLBDecoder {
	return output.NewDecoder(data, length)
}

func (o *Output) GetRecord(dec *output.FLBDecoder) (ret int, ts interface{}, rec map[interface{}]interface{}) {
	return output.GetRecord(dec)
}

//export FLBPluginRegister
func FLBPluginRegister(ctx unsafe.Pointer) int {
	return wrapper.Register(ctx, "pubsub", "output pubsub")
}

//export FLBPluginInit
func FLBPluginInit(ctx unsafe.Pointer) int {
	var err error
	project := wrapper.GetConfigKey(ctx, "Project")
	topic := wrapper.GetConfigKey(ctx, "Topic")
	jwtPath := wrapper.GetConfigKey(ctx, "JwtPath")
	dg := wrapper.GetConfigKey(ctx, "Debug")
	to := wrapper.GetConfigKey(ctx, "Timeout")
	bt := wrapper.GetConfigKey(ctx, "ByteThreshold")
	ct := wrapper.GetConfigKey(ctx, "CountThreshold")
	dt := wrapper.GetConfigKey(ctx, "DelayThreshold")
	st := wrapper.GetConfigKey(ctx, "SchemaType")
	sfp := wrapper.GetConfigKey(ctx, "SchemaFilePath")

	fmt.Printf("[pubsub-go] plugin parameter project = '%s'\n", project)
	fmt.Printf("[pubsub-go] plugin parameter topic = '%s'\n", topic)
	fmt.Printf("[pubsub-go] plugin parameter jwtPath = '%s'\n", jwtPath)
	fmt.Printf("[pubsub-go] plugin parameter debug = '%s'\n", dg)
	fmt.Printf("[pubsub-go] plugin parameter timeout = '%s'\n", to)
	fmt.Printf("[pubsub-go] plugin parameter byte threshold = '%s'\n", bt)
	fmt.Printf("[pubsub-go] plugin parameter count threshold = '%s'\n", ct)
	fmt.Printf("[pubsub-go] plugin parameter delay threshold = '%s'\n", dt)
	fmt.Printf("[pubsub-go] plugin parameter schema type = '%s'\n", st)
	fmt.Printf("[pubsub-go] plugin parameter schema definition file path = '%s'\n", sfp)

	hostname, err = os.Hostname()
	if err != nil {
		fmt.Printf("[err][init] %+v\n", err)
		return output.FLB_ERROR
	}

	fmt.Printf("[pubsub-go] plugin hostname = '%s'\n", hostname)

	if dg != "" {
		debug, err = strconv.ParseBool(dg)
		if err != nil {
			fmt.Printf("[err][init] %+v\n", err)
			return output.FLB_ERROR
		}
	}
	if to != "" {
		v, err := strconv.Atoi(to)
		if err != nil {
			fmt.Printf("[err][init] %+v\n", err)
			return output.FLB_ERROR
		}
		timeout = time.Duration(v) * time.Millisecond
	}
	if bt != "" {
		v, err := strconv.Atoi(bt)
		if err != nil {
			fmt.Printf("[err][init] %+v\n", err)
			return output.FLB_ERROR
		}
		byteThreshold = v
	}
	if ct != "" {
		v, err := strconv.Atoi(ct)
		if err != nil {
			fmt.Printf("[err][init] %+v\n", err)
			return output.FLB_ERROR
		}
		countThreshold = v
	}
	if dt != "" {
		v, err := strconv.Atoi(dt)
		if err != nil {
			fmt.Printf("[err][init] %+v\n", err)
			return output.FLB_ERROR
		}
		delayThreshold = time.Duration(v) * time.Millisecond
	}
	if st != "" {
		schemaType = pubsub.SchemaTypeUnspecified
		if st == "Avro" {
			schemaType = pubsub.SchemaAvro
		} else if st == "ProtocolBuffer" {
			schemaType = pubsub.SchemaProtocolBuffer
		} else {
			fmt.Printf("[err][init] Invalid SchemaType\n")
			return output.FLB_ERROR
		}
	}
	if sfp != "" {
		definitionSource, err := os.ReadFile(sfp)
		if err != nil {
			fmt.Printf("[err][init] %+v\n", err)
			return output.FLB_ERROR
		}
		schemaDefinition = string(definitionSource)
	}

	publishSetting := pubsub.PublishSettings{
		ByteThreshold:  byteThreshold,
		CountThreshold: countThreshold,
		DelayThreshold: delayThreshold,
		Timeout:        timeout,
	}
	schemaConfig := pubsub.SchemaConfig{
		Type:       schemaType,
		Definition: schemaDefinition,
	}

	keeper, err := NewKeeper(project, topic, jwtPath, &publishSetting, &schemaConfig)
	if err != nil {
		fmt.Printf("[err][init] %+v\n", err)
		return output.FLB_ERROR
	}
	plugin = keeper
	return output.FLB_OK
}

//export FLBPluginFlush
func FLBPluginFlush(data unsafe.Pointer, length C.int, tag *C.char) int {
	ctx := context.Background()
	tagname := ""
	if tag == nil {
		tagname = C.GoString(tag)
	}

	// Create Fluent Bit decoder
	dec := wrapper.NewDecoder(data, int(length))
	var results []*pubsub.PublishResult
	// Iterate Records
	for {
		// Extract Record
		ret, ts, record := wrapper.GetRecord(dec)
		if ret != 0 { // don't rest
			break
		}
		timestamp := ts.(output.FLBTime)

		if schemaType == pubsub.SchemaAvro || schemaType == pubsub.SchemaProtocolBuffer {
			r, err := plugin.InterfaceMapToByte(record)
			if err != nil {
				fmt.Printf("[err][publish][don't retry] %+v \n", err)
			}
			results = append(results, plugin.Send(ctx, r))
		} else {
			for k, v := range record {
				//fmt.Printf("[%s] %s %s %v \n", tagname, timestamp.String(), k, v)
				_, _, _ = k, timestamp, tagname
				results = append(results, plugin.Send(ctx, interfaceToBytes(v)))
			}
		}
	}
	for _, result := range results {
		if _, err := result.Get(ctx); err != nil {
			// if timeout is raised.
			if err == context.DeadlineExceeded || err == context.Canceled {
				fmt.Printf("[err][publish][retry] %+v \n", err)
				return output.FLB_RETRY
			}
			// else error is next
			fmt.Printf("[err][publish][don't retry] %+v \n", err)
		}
	}
	return output.FLB_OK
}

//export FLBPluginExit
func FLBPluginExit() int {
	plugin.Stop()
	return output.FLB_OK
}

func interfaceToBytes(v interface{}) []byte {
	switch d := v.(type) {
	case []byte:
		return d
	case string:
		return []byte(d)
	case int, int32, int64, uint, uint32, uint64:
		return []byte(fmt.Sprintf("%d", d))
	case float32, float64:
		return []byte(fmt.Sprintf("%f", d))
	case bool:
		return []byte(strconv.FormatBool(d))
	case time.Time:
		return []byte(d.Format(time.RFC3339))
	default:
		return []byte(fmt.Sprintf("%v", d))
	}
}

func main() {}
