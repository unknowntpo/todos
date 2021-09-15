package main

import (
	"bytes"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/unknowntpo/todos/config"

	"github.com/spf13/viper"
)

var defaultConf = []byte(`
app:
  port: 4000
  env: development
  db:
    dsn: "postgres://todos:pa55word@localhost/todos?sslmode=disable"
    max_open_conns: 25
    max_idle_conns: 25
    max_idle_time: 15m
  limiter:
    rps: 2
    burst: 4
    enabled: True
  smtp:
    host: "smtp.mailtrap.io"
    port: 25
    username: "bd2857ac6e1116"
    password: "6f9845a2b11721"
    sender: "TODOs <no-reply@todos.unknowntpo.net>"
  cors:
    trusted_origins:
      - http://localhost:8080
`)

func setConfig() *config.Config {
	viper.SetConfigType("yml")
	viper.AutomaticEnv()        // read in environment variables that match
	viper.SetEnvPrefix("TODOS") // will be uppercased automatically
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	var configPath string
	flag.StringVar(&configPath, "c", "", "Configuration file path.")
	flag.Parse()

	if configPath != "" {
		content, err := ioutil.ReadFile(configPath)
		if err != nil {
			fmt.Errorf("failed to load config file: %v", err)
			os.Exit(1)
		}

		viper.ReadConfig(bytes.NewBuffer(content))
	} else {
		// At here, user doesn't specify config file location.
		// So we load default config.
		viper.ReadConfig(bytes.NewBuffer(defaultConf))
	}

	var cfg config.Config

	cfg.Port = viper.GetInt("app.port")
	//fmt.Println("is port set ?", viper.IsSet("app.port"))
	cfg.Env = viper.GetString("app.env")
	cfg.DB = config.DB{
		Dsn:          viper.GetString("app.db.dsn"),
		MaxOpenConns: viper.GetInt("app.db.max_open_conns"),
		MaxIdleConns: viper.GetInt("app.db.max_idle_conns"),
		MaxIdleTime:  viper.GetString("app.db.max_idle_time"),
	}
	cfg.Limiter = config.Limiter{
		Rps:     viper.GetFloat64("app.limiter.rps"),
		Burst:   viper.GetInt("app.limiter.burst"),
		Enabled: viper.GetBool("app.limiter.enabled"),
	}
	cfg.Smtp = config.Smtp{
		Host:     viper.GetString("app.smtp.host"),
		Port:     viper.GetInt("app.smtp.port"),
		Username: viper.GetString("app.smtp.username"),
		Password: viper.GetString("app.smtp.password"),
		Sender:   viper.GetString("app.smtp.sender"),
	}
	cfg.Cors = config.Cors{
		TrustedOrigins: viper.GetStringSlice("app.cors.trusted_origins"),
	}

	fmt.Println(cfg)
	return &cfg
}

// Note: Read config from flag, we left it for now.
/*
func setConfig() *config.Config {
	var cfg config.Config

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

	return &cfg
}
*/
