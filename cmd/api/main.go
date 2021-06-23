package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/unknowntpo/todos/internal/logger"
	"github.com/unknowntpo/todos/internal/logger/logrus"
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
}

func main() {
	cfg := setConfig()

	err := logrus.RegisterLog()
	if err != nil {
		fmt.Errorf("Fail to set up logger: %v", err)
		os.Exit(1)
	}

	app := &application{
		config: cfg,
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
