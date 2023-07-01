package types

import (
	"btc-test-task/internal/emailSender"
	"btc-test-task/internal/emailsRepository"
	"btc-test-task/internal/helpers/templates"
	"btc-test-task/internal/rateProviders"
)

type Services struct {
	RateProvider rateProviders.RateProvider
	EmailSender  emailSender.EmailSender
	EmailStorage emailsRepository.EmailsRepository
	Templates    templates.Templates
}
