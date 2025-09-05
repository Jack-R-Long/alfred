package main

type HealthResponse struct {
	Status string `json:"status"`
	Data  struct {
		Message string `json:"message"`
	} `json:"data"`
}