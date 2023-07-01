package types

import (
	"btc-test-task/internal/server/handlers"
)

type Services struct {
	RateProvider     handlers.RateProvider
	EmailSender      handlers.EmailSender
	EmailsRepository handlers.EmailsRepository
	Templates        handlers.Templates
}

func (services *Services) GetEmailSenderService() handlers.EmailSender {
	return services.EmailSender
}
func (services *Services) GetEmailsRepositoryService() handlers.EmailsRepository {
	return services.EmailsRepository
}
func (services *Services) GetRateProviderService() handlers.RateProvider {
	return services.RateProvider
}
func (services *Services) GetTemplatesService() handlers.Templates {
	return services.Templates
}
