package model

// Response is a model
type Response struct {
	IsSuccess bool        `json:"-"`
	Status    string      `json:"status"`
	Data      interface{} `json:"data"`
	Error     string      `json:"error"`
	Meta      interface{} `json:"meta,omitempty"`
}
