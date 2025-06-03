package model

import "time"

type URL struct {
	ID          uint      `gorm:"primaryKey"`
	OriginalURL string    `gorm:"type:text;not null"`
	ShortCode   string    `gorm:"type:char(7);uniqueIndex;not null"`
	ClickCount  int       `gorm:"default:0"`
	CreatedAt   time.Time `gorm:"autoCreateTime"`
}
