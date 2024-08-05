package main

import (
	"github.com/highfive-compfest/seatudy-backend/internal/config"
	"github.com/highfive-compfest/seatudy-backend/internal/domain/auth"
	"github.com/highfive-compfest/seatudy-backend/internal/domain/user"
	"github.com/highfive-compfest/seatudy-backend/internal/middleware"
	"github.com/joho/godotenv"
	"log"
	"os"

)

func main() {
	err := godotenv.Load()
	apiEnv := os.Getenv("ENV")
	if err != nil && apiEnv == "" {
		log.Println("fail to load env", err)
	}
	config.LoadEnv()

	db := config.NewPostgresql()
	rds := config.NewRedis()

	mailDialer := config.NewMailDialer()

	engine := config.NewGin()
	engine.Use(middleware.CORS())

	config.InitializeS3()

	// User
	userRepo := user.NewRepository(db)
	userUseCase := user.NewUseCase(userRepo)
	user.NewRestController(engine, userUseCase)

	// Auth
	authRepo := auth.NewRepository(rds)
	authUseCase := auth.NewUseCase(authRepo, userRepo, mailDialer)
	auth.NewRestController(engine, authUseCase)

	if err := engine.Run(":" + config.Env.Port); err != nil {
		log.Fatalln(err)
	}
}