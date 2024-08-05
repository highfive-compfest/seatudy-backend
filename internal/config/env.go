package config

import (
	"log"
	"os"
	"time"
)

type environmentVariables struct {
	ENV  string
	Port string

	DbDsn string

	RedisAddress  string
	RedisPassword string
	RedisDatabase int

	JwtAccessSecret    []byte
	JwtAccessDuration  time.Duration
	JwtRefreshSecret   []byte
	JwtRefreshDuration time.Duration
}

var Env *environmentVariables

func LoadEnv() {
	env := &environmentVariables{}
	var err error

	env.ENV = os.Getenv("ENV")
	if env.ENV == "" {
		log.Fatal("ENV is not set")
	}
	env.Port = os.Getenv("PORT")

	env.DbDsn = os.Getenv("DB_DSN")

	env.RedisAddress = os.Getenv("REDIS_ADDRESS")
	env.RedisPassword = os.Getenv("REDIS_PASSWORD")
	env.RedisDatabase, err = strconv.Atoi(os.Getenv("REDIS_DATABASE"))
	if err != nil {
		log.Fatal("Fail to parse REDIS_DATABASE")
	}

	env.JwtAccessSecret = []byte(os.Getenv("JWT_ACCESS_SECRET"))
	env.JwtAccessDuration, err = time.ParseDuration(os.Getenv("JWT_ACCESS_DURATION"))
	if err != nil {
		log.Fatal("Fail to parse JWT_ACCESS_DURATION")
	}

	env.JwtRefreshSecret = []byte(os.Getenv("JWT_REFRESH_SECRET"))
	env.JwtRefreshDuration, err = time.ParseDuration(os.Getenv("JWT_REFRESH_DURATION"))
	if err != nil {
		log.Fatal("Fail to parse JWT_REFRESH_DURATION")
	}

	Env = env
}
