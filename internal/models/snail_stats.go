package models

import (
	"fmt"
	"math/rand"
	"strings"
)

type SnailStats struct {
	Speed   float64 `json:"speed"`
	Stamina float64 `json:"stamina"`
	Weight  float64 `json:"strength"`
}

type SnailStatLevel uint8

const (
	StartingSnail SnailStatLevel = iota
	AmateurSnail
	ProfessionalSnail
	ExpertSnail
	RandomSnail
)

func (s SnailStats) RenderStatBlock() string {
	return fmt.Sprintf(
		"%-9s%s\n%-9s%s\n%-9s%s\n",
		"Speed", renderStat(s.Speed),
		"Stamina", renderStat(s.Stamina),
		"Strength", renderStat(s.Weight),
	)
}

func (s *SnailStats) GenerateStats(level SnailStatLevel) {
	switch level {
	case StartingSnail:
		s.generateStartingStats()
	case AmateurSnail:
		s.generateAmateurStats()
	case ProfessionalSnail:
		s.generateProfessionalStats()
	case ExpertSnail:
		s.generateExpertStats()
	default:
		s.generateRandomStats()
	}
}

// Starting Snail Stats are randomly generated between 1 and 5 for each stat
func (s *SnailStats) generateStartingStats() {
	s.Speed = randFloat64(1, 5)
	s.Stamina = randFloat64(1, 5)
	s.Weight = randFloat64(1, 5)
}

// Amateur Snail Stats are randomly generated between 5 and 10 for each stat
func (s *SnailStats) generateAmateurStats() {
	s.Speed = randFloat64(5, 10)
	s.Stamina = randFloat64(5, 10)
	s.Weight = randFloat64(5, 10)
}

// Professional Snail Stats are randomly generated between 10 and 15 for each
// stat
func (s *SnailStats) generateProfessionalStats() {
	s.Speed = randFloat64(10, 15)
	s.Stamina = randFloat64(10, 15)
	s.Weight = randFloat64(10, 15)
}

// Expert Snail Stats are randomly generated between 15 and 20 for each stat
func (s *SnailStats) generateExpertStats() {
	s.Speed = randFloat64(15, 20)
	s.Stamina = randFloat64(15, 20)
	s.Weight = randFloat64(15, 20)
}

// Random Snail Stats are randomly generated between 1 and 20 for each stat
func (s *SnailStats) generateRandomStats() {
	s.Speed = randFloat64(1, 20)
	s.Stamina = randFloat64(1, 20)
	s.Weight = randFloat64(1, 20)
}

func randFloat64(min, max float64) float64 {
	return rand.Float64()*(max-min) + min
}

func renderStat(stat float64) string {
	num := int(stat)
	return fmt.Sprintf("[%s%s]", strings.Repeat("=", num), strings.Repeat(" ", 10-num))
}
