package config

import (
	"github.com/midtrans/midtrans-go"
	"log"
	"os"
	"strconv"
	"time"
)

type environmentVariables struct {
	ENV         string
	FrontendUrl string
	ApiPort     string

	PostgresHost     string
	PostgresPort     string
	PostgresUser     string
	PostgresPassword string
	PostgresDbName   string

	RedisHost     string
	RedisPort     string
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

	MidtransServerKey   string
	MidtransEnvironment midtrans.EnvironmentType
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
	env.ApiPort = os.Getenv("API_PORT")

	env.PostgresHost = os.Getenv("POSTGRES_HOST")
	env.PostgresPort = os.Getenv("POSTGRES_PORT")
	env.PostgresUser = os.Getenv("POSTGRES_USER")
	env.PostgresPassword = os.Getenv("POSTGRES_PASSWORD")
	env.PostgresDbName = os.Getenv("POSTGRES_DB")

	env.RedisHost = os.Getenv("REDIS_HOST")
	env.RedisPort = os.Getenv("REDIS_PORT")
	env.RedisPassword = os.Getenv("REDIS_PASSWORD")
	env.RedisDatabase, err = strconv.Atoi(os.Getenv("REDIS_DATABASE"))
	if err != nil && env.ENV != "test" {
		log.Fatal("Fail to parse REDIS_DATABASE")
	}

	env.JwtAccessSecret = []byte(os.Getenv("JWT_ACCESS_SECRET"))
	env.JwtAccessDuration, err = time.ParseDuration(os.Getenv("JWT_ACCESS_DURATION"))
	if err != nil && env.ENV != "test" {
		log.Fatal("Fail to parse JWT_ACCESS_DURATION")
	}

	env.JwtRefreshSecret = []byte(os.Getenv("JWT_REFRESH_SECRET"))
	env.JwtRefreshDuration, err = time.ParseDuration(os.Getenv("JWT_REFRESH_DURATION"))
	if err != nil && env.ENV != "test" {
		log.Fatal("Fail to parse JWT_REFRESH_DURATION")
	}

	env.AwsAccessId = os.Getenv("AWS_ACCESS_KEY_ID")
	env.AmsSecretAccessId = os.Getenv("AWS_SECRET_ACCESS_KEY")
	env.AwsRegion = os.Getenv("AWS_REGION")
	env.AwsBucketName = os.Getenv("AWS_BUCKET_NAME")
	if env.AwsBucketName == "" && env.ENV != "test" {
		log.Fatalf("AWS_BUCKET_NAME is not set in the environment variables")
	}

	env.SmtpHost = os.Getenv("SMTP_HOST")
	env.SmtpPort, err = strconv.Atoi(os.Getenv("SMTP_PORT"))
	if err != nil && env.ENV != "test" {
		log.Fatal("Fail to parse SMTP_PORT")
	}
	env.SmtpUsername = os.Getenv("SMTP_USERNAME")
	env.SmtpEmail = os.Getenv("SMTP_EMAIL")
	env.SmtpPassword = os.Getenv("SMTP_PASSWORD")

	env.MidtransServerKey = os.Getenv("MIDTRANS_SERVER_KEY")
	env.MidtransEnvironment = midtrans.Sandbox
	//if env.ENV == "production" {
	//	env.MidtransEnvironment = midtrans.Production
	//}

	Env = env
}
