all: test agent controller

agent:
	go build -o eyes-agent cmd/agent/eyes-agent.go

test:
	go test ./...

controller:
	go build -o eyes-controller cmd/controller/eyes-controller.go

.PHONY: agent controller