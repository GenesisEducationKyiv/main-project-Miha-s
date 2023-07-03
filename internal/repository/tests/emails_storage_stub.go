package emailsStorageTest

import (
	"btc-test-task/internal/models"
	"btc-test-task/internal/repository"
)

var (
	GoodEmail     = models.Email{Value: "goodEmail@genesis.com"}
	ExistingEmail = models.Email{Value: "existing_email@gmail.com"}
	BadEmail      = models.Email{Value: "badEmailadslkf@"}
)

type EmailsStorageStub struct {
	allEmails map[models.Email]struct{}
}

func (storage *EmailsStorageStub) AddEmail(email *models.Email) error {
	if *email == ExistingEmail {
		return repository.ErrEmailAlreadyExists
	} else if *email == BadEmail {
		return repository.ErrInvalidEmailAddress
	}
	return nil
}

func (storage *EmailsStorageStub) GetAllEmails() map[models.Email]struct{} {
	return storage.allEmails
}

func (storage *EmailsStorageStub) RemoveEmail(email *models.Email) error {
	if *email != ExistingEmail {
		return repository.ErrEmailDoesNotExists
	}
	return nil
}

func (storage *EmailsStorageStub) Close() {

}
