package routes

import (
	"net/http"

	"github.com/utkarshkrsingh/beacon/internal/handlers"
)

func NewRoutes(infoHandler *handlers.InfoHandler) http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /health", infoHandler.Health)

	return mux
}
