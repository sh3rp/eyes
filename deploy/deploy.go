package deploy

import "github.com/sh3rp/eyes/agent"

type Deployments struct {
	agents []agent.Agent
}

func (d *Deployments) AddAgent(agent agent.Agent) error {
	d.agents = append(d.agents, agent)
	return nil
}
