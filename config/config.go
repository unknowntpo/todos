package config

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

// Smtp is the configuraion used to set up the dialer in github.com/go-mail/mail/v2.
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
