package util

type StatusResponse struct {
	Status  string `json:"status"`
	Message string `json:"message,omitempty"`
}
