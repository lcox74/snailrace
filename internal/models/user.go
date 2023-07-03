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
	result := db.Where("discord_id = ?", discordID).First(user)
	return user, result.Error
}

func CreateUser(db *gorm.DB, discordID string) (*User, error) {
	user := &User{DiscordID: discordID}
	result := db.Create(user)
	return user, result.Error
}

func (user *User) RemoveMoney(db *gorm.DB, amount uint64) error {
	user.Money -= amount
	result := db.Save(user)
	return result.Error
}

func (user *User) AddMoney(db *gorm.DB, amount uint64) error {
	user.Money += amount
	result := db.Save(user)
	return result.Error
}
