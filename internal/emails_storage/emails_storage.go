package emails_storage

import (
	"btc-test-task/internal/config"
	"btc-test-task/internal/logger"
	"encoding/json"
	"errors"
	"os"
)

var (
	emails            = make(map[string]struct{})
	storage_file_path string
)

const storage_name = "emails_storage.json"

func Init(conf *config.Config) error {
	storage_file_path = conf.EmailStoragePath + "/" + storage_name
	if _, err := os.Stat(storage_file_path); err == nil {
		return openExistingStorage(storage_file_path)
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

func Close() {
	logger.LogInfo("Closing file storage")
	storage_file, err := os.Create(storage_file_path)
	if err != nil {
		logger.LogError(err)
		return
	}
	json_map := make(map[string][]string)
	json_map["emails"] = get_array_from_set(&emails)

	json_data, err := json.Marshal(json_map)
	if err != nil {
		logger.LogError(err)
		return
	}

	storage_file.Write(json_data)
	storage_file.Close()
}

func AddEmail(email string) error {
	if _, err := emails[email]; err {
		return errors.New("email alredy exists")
	}
	emails[email] = struct{}{}
	return nil
}

func GetAllEmails() *map[string]struct{} {
	return &emails
}

func openExistingStorage(path string) error {
	data, err := os.ReadFile(storage_file_path)
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
		emails[email.(string)] = struct{}{}
	}

	return nil
}
