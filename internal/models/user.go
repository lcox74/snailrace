package models

import "gorm.io/gorm"

type User struct {
	gorm.Model

	DiscordID string `gorm:"uniqueIndex"`

	Level uint64 `gorm:"default:0"`
	XP    uint64 `gorm:"default:0"`

	Races uint64 `gorm:"default:0"`
	Wins  uint64 `gorm:"default:0"`

	Money uint64 `gorm:"default:10"`
}

func GetUserByDiscordID(db *gorm.DB, discordID string) (*User, error) {
	user := &User{}
	result := db.Where("discord_id = ?", discordID).First(user).Preload("ActiveSnail")
	return user, result.Error
}