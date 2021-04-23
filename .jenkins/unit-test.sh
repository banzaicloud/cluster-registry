#!/bin/bash

set -x
set -e

# Dependency on banzaicloud/cluster-registry
export GOPRIVATE='github.cisco.com,github.com/banzaicloud'
export GOPROXY="https://proxy.golang.org,direct"

export GOPATH=$(go env GOPATH)
export PATH="${PATH}:${GOPATH}/bin"
export GOFLAGS='-mod=readonly'

go mod download

echo "Manifest generation"
make manifests 
