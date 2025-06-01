package gorm

import (
	"github.com/ruslanukhlin/SwiftTalk.post-service/internal/infrastructure/post/db/postgres"
	"github.com/ruslanukhlin/SwiftTalk.post-service/pkg/config"
)

func Migrate(config *config.Config) error {
	return DB.AutoMigrate(&postgres.Post{})
}