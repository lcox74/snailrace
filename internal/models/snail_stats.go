package models

import (
	"fmt"
	"math"
	"math/rand"
	"strings"
)

type SnailStats struct {
	Speed    float64 `json:"speed"`
	Stamina  float64 `json:"stamina"`
	Recovery float64 `json:"recovery"`
}

type SnailStatLevel uint8

const (
	StartingSnail SnailStatLevel = iota
	AmateurSnail
	ProfessionalSnail
	ExpertSnail
	RandomSnail
)

func CanUserAffordSnail(user User, tier SnailStatLevel) (bool, int) {
	prices := map[int]int{
		0: 50,
		1: 100,
		2: 150,
		3: 250,
	}
	snailPrice := prices[int(tier)]
	return int(user.Money) >= snailPrice, snailPrice
}

func (s SnailStats) RenderStatBlock(level SnailStatLevel) string {
	return fmt.Sprintf(
		"%-9s%s %.02f\n%-9s%s %.02f\n%-9s%s %.02f\n",
		"Speed", renderStat(s.Speed, level), s.Speed,
		"Stamina", renderStat(s.Stamina, level), s.Stamina,
		"Recovery", renderStat(s.Recovery, level), s.Recovery,
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
	s.Recovery = randFloat64(1, 5)
}

// Amateur Snail Stats are randomly generated between 5 and 10 for each stat
func (s *SnailStats) generateAmateurStats() {
	s.Speed = randFloat64(5, 10)
	s.Stamina = randFloat64(5, 10)
	s.Recovery = randFloat64(5, 10)
}

// Professional Snail Stats are randomly generated between 10 and 15 for each
// stat
func (s *SnailStats) generateProfessionalStats() {
	s.Speed = randFloat64(10, 15)
	s.Stamina = randFloat64(10, 15)
	s.Recovery = randFloat64(10, 15)
}

// Expert Snail Stats are randomly generated between 15 and 20 for each stat
func (s *SnailStats) generateExpertStats() {
	s.Speed = randFloat64(15, 20)
	s.Stamina = randFloat64(15, 20)
	s.Recovery = randFloat64(15, 20)
}

// Random Snail Stats are randomly generated between 1 and 20 for each stat
func (s *SnailStats) generateRandomStats() {
	s.Speed = randFloat64(1, 20)
	s.Stamina = randFloat64(1, 20)
	s.Recovery = randFloat64(1, 20)
}

func randFloat64(min, max float64) float64 {
	return rand.Float64()*(max-min) + min
}

// This fixes a bug where the stat would surpass the maximum number
// of characters rendered, thus making a negative repeat error.
// This calculates what the number and remainder would be respective
// to the snail tier
func CalcNumAndMax(stat float64, level SnailStatLevel) (float64, float64) {
	tierProps := map[SnailStatLevel]float64{
		0: 10,
		1: 20,
		2: 30,
		3: 40,
		4: 40,
	}
	maxValue := tierProps[level]
	number := math.Floor(stat * 10 / maxValue)
	return number, 10 - number
}

func renderStat(stat float64, level SnailStatLevel) string {
	num, rem := CalcNumAndMax(stat, level)
	return fmt.Sprintf("[%s%s]", strings.Repeat("=", int(num)), strings.Repeat(" ", int(rem)))
}
