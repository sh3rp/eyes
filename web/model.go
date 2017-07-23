package web

type StandardResponse struct {
	StatusCode    int    `json:"code"`
	StatusMessage string `json:"message"`
}

type AgentControlRequest struct {
	Agents   []string       `json:"agents"`
	Type     string         `json:"type"`
	Host     string         `json:"host"`
	Schedule *AgentSchedule `json:"schedule"`
}

type AgentSchedule struct {
	Repeat bool
}

type AgentControlResponse struct {
	StandardResponse
	Results []string `json:"results"`
}

type ResultResponse struct {
	StandardResponse
	ResultId      string
	AgentId       string
	AgentLabel    string
	AgentLocation string
	Timestamp     int64
	Datapoints    []int64
}
