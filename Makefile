all: protobuf agent controller client

agent:
	go build -o eyes-agent cmd/agent/eyes-agent.go

test:
	go test agent/*
	go test controller/*
	go test web/auth/*
	go test probe/*
	go test util/*

controller:
	go build -o eyes-controller cmd/controller/eyes-controller.go

client:
	go build -o eyes cmd/client/eyes.go

protobuf:
	protoc --proto_path=messages messages/messages.proto --go_out=plugins=grpc:messages

.PHONY: agent controller protobuf
