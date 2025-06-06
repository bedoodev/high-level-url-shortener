package model

import "time"

type ClickEvent struct {
	ID        uint      `gorm:"primaryKey"`
	URLID     uint      `gorm:"not null;index"` // Foreign key field
	URL       URL       `gorm:"foreignKey:URLID;constraint:OnDelete:CASCADE"`
	Timestamp time.Time `gorm:"autoCreateTime"`
}
