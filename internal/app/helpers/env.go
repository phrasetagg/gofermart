package helpers

import "os"

func GetEnv(envName string, defaultValue string) string {
	value := os.Getenv(envName)

	if value != "" {
		return value
	}

	return defaultValue
}
