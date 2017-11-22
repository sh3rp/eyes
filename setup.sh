#!/bin/sh

PROTOBUF_VERSION=3.5.0

wget https://github.com/google/protobuf/releases/download/v${PROTOBUF_VERSION}/protoc-${PROTOBUF_VERSION}-linux-x86_64.zip
sudo unzip protoc-${PROTOBUF_VERSION}-linux-x86_64.zip -d /usr/local && sudo chmod 755 /usr/local/bin/protoc
sudo add-apt-repository -y ppa:masterminds/glide && sudo apt-get update
sudo apt-get install -y glide
go get -u github.com/golang/protobuf/protoc-gen-go

