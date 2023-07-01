package handlers

import "btc-test-task/internal/helpers/models"

type EmailsRepository interface {
	AddEmail(email *models.Email) error
	GetAllEmails() map[models.Email]struct{}
	RemoveEmail(email *models.Email) error
	Close()
}
