all: agent controller

agent: eyes-agent
	go build cmd/agent/eyes-agent.go

controller: eyes-controller
	go build cmd/controller/eyes-controller.go