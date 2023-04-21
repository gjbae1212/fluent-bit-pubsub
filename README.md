# fluent-bit output plugin for google pubsub

![example workflow](https://github.com/ragi256/fluent-bit-pubsub/actions/workflows/ci.yml/badge.svg?branch=master)

This plugin is used to publish data to queue in google pubsub. 

You could easily use it.

## Build
A bin directory already has been made binaries for mac, linux.

If you should directly make binaries for mac, linux
```bash
# local machine binary
$ bash make.sh build

# Your machine is mac, and if you should do to retry cross compiling for linux.
# A command in below is required a docker.
$ bash make.sh build_linux
```

## Usage
### configuration options for fluent-bit.conf
| Key             | Description                                    | Default        |
| ----------------|------------------------------------------------|----------------|
| Project         | google cloud project id | NONE(required) |
| Topic           | google pubsub topic name | NONE(required) |
| JwtPath         | jwt file path for accessible google cloud project | NONE(required) |
| Debug           | print debug log | false(optional) |
| Timeout         | the maximum time that the client will attempt to publish a bundle of messages. (millsecond) | 60000 (optional)|
| DelayThreshold  | publish a non-empty batch after this delay has passed. (millsecond) | 1  |
| ByteThreshold   | publish a batch when its size in bytes reaches this value. | 1000000 |
| CountThreshold  | publish a batch when it has been reached count of messages. | 100  |
| SchemaType      | topic schema type Avro ~~or Protocol Buffer~~. (Only support Avro yet) | NONE(optional) |
| SchemaFilePath  | schema definition file path. | NONE(optional) |

### Example fluent-bit.conf
```conf
[Output]
    Name pubsub
    Match *
    Project your-project(custom)
    Topic your-topic-name(custom)
    Jwtpath your-jwtpath(custom)
```

### Example exec
```bash
$ fluent-bit -c [your config file] -e pubsub.so
```
