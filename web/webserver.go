package web

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"

	"github.com/sh3rp/eyes/controller"
	"github.com/sh3rp/eyes/messages"
)

type Webserver struct {
	Controller  *controller.ProbeController
	ResultCache map[string]*messages.ProbeResult
}

func NewWebserver() *Webserver {
	controller := controller.NewProbeController()
	go controller.Start()
	return &Webserver{
		Controller:  controller,
		ResultCache: make(map[string]*messages.ProbeResult),
	}
}

func (ws *Webserver) Start() {
	ws.Controller.AddResultListener(ws.handleResult)
	log.Printf("Webserver starting")
	http.HandleFunc("/agents", ws.listAgents)
	http.HandleFunc("/agent.control/", ws.controlAgent)
	http.HandleFunc("/agent.test/", ws.testAgent)
	http.HandleFunc("/results", ws.listResults)
	http.HandleFunc("/results/", ws.showResult)
	http.ListenAndServe(":8080", nil)
}

func (ws *Webserver) handleResult(result *messages.ProbeResult) {
	ws.ResultCache[result.CmdId] = result
}

func (ws *Webserver) controlAgent(w http.ResponseWriter, r *http.Request) {

}

func (ws *Webserver) listAgents(w http.ResponseWriter, r *http.Request) {
	agents := ws.Controller.Agents
	w.Header().Set("Content-type", "application/json")
	json.NewEncoder(w).Encode(agents)
}

func (ws *Webserver) testAgent(w http.ResponseWriter, r *http.Request) {
	url := r.URL.Path
	elements := strings.Split(url, "/")
	id := elements[len(elements)-1]
	log.Printf("Testing agent: %s", id)
	ws.Controller.TestProbe(id)
}

func (ws *Webserver) listResults(w http.ResponseWriter, r *http.Request) {
	results := ws.ResultCache
	w.Header().Set("Content-type", "application/json")
	json.NewEncoder(w).Encode(results)
}

func (ws *Webserver) showResult(w http.ResponseWriter, r *http.Request) {
	url := r.URL.Path
	elements := strings.Split(url, "/")
	id := elements[len(elements)-1]
	result := ws.ResultCache[id]
	w.Header().Set("Content-type", "application/json")
	json.NewEncoder(w).Encode(result)
}
