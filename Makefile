all: protobuf glide agent controller client

test:
	go test agent/*
	go test controller/*
	go test web/auth/*
	go test probe/*
	go test util/*

agent:
	GOOS=darwin GOARCH=amd64 go build -o eyes-agent.darwin cmd/agent/eyes-agent.go
	GOOS=linux GOARCH=amd64 go build -o eyes-agent.linux cmd/agent/eyes-agent.go

controller:
	GOOS=darwin GOARCH=amd64 go build -o eyes-controller.darwin cmd/controller/eyes-controller.go
	GOOS=linux GOARCH=amd64 go build -o eyes-controller.linux cmd/controller/eyes-controller.go

client:
	GOOS=darwin GOARCH=amd64 go build -o eyes.darwin cmd/client/eyes.go
	GOOS=linux GOARCH=amd64 go build -o eyes.linux cmd/client/eyes.go

protobuf:
	protoc --proto_path=messages messages/messages.proto --go_out=plugins=grpc:messages

glide:
	glide update
	glide install

.PHONY: agent controller protobuf glide
