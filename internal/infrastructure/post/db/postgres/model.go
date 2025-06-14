package postgres

import (
	"time"
)

type Image struct {
	UUID     string `gorm:"primaryKey;type:uuid"`
	URL      string `gorm:"not null;type:varchar(255)"`
	PostUUID string `gorm:"not null;type:uuid"`
	Post     Post   `gorm:"foreignKey:PostUUID;references:UUID;constraint:OnDelete:CASCADE"`
}

type Post struct {
	UUID      string    `gorm:"primaryKey;type:uuid"`
	Title     string    `gorm:"not null;type:varchar(255)"`
	Content   string    `gorm:"not null;type:text"`
	Images    []Image   `gorm:"foreignKey:PostUUID;references:UUID;constraint:OnDelete:CASCADE"`
	CreatedAt time.Time `gorm:"not null;autoCreateTime"`
	UpdatedAt time.Time `gorm:"not null;autoUpdateTime"`
}