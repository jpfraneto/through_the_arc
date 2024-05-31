package utils

import (
	"encoding/json"
	"io/ioutil"
)

// StoreData stores data in a JSON file
func StoreData(filename string, data interface{}) error {
	file, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return err
	}

	return ioutil.WriteFile(filename, file, 0644)
}

// LoadData loads data from a JSON file
func LoadData(filename string, data interface{}) error {
	file, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}

	return json.Unmarshal(file, data)
}
