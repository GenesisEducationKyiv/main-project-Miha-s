package emailsStorage

import (
	"btc-test-task/internal/config"
	"btc-test-task/internal/logger"
	"encoding/json"
	"errors"
	"net/mail"
	"os"
)

type EmailsStorageImpl struct {
	emails          map[string]struct{}
	storageFilePath string
	storageName     string
	storageFile     *os.File
}

func fileExists(filepath string) bool {
	info, err := os.Stat(filepath)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

func (storage *EmailsStorageImpl) Init(conf *config.Config) error {
	storage.storageName = "emails_storage.json"
	storage.emails = make(map[string]struct{})
	storage.storageFilePath = conf.EmailStoragePath + "/" + storage.storageName
	return storage.initStorageFile()
}

func (storage *EmailsStorageImpl) initStorageFile() error {
	err := error(nil)
	storage.storageFile, err = os.OpenFile(storage.storageFilePath, os.O_CREATE|os.O_RDWR, 0755)
	if err != nil {
		return err
	}
	fileStats, err := storage.storageFile.Stat()
	if err != nil {
		return err
	}
	if fileStats.Size() == 0 {
		return nil
	}
	return storage.uploadFromFile()
}

func getArrayFromSet(set map[string]struct{}) []string {
	result := make([]string, 0)

	for key := range set {
		result = append(result, key)
	}

	return result
}

func (storage *EmailsStorageImpl) Close() {
	logger.LogInfo("Closing file storage")
	err := storage.sync()
	if err != nil {
		logger.LogErrorStr("Was not able to sync before closing file")
	}
	err = storage.storageFile.Close()
	if err != nil {
		logger.LogErrorStr("Was not able to close the file")
	}
}

func (storage *EmailsStorageImpl) AddEmail(email string) error {
	if _, ok := storage.emails[email]; ok {
		return errors.New("email alredy exists")
	}
	storage.emails[email] = struct{}{}
	storage.sync()
	return nil
}

func (storage *EmailsStorageImpl) GetAllEmails() *map[string]struct{} {
	return &storage.emails
}

func (storage *EmailsStorageImpl) ValidateEmail(email string) bool {
	_, err := mail.ParseAddress(email)
	return err == nil
}

func (storage *EmailsStorageImpl) uploadFromFile() error {
	data, err := os.ReadFile(storage.storageFilePath)
	if err != nil {
		return err
	}
	var jsonMap map[string]interface{}
	err = json.Unmarshal(data, &jsonMap)
	if err != nil {
		return err
	}
	jsonArray := jsonMap["emails"].([]interface{})

	for _, email := range jsonArray {
		storage.emails[email.(string)] = struct{}{}
	}

	return nil
}

func (storage *EmailsStorageImpl) sync() error {
	jsonMap := make(map[string][]string)
	jsonMap["emails"] = getArrayFromSet(storage.emails)

	jsonData, err := json.Marshal(jsonMap)
	if err != nil {
		logger.LogError(err)
		return err
	}
	err = storage.storageFile.Truncate(0)
	if err != nil {
		logger.LogError(err)
		return err
	}
	_, err = storage.storageFile.Seek(0, 0)
	if err != nil {
		logger.LogError(err)
		return err
	}
	_, err = storage.storageFile.Write(jsonData)
	if err != nil {
		logger.LogErrorStr("Was not able to save to storage")
	}
	return nil
}
