package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/utkarshkrsingh/beacon/internal/handlers"
	"github.com/utkarshkrsingh/beacon/internal/routes"
)

func createServer() (*http.Server, error) {
	infoHandler := handlers.NewInfoHandler()

	mux := routes.NewRoutes(infoHandler)
	return &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}, nil
}

func runServer(
	ctx context.Context,
	server *http.Server,
	shutdownTimeout time.Duration,
) error {
	serverErr := make(chan error, 1)

	go func() {
		log.Println("Starting server at port :8080")
		if err := server.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
			serverErr <- err
		}
		close(serverErr)
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
	defer signal.Stop(stop)

	select {
	case err := <-serverErr:
		return err
	case <-stop:
		log.Println("shutdown signal recieved...")
	case <-ctx.Done():
		log.Println("context cancelled...")
	}

	shutdownContext, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancel()

	if err := server.Shutdown(shutdownContext); err != nil {
		if closeErr := server.Close(); closeErr != nil {
			return errors.Join(err, closeErr)
		}
		return err
	}

	log.Println("server exited gracefully...")
	return nil
}

func main() {
	server, err := createServer()
	if err != nil {
		log.Fatalf("server error: %v", err.Error())
	}

	if err := runServer(context.Background(), server, 3*time.Second); err != nil {
		log.Fatalf("server error: %v", err.Error())
	}
}
