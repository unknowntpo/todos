package main

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/unknowntpo/todos/config"
	"github.com/unknowntpo/todos/internal/helpers/background"
	"github.com/unknowntpo/todos/internal/logger"
	"github.com/unknowntpo/todos/internal/logger/logrus"

	_ "github.com/lib/pq"
)

var (
	version   string
	buildTime string
)

// application holds the dependencies for our HTTP handlers, helpers, and middleware.
type application struct {
	config   *config.Config
	database *sql.DB
	bg       background.Background
}

func main() {
	cfg := setConfig()

	err := logrus.RegisterLog()
	if err != nil {
		logger.Log.PrintFatal(fmt.Errorf("fail to set up logger: %v", err), nil)
	}

	// set up db.
	db, err := openDB(cfg)
	if err != nil {
		logger.Log.PrintFatal(err, nil)
	}
	defer db.Close()

	app := &application{
		config:   cfg,
		database: db,
	}

	err = app.serve()
	if err != nil {
		logger.Log.PrintFatal(fmt.Errorf("server error: %v", err), nil)
	}
}

func openDB(cfg *config.Config) (*sql.DB, error) {
	db, err := sql.Open("postgres", cfg.DB.Dsn)
	if err != nil {
		return nil, err
	}

	db.SetMaxOpenConns(cfg.DB.MaxOpenConns)
	db.SetMaxIdleConns(cfg.DB.MaxIdleConns)

	duration, err := time.ParseDuration(cfg.DB.MaxIdleTime)
	if err != nil {
		return nil, err
	}

	db.SetConnMaxIdleTime(duration)

	// Create a context with a 5-second timeout deadline.
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Use PingContext() to establish a new connection to the database, passing in the
	// context we created above as a parameter. If the connection couldn't be
	// established successfully within the 5 second deadline, then this will return an
	// error.
	err = db.PingContext(ctx)
	if err != nil {
		return nil, err
	}

	return db, nil
}
