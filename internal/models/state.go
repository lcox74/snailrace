package models

import (
	"github.com/bwmarrin/discordgo"
	"github.com/google/uuid"
	"gorm.io/gorm"

	log "github.com/sirupsen/logrus"
)

type State struct {
	DB    *gorm.DB
	Races map[string]*Race
}

func NewState(db *gorm.DB) *State {
	return &State{
		DB:    db,
		Races: make(map[string]*Race, 0),
	}
}

func (s *State) NewRace(session *discordgo.Session, channelId string, host *discordgo.User) *Race {
	// Generate Unique ID
	id := uuid.New().String()[24:]
	_, ok := s.Races[id]
	for ok {
		id = uuid.New().String()[24:]
		_, ok = s.Races[id]
	}

	// Create New Race
	race := &Race{}
	race.SetupNewRace(id, channelId, s.DB, host, func() {
		delete(s.Races, id)
		log.Printf("Race %s has come to a close", id)
	})
	s.Races[id] = race

	return race
}
