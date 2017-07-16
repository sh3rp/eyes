package main

import "github.com/sh3rp/eyes/web"

func main() {
	flags()
	webserver := web.NewWebserver()
	webserver.Start()
}

func flags() {

}
