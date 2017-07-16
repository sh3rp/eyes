package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/sh3rp/eyes/agent"
)

var controllerIP string
var label string

func main() {
	flags()
	if controllerIP == "" {
		fmt.Println("You must specify a controller IP.")
		os.Exit(1)
	}
	if label == "" {
		fmt.Println("You must specify a descriptive label.")
		os.Exit(1)
	}
	agent := agent.NewAgent(label)
	agent.Start(controllerIP)
}

func flags() {
	flag.StringVar(&controllerIP, "c", "", "IP address of the controller")
	flag.StringVar(&label, "l", "", "Label description for the agent")
	flag.Parse()
}
