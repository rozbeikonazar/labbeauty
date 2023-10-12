package main

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func (app *application) serve() error {
	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", app.config.port),
		Handler:      app.routes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		ErrorLog:     slog.NewLogLogger(app.logger.Handler(), slog.LevelError),
	}
	shutdownErr := make(chan error)
	// start a background goroutine
	go func() {
		quit := make(chan os.Signal, 1)
		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

		s := <-quit

		app.logger.Info("shutting down server", "signal", s.String())
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()
		err := srv.Shutdown(ctx)
		if err != nil {
			shutdownErr <- err
		}
		app.logger.Info("completing background tasks", "addr", srv.Addr)
		app.wg.Wait()
		shutdownErr <- nil
	}()

	app.logger.Info("Starting server", "addr", srv.Addr, "env", app.config.env)

	err := srv.ListenAndServe()
	if !errors.Is(err, http.ErrServerClosed) {
		return err
	}

	err = <-shutdownErr
	if err != nil {
		return err
	}

	app.logger.Info("stopped server", "addr", srv.Addr)
	return nil
}
