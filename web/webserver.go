package web

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/rs/zerolog/log"
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
	log.Info().Msgf("Webserver starting")
	http.HandleFunc("/agents", ws.listAgents)
	http.HandleFunc("/agent.control/", ws.controlAgent)
	http.HandleFunc("/agent.test/", ws.testAgent)
	http.HandleFunc("/results", ws.listResults)
	http.HandleFunc("/results/", ws.showResult)

	http.HandleFunc("/html/", ws.serveFile)
	http.HandleFunc("/js/", ws.serveFile)
	http.HandleFunc("/css/", ws.serveFile)
	http.ListenAndServe(":8080", nil)
}

func (ws *Webserver) serveFile(w http.ResponseWriter, r *http.Request) {
	url := r.URL.Path
	filePath := url[1:]
	data, err := ioutil.ReadFile("static/" + filePath)

	filePathSplit := strings.Split(filePath, ".")
	fileType := filePathSplit[len(filePathSplit)-1]
	log.Info().Msgf("Loading: %s", filePath)
	switch fileType {
	case "js":
		w.Header().Set("Content-Type", "text/javascript")
	case "html":
		w.Header().Set("Content-Type", "text/html")
	case "css":
		w.Header().Set("Content-Type", "text/html")
	}

	if err != nil {
		return
	}

	w.Write(data)
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
	log.Info().Msgf("Testing agent: %s", id)
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
