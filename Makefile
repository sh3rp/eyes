all: protobuf glide test

test:
	go test -cover ./...

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
	protoc --proto_path=msg msg/agent.proto --go_out=msg

glide:
	glide update
	glide install

.PHONY: agent controller protobuf glide
