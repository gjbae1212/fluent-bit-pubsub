#!/bin/bash
set -e -o pipefail
trap '[ "$?" -eq 0 ] || echo "Error Line:<$LINENO> Error Function:<${FUNCNAME}>"' EXIT

export GO111MODULE=on
cd `dirname $0`
CURRENT=`pwd`

function build
{
   go build -buildmode=c-shared -o pubsub.so .
   local osname=`go env | grep GOOS | awk -F "=" '{print $2}' | sed 's/\"//g'`
   local archname=`go env | grep GOARCH | awk -F "=" '{print $2}' | sed 's/\"//g'`
   mkdir $CURRENT/bin/${osname}_${archname} || true
   mv pubsub.so pubsub.h $CURRENT/bin/${osname}_${archname}/
}

function build_linux
{
   local plugin=`docker ps | grep pubsub | wc -l`
   if [ $plugin -eq 1 ]
   then
      docker kill pubsub
   fi
   local osname=linux
   local archname=amd64
   mkdir $CURRENT/bin/${osname}_${archname} || true
   docker build --no-cache -t pubsub:latest -f Dockerfile .
   docker run -it --rm -d --name pubsub pubsub:latest /bin/bash
   docker cp pubsub:/go/pubsub/pubsub.so $CURRENT/bin/${osname}_${archname}/
   docker cp pubsub:/go/pubsub/pubsub.h $CURRENT/bin/${osname}_${archname}/
   docker kill pubsub
}

CMD=$1
shift
$CMD $*
