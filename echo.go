package main

import (
	"go.uber.org/zap"
	"io"
	"net/http"
)

type EchoHandler struct {
	log *zap.Logger
}

func (h *EchoHandler) Pattern() string {
	return "/echo"
}

func NewEchoHandler(log *zap.Logger) *EchoHandler {
	return &EchoHandler{log}
}

func (h *EchoHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if _, err := io.Copy(w, r.Body); err != nil {
		h.log.Warn("Failed to handle request", zap.Error(err))
	}
}
