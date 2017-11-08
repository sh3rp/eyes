package cmd

import (
	"github.com/sh3rp/eyes/agent"
	"github.com/sh3rp/eyes/net"
	"github.com/spf13/cobra"
)

var AgentCmd = &cobra.Command{
	Use:   "agent",
	Short: "",
	Long:  "",
	Run: func(cmd *cobra.Command, args []string) {
		agt := agent.NewMemAgent()
		net.NewAgentServer(nil, agt)

	},
}
