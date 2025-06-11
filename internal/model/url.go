package model

import "time"

type URL struct {
	ID          uint      `gorm:"primaryKey"`
	OriginalURL string    `gorm:"type:text;not null"`
	ShortCode   string    `gorm:"type:char(7);uniqueIndex;not null"`
	CreatedAt   time.Time `gorm:"autoCreateTime"`
}
