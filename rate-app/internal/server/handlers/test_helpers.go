package handlers

import (
	"btc-test-task/internal/common/configuration/config"
	"btc-test-task/internal/common/models"
	"btc-test-task/internal/repository"
)

var conf config.Config

type templatesImplStub struct{}

func (template *templatesImplStub) CurrencyRate(rate models.Rate) string {
	return ""
}

type rateProviderStub struct {
	RateError error
	Rate      models.Rate
}

func (provider *rateProviderStub) GetCurrentRate(currency *models.Currency) (models.Rate, error) {
	return provider.Rate, provider.RateError
}

var (
	GoodEmail     = models.Email{Value: "goodEmail@genesis.com"}
	ExistingEmail = models.Email{Value: "existing_email@gmail.com"}
	BadEmail      = models.Email{Value: "badEmailadslkf@"}
)

type emailsStorageStub struct {
	allEmails map[models.Email]struct{}
}

func (storage *emailsStorageStub) AddEmail(email *models.Email) error {
	if *email == ExistingEmail {
		return repository.ErrEmailAlreadyExists
	} else if *email == BadEmail {
		return repository.ErrInvalidEmailAddress
	}
	return nil
}

func (storage *emailsStorageStub) GetAllEmails() map[models.Email]struct{} {
	return storage.allEmails
}

func (storage *emailsStorageStub) RemoveEmail(email *models.Email) error {
	if *email != ExistingEmail {
		return repository.ErrEmailDoesNotExists
	}
	return nil
}

func (storage *emailsStorageStub) Close() {}

type ServicesStub struct {
	RateProvider     RateProvider
	EmailSender      EmailSender
	EmailsRepository EmailsRepository
	Templates        Templates
}

func (services *ServicesStub) GetEmailSenderService() EmailSender {
	return services.EmailSender
}
func (services *ServicesStub) GetEmailsRepositoryService() EmailsRepository {
	return services.EmailsRepository
}
func (services *ServicesStub) GetRateProviderService() RateProvider {
	return services.RateProvider
}
func (services *ServicesStub) GetTemplatesService() Templates {
	return services.Templates
}

type LoggerWriterStub struct{}

func (_ *LoggerWriterStub) Write(p []byte) (n int, err error) { return 0, nil }
func (_ *LoggerWriterStub) Close()                            {}
