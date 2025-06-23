package models

import "time"

// ID qui est une primaryKey
// Shortcode : doit être unique, indexé pour des recherches rapide (voir doc), taille max 10 caractères
// LongURL : doit pas être null
// CreateAt : Horodatage de la créatino du lien

type Link struct {
	ID        uint      `gorm:"primaryKey"`
	ShortCode string    `gorm:"uniqueIndex;size:10;not null"`
	LongURL   string    `gorm:"not null"`
	CreatedAt time.Time `gorm:"autoCreateTime"`
}
