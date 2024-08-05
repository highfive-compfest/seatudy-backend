package config

import (
	"log"
	"os"
	"strconv"
	"time"
)

type environmentVariables struct {
	ENV         string
	FrontendUrl string
	Port        string

	DbDsn string

	RedisAddress  string
	RedisPassword string
	RedisDatabase int

	JwtAccessSecret    []byte
	JwtAccessDuration  time.Duration
	JwtRefreshSecret   []byte
	JwtRefreshDuration time.Duration

	AwsAccessId       string
	AmsSecretAccessId string
	AwsRegion         string
	AwsBucketName     string

	SmtpHost     string
	SmtpPort     int
	SmtpUsername string
	SmtpEmail    string
	SmtpPassword string
}

var Env *environmentVariables

func LoadEnv() {
	env := &environmentVariables{}
	var err error

	env.ENV = os.Getenv("ENV")
	if env.ENV == "" {
		log.Fatal("ENV is not set")
	}
	env.FrontendUrl = os.Getenv("FRONTEND_URL")
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

	env.AwsAccessId = os.Getenv("AWS_ACCESS_KEY_ID")
	env.AmsSecretAccessId = os.Getenv("AWS_SECRET_ACCESS_KEY")
	env.AwsRegion = os.Getenv("AWS_REGION")
	env.AwsBucketName = os.Getenv("AWS_BUCKET_NAME")
	if env.AwsBucketName == "" {
		log.Fatalf("AWS_BUCKET_NAME is not set in the environment variables")
	}

	env.SmtpHost = os.Getenv("SMTP_HOST")
	env.SmtpPort, err = strconv.Atoi(os.Getenv("SMTP_PORT"))
	if err != nil {
		log.Fatal("Fail to parse SMTP_PORT")
	}
	env.SmtpUsername = os.Getenv("SMTP_USERNAME")
	env.SmtpEmail = os.Getenv("SMTP_EMAIL")
	env.SmtpPassword = os.Getenv("SMTP_PASSWORD")

	Env = env
}
