package emails_storage

import (
	"btc-test-task/internal/config"
	"btc-test-task/internal/logger"
	"encoding/json"
	"errors"
	"os"
)

type EmailsStorage struct {
	emails            map[string]struct{}
	storage_file_path string
	storage_name      string "emails_storage.json"
}

func (storage *EmailsStorage) Init(conf *config.Config) error {
	storage.emails = make(map[string]struct{})
	storage.storage_file_path = conf.EmailStoragePath + "/" + storage.storage_name
	if _, err := os.Stat(storage.storage_file_path); err == nil {
		return storage.openExistingStorage()
	} else {
		if err != os.ErrNotExist {
			return err
		}
	}
	return nil
}

func get_array_from_set(set *map[string]struct{}) []string {
	result := make([]string, 0)

	for key := range *set {
		result = append(result, key)
	}

	return result
}

func (storage *EmailsStorage) Close() {
	logger.LogInfo("Closing file storage")
	storage_file, err := os.Create(storage.storage_file_path)
	if err != nil {
		logger.LogError(err)
		return
	}
	json_map := make(map[string][]string)
	json_map["emails"] = get_array_from_set(&storage.emails)

	json_data, err := json.Marshal(json_map)
	if err != nil {
		logger.LogError(err)
		return
	}

	storage_file.Write(json_data)
	storage_file.Close()
}

func (storage *EmailsStorage) AddEmail(email string) error {
	if _, err := storage.emails[email]; err {
		return errors.New("email alredy exists")
	}
	storage.emails[email] = struct{}{}
	return nil
}

func (storage *EmailsStorage) GetAllEmails() *map[string]struct{} {
	return &storage.emails
}

func (storage *EmailsStorage) openExistingStorage() error {
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
