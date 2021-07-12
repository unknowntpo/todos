package main

import (
	"context"
	"database/sql"
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/unknowntpo/todos/internal/logger"
	"github.com/unknowntpo/todos/internal/logger/logrus"

	_ "github.com/lib/pq"
)

var (
	version   string
	buildTime string
)

// config holds configuration of server
type config struct {
	port int
	env  string
	db   struct {
		dsn          string
		maxOpenConns int
		maxIdleConns int
		maxIdleTime  string
	}
}

// application holds the dependencies for our HTTP handlers, helpers, and middleware.
type application struct {
	config   config
	database *sql.DB
}

func main() {
	cfg := setConfig()

	err := logrus.RegisterLog()
	if err != nil {
		logger.Log.PrintFatal(fmt.Errorf("Fail to set up logger: %v", err), nil)
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
		logger.Log.PrintFatal(fmt.Errorf("Server error: %v", err), nil)
	}
}

// config set the configuration and  returns the config struct.
func setConfig() config {
	var cfg config

	displayVersion := flag.Bool("version", false, "Display version and exit")

	flag.IntVar(&cfg.port, "port", 4000, "API server port")
	flag.StringVar(&cfg.env, "env", "development", "Environment (development|staging|production)")

	// db setup
	// Use the empty string "" as the default value for the db-dsn command-line flag,
	// let Makefile specify it explicitly.
	flag.StringVar(&cfg.db.dsn, "db-dsn", "", "PostgreSQL DSN")

	// Read the connection pool settings from command-line flags into the config struct.
	flag.IntVar(&cfg.db.maxOpenConns, "db-max-open-conns", 25, "PostgreSQL max open connections")
	flag.IntVar(&cfg.db.maxIdleConns, "db-max-idle-conns", 25, "PostgreSQL max idle connections")
	flag.StringVar(&cfg.db.maxIdleTime, "db-max-idle-time", "15m", "PostgreSQL max connection idle time")

	flag.Parse()

	// If the version flag value is true, then print out the version number and
	// immediately exit.
	if *displayVersion == true {
		fmt.Printf("Version:\t%s\n", version)
		fmt.Printf("Build time:\t%s\n", buildTime)
		os.Exit(0)
	}

	return cfg
}

func openDB(cfg config) (*sql.DB, error) {
	db, err := sql.Open("postgres", cfg.db.dsn)
	if err != nil {
		return nil, err
	}

	db.SetMaxOpenConns(cfg.db.maxOpenConns)
	db.SetMaxIdleConns(cfg.db.maxIdleConns)

	duration, err := time.ParseDuration(cfg.db.maxIdleTime)
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
