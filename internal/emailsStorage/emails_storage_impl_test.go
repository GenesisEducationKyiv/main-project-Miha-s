package emailsStorage

import (
	"btc-test-task/internal/helpers/config"
	"btc-test-task/internal/helpers/logger"
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
	err := conf.LoadFromENV("../../.env.test")
	if err != nil {
		return err
	}
	err = os.MkdirAll(conf.EmailStoragePath, os.ModePerm)
	if err != nil {
		return err
	}
	_ = logger.Init(&conf)
	email1 = "email1@gmail.com"
	email2 = "email2@genesis.com"
	email3 = "email3@ukr.net"
	badEmail = "veryverybademail@"
	return nil
}

func tearDown() {
	_ = os.Remove(conf.EmailStoragePath + "/" + conf.EmailStorageName)
}

func TestMain(m *testing.M) {
	err := globalSetup()
	if err != nil {
		os.Exit(2)
	}
	code := m.Run()
	_ = os.RemoveAll(conf.EmailStoragePath)
	os.Exit(code)
}

func TestCreateEmailFile(t *testing.T) {
	_, err := NewJsonEmailsStorage(&conf)
	if err != nil {
		t.Errorf("failed to create storage file, %v", err)
		return
	}
	tearDown()
}

func TestAddEmail(t *testing.T) {
	storage, _ := NewJsonEmailsStorage(&conf)
	err := storage.AddEmail(email1)
	if err != nil {
		t.Errorf("failed to add email %v", err)
		return
	}
	_ = storage.AddEmail(email2)

	allEmails := storage.GetAllEmails()
	_, ok := allEmails[email1]
	if !ok {
		t.Errorf("missing email %v", email1)
	}

	_, ok = allEmails[email2]
	if !ok {
		t.Errorf("missing email %v", email2)
	}
	tearDown()
}

func TestErrorEmailExists(t *testing.T) {
	storage, _ := NewJsonEmailsStorage(&conf)
	_ = storage.AddEmail(email1)
	err := storage.AddEmail(email1)
	if !errors.Is(err, ErrEmailAlreadyExists) {
		t.Errorf("incorrect error when adding same email %v", err)
	}
	tearDown()
}

func TestValidEmail(t *testing.T) {
	storage, _ := NewJsonEmailsStorage(&conf)
	ok := storage.ValidateEmail(email3)
	if !ok {
		t.Errorf("failed to vaildate valid email %v", email1)
	}
	ok = storage.ValidateEmail(badEmail)
	if ok {
		t.Errorf("failed to validate invalid email %v", badEmail)
	}
	tearDown()
}

func TestLoadFromPersistence(t *testing.T) {
	storage, _ := NewJsonEmailsStorage(&conf)
	_ = storage.AddEmail(email1)
	_ = storage.AddEmail(email2)
	storage.Close()

	newStorage, _ := NewJsonEmailsStorage(&conf)
	allEmails := newStorage.GetAllEmails()
	_, ok := allEmails[email1]
	if !ok {
		t.Errorf("missing email %v", email1)
	}

	_, ok = allEmails[email2]
	if !ok {
		t.Errorf("missing email %v", email2)
	}
	tearDown()
}
