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
	AwsAccessId        string
	AmsSecretAccessId  string
	AwsRegion          string
	AwsBucketName      string
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

	env.AwsAccessId = os.Getenv("AWS_ACCESS_KEY_ID")
	env.AmsSecretAccessId = os.Getenv("AWS_SECRET_ACCESS_KEY")
	env.AwsRegion = os.Getenv("AWS_REGION")
	env.AwsBucketName = os.Getenv("AWS_BUCKET_NAME")
	if env.AwsBucketName == "" {
		log.Fatalf("AWS_BUCKET_NAME is not set in the environment variables")
	}

	Env = env
}
