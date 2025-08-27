package configuration

import (
	"log"

	"github.com/joeshaw/envdecode"
)

type Conf struct {
	Server      ConfServer
	DB          ConfDB
	RateLimiter ConfRL
	Mail        ConfMAIL
	Security    ConfSecurity
}

type ConfServer struct {
	Port  int  `env:"SERVER_PORT,required"`
	Debug bool `env:"SERVER_DEBUG,required"`
}

type ConfDB struct {
	DSN          string `env:"DB_DSN,required"`
	MaxOpenConns int    `env:"DB_MAX_OPEN_CONNS,required"`
	MaxIdleConns int    `env:"DB_MAX_IDLE_CONNS,required"`
	MaxIdleTime  string `env:"DB_MAX_IDLE_TIME,required"`
}

type ConfRL struct {
	RPS     float64 `env:"LIMITER_RPS,required"`
	Burst   int     `env:"LIMITER_BURST,required"`
	Enabled bool    `env:"LIMITER_ENABLED,required"`
}

type ConfMAIL struct {
	HOST     string `env:"EMAIL_HOST,required"`
	PORT     int    `env:"EMAIL_POST,required"`
	USERNAME string `env:"EMAIL_USERNAME,required"`
	PASSWORD string `env:"EMAIL_PASSWORD,required"`
}

type ConfSecurity struct {
	SecretKey string `env:"SECRET_KEY,required"`
}

func New() *Conf {
	var c Conf
	if err := envdecode.StrictDecode(&c); err != nil {
		log.Fatalf("Failed to decode: %s", err)
	}

	return &c
}

func NewDB() *ConfDB {
	var c ConfDB
	if err := envdecode.StrictDecode(&c); err != nil {
		log.Fatalf("Failed to decode: %s", err)
	}

	return &c
}
