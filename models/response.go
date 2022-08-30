package models

type Response struct {
	Code    string      `json:"code"`
	Message string      `json:"message"`
	Results interface{} `json:"results,omitempty"`
}
