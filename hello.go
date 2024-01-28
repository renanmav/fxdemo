package main

import (
	"fmt"
	"go.uber.org/zap"
	"io"
	"net/http"
)

type HelloHandler struct {
	log *zap.Logger
}

func (h *HelloHandler) Pattern() string {
	return "/hello"
}

func NewHelloHandler(log *zap.Logger) *HelloHandler {
	return &HelloHandler{log}
}

func (h *HelloHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		h.log.Warn("Failed to read request body", zap.Error(err))
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	if _, err := fmt.Fprintf(w, "Hello, %s\n", body); err != nil {
		h.log.Warn("Failed to write response", zap.Error(err))
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
}
