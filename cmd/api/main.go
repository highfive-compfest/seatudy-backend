package main

import (
	"github.com/highfive-compfest/seatudy-backend/internal/domain/wallet"
	"log"
	"os"

	"github.com/highfive-compfest/seatudy-backend/internal/config"
	"github.com/highfive-compfest/seatudy-backend/internal/domain/auth"
	"github.com/highfive-compfest/seatudy-backend/internal/domain/course"
	"github.com/highfive-compfest/seatudy-backend/internal/domain/user"
	"github.com/highfive-compfest/seatudy-backend/internal/middleware"

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
		&wallet.Wallet{},
		&wallet.MidtransTransaction{},
		&user.User{},
		&course.Course{},
	)
	rds := config.NewRedis()

	mailDialer := config.NewMailDialer()
	config.SetupMidtrans()

	engine := config.NewGin()
	engine.Use(middleware.CORS())

	config.InitializeS3()

	// Wallet
	walletRepo := wallet.NewRepository(db)
	walletUseCase := wallet.NewUseCase(walletRepo, nil)
	midtUseCase := wallet.NewMidtransUseCase(walletUseCase)
	walletUseCase.MidtUc = midtUseCase
	wallet.NewRestController(engine, walletUseCase, midtUseCase)

	// User
	userRepo := user.NewRepository(db, walletRepo)
	userUseCase := user.NewUseCase(userRepo)
	user.NewRestController(engine, userUseCase)

	// Auth
	authRepo := auth.NewRepository(rds)
	authUseCase := auth.NewUseCase(authRepo, userRepo, mailDialer)
	auth.NewRestController(engine, authUseCase)

	// Course
	courseRepo := course.NewRepository(db)
	courseUseCase := course.NewUseCase(courseRepo)
	course.NewRestController(engine, courseUseCase)

	if err := engine.Run(":" + config.Env.Port); err != nil {
		log.Fatalln(err)
	}
}
