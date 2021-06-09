package main

import (
	"expvar"
	"flag"
	"fmt"
	"os"

	"github.com/unknowntpo/todos/internal/jsonlog"
)

var (
	version   string
	buildTime string
)

// config holds configuration of server
type config struct {
	port int
	env  string
}

// application holds the dependencies for our HTTP handlers, helpers, and middleware.
type application struct {
	config config
	logger *jsonlog.Logger
}

func main() {
	cfg := setConfig()

	logger := jsonlog.New(os.Stdout, jsonlog.LevelInfo)

	// Publish a new "version" variable in the expvar handler containing our application
	// version number (currently the constant "1.0.0").
	expvar.NewString("version").Set(version)

	app := &application{
		config: cfg,
		logger: logger,
	}

	err := app.serve()
	if err != nil {
		logger.PrintFatal(err, nil)
	}
}

// config set the configuration and  returns the config struct.
func setConfig() config {
	var cfg config

	flag.IntVar(&cfg.port, "port", 4000, "API server port")
	flag.StringVar(&cfg.env, "env", "development", "Environment (development|staging|production)")

	displayVersion := flag.Bool("version", false, "Display version and exit")

	// If the version flag value is true, then print out the version number and
	// immediately exit.
	if *displayVersion {
		fmt.Printf("Version:\t%s\n", version)
		fmt.Printf("Build time:\t%s\n", buildTime)
		os.Exit(0)
	}

	flag.Parse()

	return cfg
}
