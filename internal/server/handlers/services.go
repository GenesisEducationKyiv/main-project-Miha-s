package handlers

import (
	models2 "btc-test-task/internal/common/models"
)

type EmailSender interface {
	BroadcastEmails(recipients map[models2.Email]struct{}, body string)
	SendEmail(recipient models2.Email, body string) error
}

type EmailsRepository interface {
	AddEmail(email *models2.Email) error
	GetAllEmails() map[models2.Email]struct{}
	RemoveEmail(email *models2.Email) error
	Close()
}

type RateProvider interface {
	GetCurrentRate(currency *models2.Currency) (models2.Rate, error)
}

type Templates interface {
	CurrencyRate(rate models2.Rate) string
}

type Services interface {
	GetEmailSenderService() EmailSender
	GetEmailsRepositoryService() EmailsRepository
	GetRateProviderService() RateProvider
	GetTemplatesService() Templates
}
