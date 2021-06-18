package main

import (
	"flag"
	"fmt"
	"os"

	log "github.com/sirupsen/logrus"
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
	config  config
	infoLog *log.Logger
	errLog  *log.Logger
}

func main() {
	cfg := setConfig()

	infoLog := &log.Logger{
		Out:       os.Stdout,
		Formatter: new(log.JSONFormatter),
		Hooks:     make(log.LevelHooks),
		Level:     log.DebugLevel,
	}

	errLog := &log.Logger{
		Out:       os.Stderr,
		Formatter: new(log.JSONFormatter),
		Hooks:     make(log.LevelHooks),
		Level:     log.DebugLevel,
	}

	app := &application{
		config:  cfg,
		infoLog: infoLog,
		errLog:  errLog,
	}

	err := app.serve()
	if err != nil {
		app.errLog.Fatal(err)
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
