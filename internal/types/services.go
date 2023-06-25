package types

import (
	"btc-test-task/internal/emailSender"
	"btc-test-task/internal/emailsStorage"
	"btc-test-task/internal/rateAccessors"
	"btc-test-task/internal/templates"
)

type Services struct {
	RateAccessor rateAccessors.RateAccessor
	EmailSender  emailSender.EmailSender
	EmailStorage emailsStorage.EmailsStorage
	Templates    templates.Templates
}
