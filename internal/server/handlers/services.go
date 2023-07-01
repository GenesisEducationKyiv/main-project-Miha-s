package handlers

type Services interface {
	GetEmailSenderService() EmailSender
	GetEmailsRepositoryService() EmailsRepository
	GetRateProviderService() RateProvider
	GetTemplatesService() Templates
}
