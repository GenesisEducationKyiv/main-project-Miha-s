package validators

import "btc-test-task/internal/helpers/models"

type EmailValidator interface {
	ValidateEmail(email *models.Email) bool
}
