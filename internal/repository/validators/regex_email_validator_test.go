package validators

import (
	"btc-test-task/internal/common/models"
	"testing"
)

func TestRegexValidator(t *testing.T) {
	validator := RegexEmailValidator{}

	var tests = []struct {
		email models.Email
		valid bool
	}{
		{email: models.Email{Value: "someemail@gmail.com"}, valid: true},
		{email: models.Email{Value: "another_email@genesis.com"}, valid: true},
		{email: models.Email{Value: "gmail.com"}, valid: false},
		{email: models.Email{Value: "someemail@.com"}, valid: false},
		{email: models.Email{Value: ""}, valid: false},
	}

	for _, testEmail := range tests {
		t.Run(testEmail.email.Value, func(t *testing.T) {
			valid := validator.ValidateEmail(&testEmail.email)
			if valid != testEmail.valid {
				t.Errorf("expected email %v, valid %v", testEmail.email, testEmail.valid)
			}
		})
	}
}
