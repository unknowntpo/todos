package main

import (
	"flag"
	"fmt"
	"os"
	"strings"
)

// Config holds configuration of server
type Config struct {
	Port    int
	Env     string
	DB      DB
	Limiter Limiter
	Smtp    Smtp
	Cors    Cors
}

type DB struct {
	Dsn          string
	MaxOpenConns int
	MaxIdleConns int
	MaxIdleTime  string
}
type Limiter struct {
	Rps     float64
	Burst   int
	Enabled bool
}
type Smtp struct {
	Host     string
	Port     int
	Username string
	Password string
	Sender   string
}
type Cors struct {
	TrustedOrigins []string
}

// config set the configuration and  returns the config struct.
func setConfig() Config {
	var cfg Config

	displayVersion := flag.Bool("version", false, "Display version and exit")

	flag.IntVar(&cfg.Port, "port", 4000, "API server port")
	flag.StringVar(&cfg.Env, "env", "development", "Environment (development|staging|production)")

	// db setup
	// Use the empty string "" as the default value for the db-dsn command-line flag,
	// let Makefile specify it explicitly.
	flag.StringVar(&cfg.DB.Dsn, "db-dsn", "", "PostgreSQL DSN")

	// Read the connection pool settings from command-line flags into the config struct.
	flag.IntVar(&cfg.DB.MaxOpenConns, "db-max-open-conns", 25, "PostgreSQL max open connections")
	flag.IntVar(&cfg.DB.MaxIdleConns, "db-max-idle-conns", 25, "PostgreSQL max idle connections")
	flag.StringVar(&cfg.DB.MaxIdleTime, "db-max-idle-time", "15m", "PostgreSQL max connection idle time")

	// Config the rate limiter.
	flag.Float64Var(&cfg.Limiter.Rps, "limiter-rps", 2, "Rate limiter maximum requests per second")
	flag.IntVar(&cfg.Limiter.Burst, "limiter-burst", 4, "Rate limiter maximum burst")
	flag.BoolVar(&cfg.Limiter.Enabled, "limiter-enabled", true, "Enable rate limiter")

	// Config the SMTP server.
	flag.StringVar(&cfg.Smtp.Host, "smtp-host", "smtp.mailtrap.io", "SMTP host")
	flag.IntVar(&cfg.Smtp.Port, "smtp-port", 25, "SMTP port")
	flag.StringVar(&cfg.Smtp.Username, "smtp-username", "bd2857ac6e1116", "SMTP username")
	flag.StringVar(&cfg.Smtp.Password, "smtp-password", "6f9845a2b11721", "SMTP password")
	flag.StringVar(&cfg.Smtp.Sender, "smtp-sender", "TODOs <no-reply@todos.unknowntpo.net>", "SMTP sender")

	// Config the cors trusted origins.
	flag.Func("cors-trusted-origins", "Trusted CORS origins (space separated)", func(val string) error {
		cfg.Cors.TrustedOrigins = strings.Fields(val)
		return nil
	})

	flag.Parse()

	// If the version flag value is true, then print out the version number and
	// immediately exit.
	if *displayVersion {
		fmt.Printf("Version:\t%s\n", version)
		fmt.Printf("Build time:\t%s\n", buildTime)
		os.Exit(0)
	}

	return cfg
}
