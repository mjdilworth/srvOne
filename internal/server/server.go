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
	port string
}

func New(p string) server {
	return server{port: p}
}

func (s server) Run(ctx context.Context) {

	srv := &http.Server{
		//Addr: ":8080",
		Addr: ":" + s.port,
	}

	//simple
	http.HandleFunc("/health/", s.Health)

	go func() {
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatal(err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		slog.Info("Error shutting down server", "err", err)
	}

}
