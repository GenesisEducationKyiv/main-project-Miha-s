package emails_storage

import "btc-test-task/internal/config"

type EmailsStorage interface {
	Init(conf *config.Config) error
	AddEmail(email string) error
	GetAllEmails() *map[string]struct{}
	Close()
}
