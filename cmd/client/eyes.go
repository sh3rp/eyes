package main

import (
	"context"
	"fmt"
	"strconv"

	"github.com/sh3rp/eyes/config"
	"github.com/sh3rp/eyes/messages"
	"google.golang.org/grpc"
)

func main() {
	cfg := &config.ClientConfig{}
	cfg.Read("./eyes-client.json")

	var connections []*grpc.ClientConn

	controllers := make(map[string]messages.ControllerClient)

	for _, host := range cfg.ControllerHosts {
		c, err := grpc.Dial(host.Host+":"+strconv.Itoa(host.Port), grpc.WithInsecure())
		if err != nil {
			fmt.Printf("%v\n", err)
		}
		connections = append(connections, c)
		controllers[host.Label] = messages.NewControllerClient(c)
	}

	defer func() {
		for _, c := range connections {
			c.Close()
		}
	}()

	for k, v := range controllers {
		controllerInfo, err := v.GetControllerInfo(context.Background(), &messages.Empty{})

		if err != nil {
			panic(err)
		}

		fmt.Printf("%s - controller %s\n\n", k, controllerInfo.Version)

		agents, err := v.GetAgents(context.Background(), &messages.AgentRequest{})

		if err != nil {
			panic(err)
		}

		for id, agent := range agents.Agents {
			fmt.Printf("\t[%s]\t%s\t%s\t%s\t%s\n", id, agent.Label, agent.Ipaddress, agent.Os, agent.Location)
		}
	}

}
