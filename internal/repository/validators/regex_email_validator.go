package validators

import (
	"btc-test-task/internal/models"
	"net/mail"
)

type RegexEmailValidator struct {
}

func (validator *RegexEmailValidator) ValidateEmail(email *models.Email) bool {
	_, err := mail.ParseAddress(email.Value)
	return err == nil
}
