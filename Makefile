all: test agent controller

agent:
	go build -o eyes-agent cmd/agent/eyes-agent.go

test:
	go test agent/*
	go test controller/*
	go test web/*
	go test probe/*
	go test util/*

controller:
	go build -o eyes-controller cmd/controller/eyes-controller.go

.PHONY: agent controller
