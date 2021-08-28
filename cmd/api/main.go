// TODOS API Documentation
//
// Terms Of Service:
//
// there are no TOS at this moment, use at your own risk we take no responsibility
//
//     Schemes: http, https
//     Host: localhost:4000
//     BasePath: /
//     Version: 1.0.0
//     License: MIT http://opensource.org/licenses/MIT
//     Contact: Chang Chen Chien<e850506@gmail.com>
//     Consumes:
//     - application/json
//
//     Produces:
//     - application/json
//
//     Security:
//     - api_key:
//
//     SecurityDefinitions:
//     api_key:
//          type: apiKey
//          name: KEY
//          in: header
// swagger:meta
package main

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/unknowntpo/todos/internal/logger"
	"github.com/unknowntpo/todos/internal/logger/logrus"

	"github.com/unknowntpo/todos/internal/helpers/background"

	_ "github.com/lib/pq"
)

var (
	version   string
	buildTime string
)

// application holds the dependencies for our HTTP handlers, helpers, and middleware.
type application struct {
	config   Config
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

func openDB(cfg Config) (*sql.DB, error) {
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
