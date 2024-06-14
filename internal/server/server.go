package server

import (
	"context"
	"errors"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"time"
)

type server struct {
	port   string
	router *http.ServeMux
}

func New(p string) *server {
	mux := http.NewServeMux()
	s := server{
		port:   p,
		router: mux,
	}

	// setup routes
	//	s.routes()

	return &s
}

func (s server) Run(ctx context.Context) {

	srv := &http.Server{
		//Addr: ":8080",
		Addr: ":" + s.port,
	}

	//simple
	http.HandleFunc("/health/", s.Health)

	shutdownChan := make(chan bool, 1)

	go func() {
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatal(err)
		}
		slog.Info("Stopped serving new connections.")
		shutdownChan <- true
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		slog.Info("Error shutting down server", "err", err)
	}

	<-shutdownChan
	slog.Info("Graceful shutdown complete.")

}
