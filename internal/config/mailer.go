package config

import "gopkg.in/gomail.v2"

func NewMailDialer() *gomail.Dialer {
	return gomail.NewDialer(
		Env.SmtpHost,
		Env.SmtpPort,
		Env.SmtpUsername,
		Env.SmtpPassword,
	)
}
