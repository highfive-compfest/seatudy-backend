package main

import (
	"github.com/highfive-compfest/seatudy-backend/internal/domain/forum"
	"log"
	"os"

	"github.com/highfive-compfest/seatudy-backend/internal/domain/courseenroll"
	"github.com/highfive-compfest/seatudy-backend/internal/domain/submission"

	"github.com/highfive-compfest/seatudy-backend/internal/domain/assignment"
	"github.com/highfive-compfest/seatudy-backend/internal/domain/attachment"
	"github.com/highfive-compfest/seatudy-backend/internal/domain/material"
	"github.com/highfive-compfest/seatudy-backend/internal/domain/review"
	"github.com/highfive-compfest/seatudy-backend/internal/domain/wallet"
	"github.com/highfive-compfest/seatudy-backend/internal/schema"

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
		&schema.Wallet{},
		&schema.MidtransTransaction{},
		&schema.User{},
		&schema.Course{},
		&schema.Material{},
		&schema.Assignment{},
		&schema.Submission{},
		&schema.Attachment{},
		&schema.Review{},
		&schema.CourseEnroll{},
		&schema.ForumDiscussion{},
		&schema.ForumReply{},
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

	courseEnrollRepo := courseenroll.NewRepository(db)
	courseEnrollUseCase := courseenroll.NewUseCase(courseEnrollRepo)

	// Course
	courseRepo := course.NewRepository(db)
	courseUseCase := course.NewUseCase(courseRepo, walletRepo, *courseEnrollUseCase)
	course.NewRestController(engine, courseUseCase, walletUseCase)

	// Attachment
	attachmentRepo := attachment.NewRepository(db)
	attachmentUseCase := attachment.NewUseCase(attachmentRepo)
	attachment.NewRestController(engine, attachmentUseCase)

	// Assignment
	assignmentRepo := assignment.NewRepository(db)
	assignmentUseCase := assignment.NewUseCase(assignmentRepo, attachmentUseCase)
	assignment.NewRestController(engine, assignmentUseCase, courseUseCase)

	// Submission
	submissionRepo := submission.NewRepository(db)
	submissionUseCase := submission.NewUseCase(submissionRepo, assignmentRepo, *attachmentUseCase, courseRepo, courseEnrollRepo)
	submission.NewRestController(engine, submissionUseCase)
	//Material
	materialRepo := material.NewRepository(db)
	materialUsecase := material.NewUseCase(materialRepo, attachmentUseCase)
	material.NewRestController(engine, materialUsecase, courseUseCase)

	// Review
	reviewRepo := review.NewRepository(db)
	reviewUseCase := review.NewUseCase(reviewRepo, courseRepo, courseEnrollUseCase)
	review.NewRestController(engine, reviewUseCase)

	// Forum
	forumRepo := forum.NewRepository(db)
	forumUseCase := forum.NewUseCase(forumRepo, courseEnrollUseCase)
	forum.NewRestController(engine, forumUseCase)

	if err := engine.Run(":" + config.Env.ApiPort); err != nil {
		log.Fatalln(err)
	}
}
