package emailsStorageTest

import (
	"btc-test-task/internal/emailsStorage"
	"btc-test-task/internal/helpers/config"
	"btc-test-task/internal/helpers/logger"
	"fmt"
	"os"
	"testing"

	"github.com/pkg/errors"
)

var conf config.Config
var (
	email1   string
	email2   string
	email3   string
	badEmail string
)

func globalSetup() error {
	err := conf.LoadFromENV("../../../.env.test")
	if err != nil {
		return err
	}
	err = os.MkdirAll(conf.EmailStoragePath, os.ModePerm)
	if err != nil {
		return err
	}
	err = logger.Init(&conf)
	if err != nil {
		return err
	}
	email1 = "email1@gmail.com"
	email2 = "email2@genesis.com"
	email3 = "email3@ukr.net"
	badEmail = "veryverybademail@"
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

func addEmailTest(err error, email string, t *testing.T) {
	if err != nil {
		t.Fatalf("failed to add email %v, because of error %v", email, err)
	}
}

func missingEmailTest(exists bool, email string, t *testing.T) {
	if !exists {
		t.Fatalf("missing email %v", email)
	}
}

func TestCreateEmailFile(t *testing.T) {
	_, err := emailsStorage.NewJsonEmailsStorage(&conf)
	storageInitializationTest(err, t)
	tearDown(t)
}

func TestAddEmail(t *testing.T) {
	storage, err := emailsStorage.NewJsonEmailsStorage(&conf)
	storageInitializationTest(err, t)
	err = storage.AddEmail(email1)
	addEmailTest(err, email1, t)
	err = storage.AddEmail(email2)
	addEmailTest(err, email2, t)

	allEmails := storage.GetAllEmails()
	_, ok := allEmails[email1]
	missingEmailTest(ok, email1, t)

	_, ok = allEmails[email2]
	missingEmailTest(ok, email2, t)
	tearDown(t)
}

func TestErrorEmailExists(t *testing.T) {
	storage, err := emailsStorage.NewJsonEmailsStorage(&conf)
	storageInitializationTest(err, t)
	err = storage.AddEmail(email1)
	addEmailTest(err, email1, t)
	err = storage.AddEmail(email1)
	if !errors.Is(err, emailsStorage.ErrEmailAlreadyExists) {
		t.Errorf("incorrect error when adding same email %v", err)
	}
	tearDown(t)
}

func TestValidEmail(t *testing.T) {
	storage, err := emailsStorage.NewJsonEmailsStorage(&conf)
	storageInitializationTest(err, t)
	ok := storage.ValidateEmail(email3)
	if !ok {
		t.Errorf("failed to vaildate valid email %v", email1)
	}
	ok = storage.ValidateEmail(badEmail)
	if ok {
		t.Errorf("failed to validate invalid email %v", badEmail)
	}
	tearDown(t)
}

func TestLoadFromPersistence(t *testing.T) {
	storage, err := emailsStorage.NewJsonEmailsStorage(&conf)
	storageInitializationTest(err, t)
	err = storage.AddEmail(email1)
	addEmailTest(err, email1, t)
	err = storage.AddEmail(email2)
	addEmailTest(err, email2, t)
	storage.Close()

	newStorage, err := emailsStorage.NewJsonEmailsStorage(&conf)
	storageInitializationTest(err, t)
	allEmails := newStorage.GetAllEmails()
	_, ok := allEmails[email1]
	missingEmailTest(ok, email1, t)

	_, ok = allEmails[email2]
	missingEmailTest(ok, email2, t)
	tearDown(t)
}
