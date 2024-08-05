package mailer

import (
	"github.com/highfive-compfest/seatudy-backend/internal/config"
	"gopkg.in/gomail.v2"
)

func NewMail() *gomail.Message {
	mail := gomail.NewMessage()
	mail.SetHeader("From", "Seatudy <"+config.Env.SmtpEmail+">")
	return mail
}
