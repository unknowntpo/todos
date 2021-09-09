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

	// Starting worker pool.
	poolCtx, poolCancel := context.WithCancel(context.Background())
	defer poolCancel()

	app.pool.Start(poolCtx)

	go func() {
		// Why we need buffered channel ?
		quit := make(chan os.Signal, 1)

		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

		s := <-quit

		// Shutdown worker pool.
		poolCancel()

		app.logger.PrintInfo("shutting down server", map[string]interface{}{
			"signal": s.String(),
		})

		// do server shutdown routine.
		serverCtx, serverCancel := context.WithTimeout(context.Background(), 3*time.Second)
		defer serverCancel()

		err := srv.Shutdown(serverCtx)
		if err != nil {
			shutdownErr <- err
		}

		app.pool.Wait()
		shutdownErr <- nil
	}()

	app.logger.PrintInfo("starting server", map[string]interface{}{
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

	app.logger.PrintInfo("stopped server", map[string]interface{}{
		"addr": srv.Addr,
	})

	return nil
}
