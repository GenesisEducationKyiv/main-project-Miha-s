package handlers

import "btc-test-task/internal/helpers/models"

type EmailSender interface {
	BroadcastEmails(recipients map[models.Email]struct{}, body string)
	SendEmail(recipient models.Email, body string) error
}

type EmailsRepository interface {
	AddEmail(email *models.Email) error
	GetAllEmails() map[models.Email]struct{}
	RemoveEmail(email *models.Email) error
	Close()
}

type RateProvider interface {
	GetCurrentRate(currency *models.Currency) (models.Rate, error)
}

type Templates interface {
	CurrencyRate(rate models.Rate) string
}

type Services interface {
	GetEmailSenderService() EmailSender
	GetEmailsRepositoryService() EmailsRepository
	GetRateProviderService() RateProvider
	GetTemplatesService() Templates
}
