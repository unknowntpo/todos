package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/unknowntpo/todos/internal/logger"
)

func (app *application) serve() error {
	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", app.config.Port),
		Handler:      app.newRoutes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	shutdownErr := make(chan error)

	// Start background goroutine handler.
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	app.bg.Start(ctx)

	go func() {
		// Why we need buffered channel ?
		quit := make(chan os.Signal, 1)

		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

		s := <-quit
		logger.Log.PrintInfo("shutting down server", map[string]string{
			"signal": s.String(),
		})

		// do shutdown routine.
		ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
		defer cancel()

		err := srv.Shutdown(ctx)
		if err != nil {
			shutdownErr <- err
		}

		app.bg.Wait()
		shutdownErr <- nil
	}()

	logger.Log.PrintInfo("starting server", map[string]string{
		"addr": srv.Addr,
		"env":  app.config.Env,
	})

	err := srv.ListenAndServe()
	if err != nil {
		if !errors.Is(err, http.ErrServerClosed) {
			// Don't need to do graceful shutdown process, just return the error.
			return err
		}
	}

	err = <-shutdownErr
	if err != nil {
		return err
	}

	logger.Log.PrintInfo("stopped server", map[string]string{
		"addr": srv.Addr,
	})

	return nil
}
