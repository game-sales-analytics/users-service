#!/bin/sh

set -ex

protoc --go_out=./internal  --go-grpc_out=./internal api/api.proto
