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
}

type AgentControlResponse struct {
	StandardResponse
	Results []string `json:"results"`
}
