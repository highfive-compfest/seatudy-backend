package main

import (
	"log"
	"os"

	"github.com/highfive-compfest/seatudy-backend/internal/config"
	"github.com/highfive-compfest/seatudy-backend/internal/domain/attachment"
	"github.com/highfive-compfest/seatudy-backend/internal/domain/auth"
	"github.com/highfive-compfest/seatudy-backend/internal/domain/course"
	"github.com/highfive-compfest/seatudy-backend/internal/domain/material"
	"github.com/highfive-compfest/seatudy-backend/internal/domain/user"
	"github.com/highfive-compfest/seatudy-backend/internal/middleware"
	"github.com/highfive-compfest/seatudy-backend/internal/schema"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	apiEnv := os.Getenv("ENV")
	if err != nil && apiEnv == "" {
		log.Println("fail to load env", err)
	}
	config.LoadEnv()

	db := config.NewPostgresql(
		&schema.User{},
		&schema.Course{},
		&schema.Material{},
		&schema.Attachment{},

	)
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

	// Course
	courseRepo := course.NewRepository(db)
	courseUseCase := course.NewUseCase(courseRepo)
	course.NewRestController(engine,courseUseCase)

	// Attachment
	attachmentRepo := attachment.NewRepository(db)
	attachmentUseCase := attachment.NewUseCase(attachmentRepo)
	attachment.NewRestController(engine,attachmentUseCase)

	//Material
	materialRepo := material.NewRepository(db)
	materialUsecase := material.NewUseCase(materialRepo,attachmentUseCase)
	material.NewRestController(engine, materialUsecase,courseUseCase)

	if err := engine.Run(":" + config.Env.Port); err != nil {
		log.Fatalln(err)
	}
}
