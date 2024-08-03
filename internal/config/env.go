package config

import (
	"log"
	"os"
	"time"
)

type environmentVariables struct {
	ENV                string
	Port               string
	DbDsn              string
	JwtAccessSecret    []byte
	JwtAccessDuration  time.Duration
	JwtRefreshSecret   []byte
	JwtRefreshDuration time.Duration
}

var Env *environmentVariables

func LoadEnv() {
	env := &environmentVariables{}

	env.ENV = os.Getenv("ENV")
	if env.ENV == "" {
		log.Fatal("ENV is not set")
	}

	env.Port = os.Getenv("PORT")
	env.DbDsn = os.Getenv("DB_DSN")

	env.JwtAccessSecret = []byte(os.Getenv("JWT_ACCESS_SECRET"))
	dur, err := time.ParseDuration(os.Getenv("JWT_ACCESS_DURATION"))
	if err != nil {
		log.Fatal("Fail to parse JWT_ACCESS_DURATION")
	}
	env.JwtAccessDuration = dur

	env.JwtRefreshSecret = []byte(os.Getenv("JWT_REFRESH_SECRET"))
	dur, err = time.ParseDuration(os.Getenv("JWT_REFRESH_DURATION"))
	if err != nil {
		log.Fatal("Fail to parse JWT_REFRESH_DURATION")
	}

	Env = env
}
