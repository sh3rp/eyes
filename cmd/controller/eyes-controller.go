package main

import (
	"github.com/rs/zerolog/log"
	"github.com/sh3rp/eyes/controller"
)

func main() {
	flags()
	c := controller.NewProbeController()
	go c.Start()
	maj, min, pat := c.GetVersion()
	log.Info().Msgf("Net-Eyes controller v%d.%d.%d", maj, min, pat)
	server := controller.NewGRPCServer(9999, c)
	server.Start()
}

func flags() {

}
