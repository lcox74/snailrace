package models

import (
	"github.com/bwmarrin/discordgo"
	"github.com/google/uuid"
	"gorm.io/gorm"
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
	race := &Race{
		Id:        id,
		ChannelId: channelId,
		Host:      host,
	}
	s.Races[id] = race

	return race
}
