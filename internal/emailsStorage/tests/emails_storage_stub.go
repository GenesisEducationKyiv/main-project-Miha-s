package emailsStorageTest

import "btc-test-task/internal/emailsStorage"

const (
	GoodEmail     = "goodEmail@genesis.com"
	ExistingEmail = "existing_email@gmail.com"
	BadEmail      = "badEmailadslkf@"
)

type EmailsStorageStub struct {
	allEmails map[string]struct{}
}

func (storage *EmailsStorageStub) ValidateEmail(email string) bool {
	if email == BadEmail {
		return false
	} else {
		return true
	}
}

func (storage *EmailsStorageStub) AddEmail(email string) error {
	if email == ExistingEmail {
		return emailsStorage.ErrEmailAlreadyExists
	} else {
		return nil
	}
}

func (storage *EmailsStorageStub) GetAllEmails() map[string]struct{} {
	return storage.allEmails
}

func (storage *EmailsStorageStub) Close() {

}
