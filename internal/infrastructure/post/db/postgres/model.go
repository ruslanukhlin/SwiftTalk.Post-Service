package postgres

import (
	"time"

	"github.com/lib/pq"
)

type Post struct {
	UUID      string         `gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	Title     string         `gorm:"not null;type:varchar(255)"`
	Content   string         `gorm:"not null;type:text"`
	Images    pq.StringArray `gorm:"type:text[]"`
	CreatedAt time.Time      `gorm:"not null;autoCreateTime"`
	UpdatedAt time.Time      `gorm:"not null;autoUpdateTime"`
}