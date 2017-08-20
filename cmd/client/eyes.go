package main

import (
	"context"
	"fmt"

	"github.com/sh3rp/eyes/messages"
	"google.golang.org/grpc"
)

func main() {
	c, err := grpc.Dial("localhost:9999", grpc.WithInsecure())

	if err != nil {
		panic(err)
	}

	defer c.Close()

	client := messages.NewControllerClient(c)

	agents, err := client.GetAgents(context.Background(), &messages.AgentRequest{})

	if err != nil {
		panic(err)
	}

	for id, agent := range agents.Agents {
		fmt.Println(id, "=>", agent.Ipaddress)
	}
}
