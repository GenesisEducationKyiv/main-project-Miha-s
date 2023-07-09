package emailsStorageTest

import (
	"btc-test-task/internal/common/configuration/config"
	"btc-test-task/internal/common/configuration/logger"
	"btc-test-task/internal/common/models"
	"btc-test-task/internal/repository"
	"fmt"
	"os"
	"testing"

	"github.com/pkg/errors"
)

type EmailValidatorStub struct {
}

func (validator *EmailValidatorStub) ValidateEmail(email *models.Email) bool {
	return true
}

var validator = &EmailValidatorStub{}
var conf config.Config
var (
	email1   = &models.Email{}
	email2   = &models.Email{}
	email3   = &models.Email{}
	badEmail = &models.Email{}
)

func globalSetup() error {
	conf.EmailStoragePath = "./tests"
	conf.EmailStorageName = "test_storage.json"
	err := os.MkdirAll(conf.EmailStoragePath, os.ModePerm)
	if err != nil {
		return err
	}
	err = logger.Init(&conf)
	if err != nil {
		return err
	}
	email1.Value = "email1@gmail.com"
	email2.Value = "email2@genesis.com"
	email3.Value = "email3@ukr.net"
	badEmail.Value = "veryverybademail@"
	return nil
}

func tearDown(t *testing.T) {
	err := os.Remove(conf.EmailStoragePath + "/" + conf.EmailStorageName)
	if err != nil {
		t.Fatalf("Failed to remove starage %v", err)
	}
}

func TestMain(m *testing.M) {
	err := globalSetup()
	if err != nil {
		panic(fmt.Sprintf("failed to setup test %v", err))
	}
	code := m.Run()
	_ = os.RemoveAll(conf.EmailStoragePath)
	os.Exit(code)
}

func storageInitializationTest(err error, t *testing.T) {
	if err != nil {
		t.Fatalf("failed to initialize storage %v", err)
	}
}

func addEmailTest(err error, email *models.Email, t *testing.T) {
	if err != nil {
		t.Fatalf("failed to add email %v, because of error %v", email, err)
	}
}

func missingEmailTest(exists bool, email *models.Email, t *testing.T) {
	if !exists {
		t.Fatalf("missing email %v", email)
	}
}

func TestCreateEmailFile(t *testing.T) {
	_, err := repository.NewJsonEmailsStorage(&conf, validator)
	storageInitializationTest(err, t)
	tearDown(t)
}

func TestAddEmail(t *testing.T) {
	storage, err := repository.NewJsonEmailsStorage(&conf, validator)
	storageInitializationTest(err, t)
	err = storage.AddEmail(email1)
	addEmailTest(err, email1, t)
	err = storage.AddEmail(email2)
	addEmailTest(err, email2, t)

	allEmails := storage.GetAllEmails()
	_, ok := allEmails[*email1]
	missingEmailTest(ok, email1, t)

	_, ok = allEmails[*email2]
	missingEmailTest(ok, email2, t)
	tearDown(t)
}

func TestErrorEmailExists(t *testing.T) {
	storage, err := repository.NewJsonEmailsStorage(&conf, validator)
	storageInitializationTest(err, t)
	err = storage.AddEmail(email1)
	addEmailTest(err, email1, t)
	err = storage.AddEmail(email1)
	if !errors.Is(err, repository.ErrEmailAlreadyExists) {
		t.Errorf("incorrect error when adding same email %v", err)
	}
	tearDown(t)
}

func TestRemoveEmail(t *testing.T) {
	storage, err := repository.NewJsonEmailsStorage(&conf, validator)
	storageInitializationTest(err, t)
	err = storage.AddEmail(email1)
	addEmailTest(err, email1, t)
	err = storage.RemoveEmail(email1)
	if err != nil {
		t.Errorf("error when removing existing email %v", err)
	}
	tearDown(t)
}

func TestErrEmailNotExists(t *testing.T) {
	storage, err := repository.NewJsonEmailsStorage(&conf, validator)
	storageInitializationTest(err, t)
	err = storage.RemoveEmail(email1)
	if !errors.Is(err, repository.ErrEmailDoesNotExists) {
		t.Errorf("removing non existing email: error expected %v, got %v",
			repository.ErrEmailDoesNotExists, err)
	}
	tearDown(t)
}

func TestLoadFromPersistence(t *testing.T) {
	storage, err := repository.NewJsonEmailsStorage(&conf, validator)
	storageInitializationTest(err, t)
	err = storage.AddEmail(email1)
	addEmailTest(err, email1, t)
	err = storage.AddEmail(email2)
	addEmailTest(err, email2, t)
	storage.Close()

	newStorage, err := repository.NewJsonEmailsStorage(&conf, validator)
	storageInitializationTest(err, t)
	allEmails := newStorage.GetAllEmails()
	_, ok := allEmails[*email1]
	missingEmailTest(ok, email1, t)

	_, ok = allEmails[*email2]
	missingEmailTest(ok, email2, t)
	tearDown(t)
}
