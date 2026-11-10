package handlers

import (
	"fmt"
	"net/http"
)

type InfoHandler struct{}

func NewInfoHandler() *InfoHandler {
	return &InfoHandler{}
}

func (h *InfoHandler) Health(w http.ResponseWriter, r *http.Request) {
	_, err := fmt.Fprint(w, "server running in development mode")
	if err != nil {
		http.Error(w, "unable to write the server health status", http.StatusInternalServerError)
	}
}
