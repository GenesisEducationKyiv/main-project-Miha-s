package repository

import (
	"btc-test-task/internal/configuration/config"
	"btc-test-task/internal/configuration/logger"
	"btc-test-task/internal/models"
	"encoding/json"
	"os"

	"github.com/pkg/errors"
)

type EmailValidator interface {
	ValidateEmail(email *models.Email) bool
}

type JsonEmailsStorage struct {
	emails          map[models.Email]struct{}
	storageFilePath string
	storageName     string
	storageFile     *os.File
	validator       EmailValidator
}

func NewJsonEmailsStorage(conf *config.Config, emailValidator EmailValidator) (*JsonEmailsStorage, error) {
	newJsonEmailsStorage := new(JsonEmailsStorage)
	err := newJsonEmailsStorage.init(conf)
	newJsonEmailsStorage.validator = emailValidator
	if err != nil {
		return nil, err
	}
	return newJsonEmailsStorage, nil
}

func (storage *JsonEmailsStorage) init(conf *config.Config) error {
	storage.storageName = conf.EmailStorageName
	storage.emails = make(map[models.Email]struct{})
	storage.storageFilePath = conf.EmailStoragePath + "/" + storage.storageName
	return storage.initStorageFile()
}

func (storage *JsonEmailsStorage) initStorageFile() error {
	err := error(nil)
	storage.storageFile, err = os.OpenFile(storage.storageFilePath, os.O_CREATE|os.O_RDWR, 0755)
	if err != nil {
		return errors.Wrap(err, "initStorageFile")
	}
	fileStats, err := storage.storageFile.Stat()
	if err != nil {
		return errors.Wrap(err, "initStorageFile")
	}
	if fileStats.Size() == 0 {
		return nil
	}
	err = storage.uploadFromFile()
	if err != nil {
		return errors.Wrap(err, "initStorageFile")
	}

	return nil
}

func getArrayFromSet(set map[models.Email]struct{}) []models.Email {
	result := make([]models.Email, 0)

	for key := range set {
		result = append(result, key)
	}

	return result
}

func (storage *JsonEmailsStorage) Close() {
	logger.Log.Info("Closing file storage")
	err := storage.sync()
	if err != nil {
		logger.Log.Error("Was not able to sync before closing file")
	}
	err = storage.storageFile.Close()
	if err != nil {
		logger.Log.Error("Was not able to close the file")
	}
}

func (storage *JsonEmailsStorage) emailExists(email *models.Email) bool {
	_, ok := storage.emails[*email]
	return ok
}

func (storage *JsonEmailsStorage) AddEmail(email *models.Email) error {
	if !storage.validator.ValidateEmail(email) {
		return ErrInvalidEmailAddress
	}

	if storage.emailExists(email) {
		return ErrEmailAlreadyExists
	}
	storage.emails[*email] = struct{}{}
	return storage.sync()
}

func (storage *JsonEmailsStorage) GetAllEmails() map[models.Email]struct{} {
	return storage.emails
}

func (storage *JsonEmailsStorage) RemoveEmail(email *models.Email) error {
	if !storage.emailExists(email) {
		return ErrEmailDoesNotExists
	}
	delete(storage.emails, *email)
	return storage.sync()
}

func (storage *JsonEmailsStorage) uploadFromFile() error {
	data, err := os.ReadFile(storage.storageFilePath)
	if err != nil {
		return errors.Wrap(err, "uploadFromFile")
	}
	var jsonArray map[string][]models.Email
	err = json.Unmarshal(data, &jsonArray)
	if err != nil {
		return errors.Wrap(err, "uploadFromFile")
	}

	for _, email := range jsonArray["emails"] {
		storage.emails[email] = struct{}{}
	}

	return nil
}

func (storage *JsonEmailsStorage) sync() error {
	jsonMap := make(map[string][]models.Email)
	jsonMap["emails"] = getArrayFromSet(storage.emails)

	jsonData, err := json.Marshal(jsonMap)
	if err != nil {
		logger.Log.Error(err)
		return ErrFailedSyncStorage
	}
	err = storage.storageFile.Truncate(0)
	if err != nil {
		logger.Log.Error(err)
		return ErrFailedSyncStorage
	}
	_, err = storage.storageFile.Seek(0, 0)
	if err != nil {
		logger.Log.Error(err)
		return ErrFailedSyncStorage
	}
	_, err = storage.storageFile.Write(jsonData)
	if err != nil {
		logger.Log.Error("Was not able to save to storage")
		return ErrFailedSyncStorage
	}
	return nil
}
