package main

import (
	"context"
	"database/sql"
	"expvar"
	"fmt"
	"os"
	"runtime"
	"time"

	"github.com/unknowntpo/todos/config"
	"github.com/unknowntpo/todos/internal/mailer"
	"github.com/unknowntpo/todos/pkg/naivepool"

	"github.com/unknowntpo/todos/internal/logger"
	"github.com/unknowntpo/todos/internal/logger/zerolog"

	_ "github.com/lib/pq"
)

var (
	version string
	//lint:ignore U1000 buildTime is used during compilation at -ldflags.
	buildTime string
)

// application holds the dependencies for our HTTP handlers, helpers, and middleware.
type application struct {
	config   *config.Config
	database *sql.DB
	pool     *naivepool.Pool
	mailer   *mailer.Mailer
	logger   logger.Logger
}

// @title TODOS API
// @version 1.0
// @description This is the api documentation of TODOS server.
// @host localhost:4000
// @BasePath /
func main() {
	cfg := setConfig()

	// Use zerolog
	logger := zerolog.New(os.Stdout)

	// set up db.
	db, err := openDB(cfg)
	if err != nil {
		logger.PrintFatal(err, nil)
	}
	defer db.Close()

	// Set up workerpool with max jobs and max workers.
	// TODO: Do this in config.yml
	maxJobs := 100
	maxWorkers := 10
	workerChanSize := 10
	pool := naivepool.New(maxJobs, maxWorkers, workerChanSize)

	// Publish a new "version" variable in the expvar handler containing our application
	// version number.
	expvar.NewString("version").Set(version)

	// Publish the number of active goroutines.
	expvar.Publish("goroutines", expvar.Func(func() interface{} {
		return runtime.NumGoroutine()
	}))

	// Publish the database connection pool statistics.
	expvar.Publish("database", expvar.Func(func() interface{} {
		return db.Stats()
	}))

	// Publish the current Unix timestamp.
	expvar.Publish("timestamp", expvar.Func(func() interface{} {
		return time.Now().Unix()
	}))

	app := &application{
		config:   cfg,
		database: db,
		pool:     pool,
		mailer:   mailer.New(&cfg.Smtp),
		logger:   logger,
	}

	err = app.serve()
	if err != nil {
		logger.PrintFatal(fmt.Errorf("server error: %v", err), nil)
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
