package handlers

import (
	"net/http"
)

type healthCheckHandler struct{}

func NewHealthCheckHandler() Handler {
	return &healthCheckHandler{}
}

func (h healthCheckHandler) Handle(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(200)
	w.Write([]byte("ok"))
}
