package main

import (
	"flag"
	"fmt"

	"github.com/sh3rp/eyes/agent"
)

var controllerIP string
var label string

func main() {
	flags()
	fmt.Printf("%s\n", label)
	agent := agent.NewAgent(label)
	agent.Start(controllerIP)
}

func flags() {
	flag.StringVar(&controllerIP, "c", "", "IP address of the controller")
	flag.StringVar(&label, "l", "", "Label description for the agent")
	flag.Parse()
}
