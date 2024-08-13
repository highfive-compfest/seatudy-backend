package config

import (
	"gopkg.in/gomail.v2"
)

type IMailer interface {
	DialAndSend(m ...*gomail.Message) error
}

func NewMailDialer() IMailer {
	return gomail.NewDialer(
		Env.SmtpHost,
		Env.SmtpPort,
		Env.SmtpUsername,
		Env.SmtpPassword,
	)
}
