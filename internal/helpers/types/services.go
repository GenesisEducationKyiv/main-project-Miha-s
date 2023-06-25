package types

import (
	"btc-test-task/internal/emailSender"
	"btc-test-task/internal/emailsStorage"
	"btc-test-task/internal/helpers/templates"
	"btc-test-task/internal/rateAccessors"
)

type Services struct {
	RateAccessor rateAccessors.RateAccessor
	EmailSender  emailSender.EmailSender
	EmailStorage emailsStorage.EmailsStorage
	Templates    templates.Templates
}
