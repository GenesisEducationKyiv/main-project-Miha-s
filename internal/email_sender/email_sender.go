package email_sender

import (
	"btc-test-task/internal/config"
	"btc-test-task/internal/logger"

	"crypto/tls"

	gomail "gopkg.in/gomail.v2"
)

type EmailSender struct {
	email   string
	subject string
	dialer  gomail.Dialer
}

var Sender EmailSender

func Init(conf *config.Config) {
	Sender.dialer = *gomail.NewDialer(conf.EmailServiceUrl, conf.EmailServicePort,
		conf.EmailToSendFrom, conf.EmailToSendFromPassword)

	Sender.email = conf.EmailToSendFrom
	Sender.subject = conf.EmailSubject
	Sender.dialer.TLSConfig = &tls.Config{InsecureSkipVerify: true}
	
	Sender.subject = conf.EmailSubject
}

func SendEmail(recipient, body string) error {
	message := gomail.NewMessage()
	message.SetHeader("From", Sender.email)
	message.SetHeader("To", recipient)
	message.SetHeader("Subject", Sender.subject)
	message.SetBody("text/plain", body)

	if err := Sender.dialer.DialAndSend(message); err != nil {
		logger.LogError(err)
		return err
	}
	return nil
}
