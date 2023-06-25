package email_sender

import "btc-test-task/internal/config"

type EmailSender interface {
	Init(conf *config.Config) error
	BroadcastEmails(recipients *map[string]struct{}, body string)
	SendEmail(recipient, body string) error
}
