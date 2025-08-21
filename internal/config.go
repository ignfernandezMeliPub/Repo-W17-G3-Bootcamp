package internal

import (
	"github.com/melisource/fury_go-platform/pkg/fury/secret"
	"github.com/melisource/fury_go-toolkit-secrets/pkg/secrets"
)

func getEnvOrDefault(client secrets.Client, key string, defaultValue string) string {
	value, exists := client.GetSecret(key)

	if !exists {
		return defaultValue
	}
	return value
}

func newSecretsClient() secrets.Client {
	client, err := secrets.NewClient()
	if err != nil {
		panic(err)
	}
	return client
}

var secretsClient = newSecretsClient()

// ? Server config
var ServerAddress = getEnvOrDefault(secretsClient, "ServerAddress", ":8080")

// ? DB Config
var DbUser = getEnvOrDefault(secretsClient, "DB_MYSQL_DESAENV10_BGOW17S434_BGOW17S434_WPROD_USER", "root")
var DbPassword = getEnvOrDefault(secretsClient, "DB_MYSQL_DESAENV10_BGOW17S434_BGOW17S434_WPROD", "")
var DbProtocol = getEnvOrDefault(secretsClient, "DbProtocol", "tcp")
var DbAddress = secret.FromEnv("DB_MYSQL_DESAENV10_BGOW17S434_BGOW17S434_ENDPOINT")

//getEnvOrDefault(secretsClient, "DB_MYSQL_DESAENV10_BGOW17S434_BGOW17S434_ENDPOINT", "localhost:3306")

var DbName = getEnvOrDefault(secretsClient, "DB_MYSQL_NAME", "fresh_db") //  bgow17s434
var DbDriverName = getEnvOrDefault(secretsClient, "DbDriverName", "mysql")

// Config estructura que contiene toda la configuración
type Config struct {
	Server ServerConfig
	DB     DBConfig
}

type ServerConfig struct {
	Address string
}

type DBConfig struct {
	User                 string
	Password             string
	Protocol             string
	Address              string
	Name                 string
	Driver               string
	AllowNativePasswords bool
}

// GetConfig retorna la configuración completa
func GetConfig() *Config {
	return &Config{
		Server: ServerConfig{
			Address: ServerAddress,
		},
		DB: DBConfig{
			User:                 DbUser,
			Password:             DbPassword,
			Protocol:             DbProtocol,
			Address:              DbAddress,
			Name:                 DbName,
			Driver:               DbDriverName,
			AllowNativePasswords: true,
		},
	}
}
