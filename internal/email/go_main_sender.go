package email

import (
	"btc-test-task/internal/configuration/config"
	"btc-test-task/internal/configuration/logger"
	"btc-test-task/internal/models"
	"crypto/tls"

	"gopkg.in/gomail.v2"
)

type GoMailSender struct {
	email   string
	subject string
	dialer  gomail.Dialer
}

func NewGoMailSender(conf *config.Config) *GoMailSender {
	newEmailSender := new(GoMailSender)
	newEmailSender.init(conf)

	return newEmailSender
}

func (sender *GoMailSender) init(conf *config.Config) {
	sender.dialer = *gomail.NewDialer(conf.EmailServiceUrl, conf.EmailServicePort,
		conf.EmailToSendFrom, conf.EmailToSendFromPassword)

	sender.email = conf.EmailToSendFrom
	sender.subject = conf.EmailSubject
	sender.dialer.TLSConfig = &tls.Config{InsecureSkipVerify: true}

	sender.subject = conf.EmailSubject
}

func (sender *GoMailSender) SendEmail(recipient models.Email, body string) error {
	message := gomail.NewMessage()
	message.SetHeader("From", sender.email)
	message.SetHeader("To", recipient.Value)
	message.SetHeader("Subject", sender.subject)
	message.SetBody("text/plain", body)

	if err := sender.dialer.DialAndSend(message); err != nil {
		logger.Log.Error(err)
		return ErrFailedToSendEmail
	}
	return nil
}

func (sender *GoMailSender) BroadcastEmails(recipients map[models.Email]struct{}, body string) {
	for email := range recipients {
		err := sender.SendEmail(email, body)
		if err != nil {
			logger.Log.Warn(err)
		}
	}
}
