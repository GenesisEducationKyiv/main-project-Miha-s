package types

import (
	"btc-test-task/internal/course_accessors"
	"btc-test-task/internal/email_sender"
	"btc-test-task/internal/emails_storage"
)

type Services struct {
	CourseAccessor course_accessors.CoinApI
	EmailSender    email_sender.EmailSender
	EmailStorage   emails_storage.EmailsStorage
}
