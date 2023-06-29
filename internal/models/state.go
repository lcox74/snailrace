package models

import "gorm.io/gorm"

type State struct {
	DB *gorm.DB
}

func NewState(db *gorm.DB) *State {
	return &State{DB: db}
}