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

type S3Config struct {
	Bucket string
	BucketUrl string
}

type Config struct {
	Mode         string
	PortGrpc     string
	PortHttp     string
	Postgres     *PostgresConfig
	S3           *S3Config
}

func LoadConfigFromEnv() *Config {
	_ = godotenv.Load(".env.local")

	return &Config{
		Mode:         os.Getenv("MODE"),
		PortGrpc:     os.Getenv("PORT_GRPC"),
		PortHttp:     os.Getenv("PORT_HTTP"),
		Postgres: &PostgresConfig{
			Host:     os.Getenv("POSTGRES_HOST"),
			Port:     os.Getenv("POSTGRES_PORT"),
			User:     os.Getenv("POSTGRES_USER"),
			Password: os.Getenv("POSTGRES_PASSWORD"),
			DBName:   os.Getenv("POSTGRES_DB"),
		},
		S3: &S3Config{
			Bucket:    os.Getenv("S3_BUCKET"),
			BucketUrl: os.Getenv("S3_BUCKET_URL"),
		},
	}
}

func DNS(c *PostgresConfig) string {
	return fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		c.Host, c.User, c.Password, c.DBName, c.Port,
	)
}