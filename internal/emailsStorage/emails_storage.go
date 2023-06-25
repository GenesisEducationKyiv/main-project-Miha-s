package emailsStorage

import "btc-test-task/internal/config"

type EmailsStorage interface {
	Init(conf *config.Config) error
	AddEmail(email string) error
	GetAllEmails() *map[string]struct{}
	ValidateEmail(email string) bool
	Close()
}
