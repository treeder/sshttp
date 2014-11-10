#!/bin/bash
set -eu

# This script makes certain assumptions about your deployment of sshttp
# - You are deploying to at least Ubuntu 14.04 LTS
# - User is authorized to read github.com/treeder/sshttp
# - Go is installed

export GOPATH=$HOME/go
export PATH=$GOPATH/bin:$PATH

repo="github.com/treeder/sshttp"
go get $repo
go install $repo

go run $(go env GOROOT)/src/pkg/crypto/tls/generate_cert.go --host localhost
(sshttp -ssl -p ${SS_PORT} -t ${SS_TOKEN} &) &
