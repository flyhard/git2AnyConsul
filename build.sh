#!/usr/bin/env bash

if [[ -z $GOPATH ]]
then
    export GOPATH=$HOME/go
    mkdir -p $GOPATH
fi
if [[ -z $(which gb) ]]
then
    go get github.com/constabulary/gb/...
    go clean
    go install -v github.com/constabulary/gb/
    sudo cp $GOPATH/bin/gb* /usr/local/bin
fi

gb build github.com/flyhard/gitProperties2Consul