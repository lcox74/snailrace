package models

import (
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

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
	log.Debugf("GetUserByDiscordID(id: %s)", discordID)

	user := &User{}
	result := db.Where("discord_id = ?", discordID).First(user)
	return user, result.Error
}

func CreateUser(db *gorm.DB, discordID string) (*User, error) {
	log.Debugf("CreateUser(id: %s)", discordID)

	user := &User{DiscordID: discordID}
	result := db.Create(user)
	return user, result.Error
}

func (user *User) RemoveMoney(db *gorm.DB, amount uint64) error {
	log.Debugf("RemoveMoney(id: %s, amount: %d)", user.DiscordID, amount)

	// Fetch the Current User Data
	db.Where("discord_id = ?", user.DiscordID).First(user)

	// Modify the users money
	user.Money -= amount
	result := db.Save(user)
	return result.Error
}

func (user *User) AddMoney(db *gorm.DB, amount uint64) error {
	log.Debugf("AddMoney(id: %s, amount: %d)", user.DiscordID, amount)

	// Fetch the Current User Data
	db.Where("discord_id = ?", user.DiscordID).First(user)

	// Modify the users money
	user.Money += amount
	result := db.Save(user)
	return result.Error
}
func (user *User) AddXP(db *gorm.DB, amount uint64) error {
	log.Debugf("AddXP(id: %s, amount: %d)", user.DiscordID, amount)

	// Fetch the Current User Data
	db.Where("discord_id = ?", user.DiscordID).First(user)

	user.XP += amount
	if user.XP >= user.Level*100 {
		user.XP -= user.Level * 100
		user.Level++
	}

	result := db.Save(user)
	return result.Error
}

func (user *User) AddRace(db *gorm.DB, win bool) error {
	log.Debugf("AddRace(id: %s, win: %v)", user.DiscordID, win)

	// Fetch the Current User Data
	db.Where("discord_id = ?", user.DiscordID).First(user)

	user.Races++
	if win {
		user.Wins++
	}

	result := db.Save(user)
	return result.Error
}
