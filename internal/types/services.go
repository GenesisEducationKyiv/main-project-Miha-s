package types

import (
	"btc-test-task/internal/email_sender"
	"btc-test-task/internal/emails_storage"
	"btc-test-task/internal/rate_accessors"
	"btc-test-task/internal/templates"
)

type Services struct {
	RateAccessor rate_accessors.RateAccessor
	EmailSender  email_sender.EmailSender
	EmailStorage emails_storage.EmailsStorage
	Templates    templates.Templates
}
