package emails_storage

import (
	"btc-test-task/internal/config"
	"btc-test-task/internal/logger"
	"encoding/json"
	"errors"
	"os"
)

type EmailsStorageImpl struct {
	emails            map[string]struct{}
	storage_file_path string
	storage_name      string
}

func fileExists(filepath string) bool {
	info, err := os.Stat(filepath)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

func (storage *EmailsStorageImpl) Init(conf *config.Config) error {
	storage.storage_name = "emails_storage.json"
	storage.emails = make(map[string]struct{})
	storage.storage_file_path = conf.EmailStoragePath + "/" + storage.storage_name
	if !fileExists(storage.storage_file_path) {
		return nil
	}
	return storage.openExistingStorage()
}

func getArrayFromSet(set *map[string]struct{}) []string {
	result := make([]string, 0)

	for key := range *set {
		result = append(result, key)
	}

	return result
}

func (storage *EmailsStorageImpl) Close() {
	logger.LogInfo("Closing file storage")
	storage_file, err := os.Create(storage.storage_file_path)
	if err != nil {
		logger.LogError(err)
		return
	}
	json_map := make(map[string][]string)
	json_map["emails"] = getArrayFromSet(&storage.emails)

	json_data, err := json.Marshal(json_map)
	if err != nil {
		logger.LogError(err)
		return
	}

	_, err = storage_file.Write(json_data)
	if err != nil {
		logger.LogErrorStr("Was not able to save storage")
	}
	storage_file.Close()
}

func (storage *EmailsStorageImpl) AddEmail(email string) error {
	if _, ok := storage.emails[email]; ok {
		return errors.New("email alredy exists")
	}
	storage.emails[email] = struct{}{}
	return nil
}

func (storage *EmailsStorageImpl) GetAllEmails() *map[string]struct{} {
	return &storage.emails
}

func (storage *EmailsStorageImpl) openExistingStorage() error {
	data, err := os.ReadFile(storage.storage_file_path)
	if err != nil {
		return err
	}
	var json_map map[string]interface{}
	err = json.Unmarshal(data, &json_map)
	if err != nil {
		return err
	}
	json_array := json_map["emails"].([]interface{})

	for _, email := range json_array {
		storage.emails[email.(string)] = struct{}{}
	}

	return nil
}
