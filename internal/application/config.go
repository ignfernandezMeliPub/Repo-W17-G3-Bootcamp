package application

import "os"

func getEnvOrDefault(key string, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}

	return defaultValue
}

// ? Server config
var ServerAddress = getEnvOrDefault("ServerAddress", ":8080")

// ? DB Config
var DbUser = getEnvOrDefault("DbUser", "root")
var DbPassword = getEnvOrDefault("DbPassword", "")
var DbProtocol = getEnvOrDefault("DbProtocol", "tcp")
var DbAddress = getEnvOrDefault("DbAddress", "localhost:3306")
var DbName = getEnvOrDefault("DbName", "fresh_db")
var DbDriverName = getEnvOrDefault("DbDriverName", "mysql")
