package emailSender

import "btc-test-task/internal/helpers/models"

type EmailSender interface {
	BroadcastEmails(recipients map[models.Email]struct{}, body string)
	SendEmail(recipient models.Email, body string) error
}
