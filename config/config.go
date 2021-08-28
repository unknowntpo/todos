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
