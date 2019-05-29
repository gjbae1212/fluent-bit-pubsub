# fluent-bit output plugin for google pubsub

<p align="left">    
  <a href="https://circleci.com/gh/gjbae1212/fluent-bit-pubsub/tree/master"><img src="https://circleci.com/gh/gjbae1212/fluent-bit-pubsub/tree/master.svg?style=svg"/></a>
  <a href="https://hits.seeyoufarm.com"/><img src="https://hits.seeyoufarm.com/api/count/incr/badge.svg?url=https%3A%2F%2Fgithub.com%2Fgjbae1212%2Ffluent-bit-pubsub"/></a>
  <a href="/LICENSE"><img src="https://img.shields.io/badge/license-MIT-GREEN.svg" alt="license" /></a>
  <a href="https://goreportcard.com/report/github.com/gjbae1212/fluent-bit-pubsub"><img src="https://goreportcard.com/badge/github.com/gjbae1212/fluent-bit-pubsub" alt="Go Report Card" /></a>
  <a href="https://codecov.io/gh/gjbae1212/fluent-bit-pubsub"><img src="https://codecov.io/gh/gjbae1212/fluent-bit-pubsub/branch/master/graph/badge.svg"/></a>        
</p>

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
| Key           | Description                                    | Default        |
| ----------------|------------------------------------------------|----------------|
| Project         | google cloud project id | NONE(required) |
| Topic           | google pubsub topic name | NONE(required) |
| JwtPath         | jwt file path for accessible google cloud project | NONE(required) |
| Debug           | print debug log | false(optional) |
| Timeout         | the maximum time that the client will attempt to publish a bundle of messages. (millsecond) | 60000 (optional)|
| DelayThreshold  | publish a non-empty batch after this delay has passed. (millsecond) | 1  |
| ByteThreshold   | publish a batch when its size in bytes reaches this value. | 1000000 |
| CountThreshold  | publish a batch when it has been reached count of messages. | 100  |

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
