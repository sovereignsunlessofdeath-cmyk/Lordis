package database

import (
	"encoding/json"
	"os"
	"sync"

	"lordis/internal/models"
)

var (
	mutex  sync.Mutex
	dbPath = "internal/database/data.json"
)

// LoadData safely reads the JSON file into structs
func LoadData() (models.AppData, error) {
	mutex.Lock()
	defer mutex.Unlock()

	var data models.AppData
	file, err := os.ReadFile(dbPath)
	if err != nil {
		// If file doesn't exist yet, return empty structure cleanly
		if os.IsNotExist(err) {
			return data, nil
		}
		return data, err
	}

	err = json.Unmarshal(file, &data)
	return data, err
}

// SaveData safely marshals structs back down to the JSON file
func SaveData(data models.AppData) error {
	mutex.Lock()
	defer mutex.Unlock()

	file, err := json.MarshalIndent(data, "", "    ")
	if err != nil {
		return err
	}

	return os.WriteFile(dbPath, file, 0644)
}