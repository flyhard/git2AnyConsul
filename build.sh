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
echo --------------------
echo Building GO project
echo --------------------
echo
gb build -f github.com/flyhard/gitProperties2Consul || exit 1
echo
echo --------------------
echo Building Dockerfile project
echo --------------------
echo
#docker build -t git2consul . || exit

#docker run --rm -it --link consul:consul git2consul $1
docker-compose build
docker-compose rm --force
docker-compose up