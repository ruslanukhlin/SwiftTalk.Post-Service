package config

import (
	"fmt"
	"log"
	"os"
	"sync"
)

var (
	once sync.Once
	cfg  *Config
)

type AuthConfig struct {
	Host string
	Port string
}

type PostgresConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
}

type S3Config struct {
	Bucket       string
	BucketUrl    string
	BucketFolder string
}

type Config struct {
	Mode     string
	PortGrpc string
	PortHttp string
	Postgres *PostgresConfig
	S3       *S3Config
	Auth     *AuthConfig
}

func LoadConfigFromEnv() *Config {
	once.Do(func() {
		// _ = godotenv.Load(".env.local")

		cfg = &Config{
			Mode:     os.Getenv("MODE"),
			PortGrpc: os.Getenv("PORT_GRPC"),
			PortHttp: os.Getenv("PORT_HTTP"),
			Postgres: &PostgresConfig{
				Host:     os.Getenv("POSTGRES_HOST"),
				Port:     os.Getenv("POSTGRES_PORT"),
				User:     os.Getenv("POSTGRES_USER"),
				Password: os.Getenv("POSTGRES_PASSWORD"),
				DBName:   os.Getenv("POSTGRES_DB"),
			},
			S3: &S3Config{
				Bucket:       os.Getenv("S3_BUCKET"),
				BucketUrl:    os.Getenv("S3_BUCKET_URL"),
				BucketFolder: os.Getenv("S3_BUCKET_FOLDER"),
			},
			Auth: &AuthConfig{
				Host: os.Getenv("AUTH_HOST"),
				Port: os.Getenv("AUTH_PORT"),
			},
		}
	})

	return cfg
}

func DNS(c *PostgresConfig) string {
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		c.Host, c.User, c.Password, c.DBName, c.Port,
	)

	// Логируем строку подключения (без пароля)
	log.Printf("Attempting to connect to PostgreSQL with: host=%s user=%s dbname=%s port=%s sslmode=disable",
		c.Host, c.User, c.DBName, c.Port)

	return dsn
}
