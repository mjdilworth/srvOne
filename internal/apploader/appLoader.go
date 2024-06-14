package apploader

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
)

type AppLoader struct {
	app    App
	ctx    context.Context
	cancel context.CancelFunc
}

func New(app App) *AppLoader {

	ctx, cancel := context.WithCancel(context.Background())

	return &AppLoader{
		app:    app,
		ctx:    ctx,
		cancel: cancel,
	}
}

func (al AppLoader) LogMe(s string, k string, v int) {
	slog.Info(s, k, v)
	// log.Fatalf("Fatal string %s, key %s and value %d", s, k, v)
}

func (al AppLoader) Run() error {
	defer al.cancel()

	exitCh := make(chan struct{})

	go func() {
		// This starts the app
		// and will wait until the child process is done
		al.app.Run(al.ctx)

		exitCh <- struct{}{}
	}()

	termCh := make(chan os.Signal, 3)
	signal.Notify(termCh, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-termCh
		al.cancel()
	}()

	<-exitCh
	return nil
}
