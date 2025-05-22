package config

import (
	"encoding/json"
	"io"
	"os"
)

// Parse reads and unmarshals config file
func Parse(file string, model any) error {
	// Open json file.
	jsonFile, err := os.Open(file)

	if err != nil {
		return err
	}

	defer jsonFile.Close()
	fileBytes, err := io.ReadAll(jsonFile)

	if err != nil {
		return err
	}

	// Attempt to decode json.
	return json.Unmarshal(fileBytes, model)
}
