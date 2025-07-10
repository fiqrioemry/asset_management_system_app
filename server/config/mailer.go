package config

import (
	"fmt"

	gomail "gopkg.in/gomail.v2"
)

var MailDialer *gomail.Dialer

func InitMailer() {
	MailDialer = gomail.NewDialer(AppConfig.SMTPHost, AppConfig.SMTPPort, AppConfig.SMTPEmail, AppConfig.SMTPPassword)
	MailDialer.TLSConfig = nil

	fmt.Println("✅ Mailer configured")
}
