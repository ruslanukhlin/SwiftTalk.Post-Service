package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

type PostgresConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
}

type Config struct {
	Mode     string
	Port     string
	Postgres *PostgresConfig
}

func LoadConfigFromEnv() *Config {
	_ = godotenv.Load(".env.local")

	fmt.Println(os.Getenv("POSTGRES_HOST"))

	return &Config{
		Mode:     os.Getenv("MODE"),
		Port:     os.Getenv("PORT"),
		Postgres: &PostgresConfig{
			Host:     os.Getenv("POSTGRES_HOST"),
			Port:     os.Getenv("POSTGRES_PORT"),
			User:     os.Getenv("POSTGRES_USER"),
			Password: os.Getenv("POSTGRES_PASSWORD"),
			DBName:   os.Getenv("POSTGRES_DB"),
		},
	}
}

func DNS(c *PostgresConfig) string {
	return fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		c.Host, c.User, c.Password, c.DBName, c.Port,
	)
}