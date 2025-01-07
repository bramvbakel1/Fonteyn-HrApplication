package utils

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
)

func LoadEnvironment() (map[string]string, error) {
	err := godotenv.Load()
	if err != nil {
		return nil, fmt.Errorf("could not load .env file")
	}

	vars := map[string]string{
		"AZURE_TENANT_ID":     os.Getenv("AZURE_TENANT_ID"),
		"AZURE_CLIENT_ID":     os.Getenv("AZURE_CLIENT_ID"),
		"AZURE_CLIENT_SECRET": os.Getenv("AZURE_CLIENT_SECRET"),
		"AZURE_KEYVAULT_URL":  os.Getenv("AZURE_KEYVAULT_URL"),
		"AZURE_CERT_NAME":     os.Getenv("AZURE_CERT_NAME"),
	}

	for key, value := range vars {
		if value == "" {
			return nil, fmt.Errorf("missing environment variable: %s", key)
		}
	}

	return vars, nil
}

// SaveToTempFile saves data to a temporary file with the given pattern.
func SaveToTempFile(data []byte, pattern string) (string, error) {
	// Create a temporary file with the provided pattern.
	file, err := os.CreateTemp("", pattern)
	if err != nil {
		return "", err
	}
	defer file.Close()

	// Write the data to the file.
	_, err = file.Write(data)
	if err != nil {
		return "", err
	}

	// Return the name of the temporary file.
	return file.Name(), nil
}

func RemoveTempFile(filePath string) {
	err := os.Remove(filePath)
	if err != nil {
		log.Printf("Error removing temporary file: %v", err)
	}
}
