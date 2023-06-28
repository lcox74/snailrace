package models

import (
	"math"
	"math/rand"

	"gorm.io/gorm"
)

type SnailMood int

const (
	MoodSad     SnailMood = -1
	MoodHappy   SnailMood = 1
	MoodFocused SnailMood = 0
)

type Snail struct {
	gorm.Model

	Name  string `json:"name"`
	Owner uint64 `json:"owner"`

	Level uint64 `json:"level"`
	Exp   uint64 `json:"exp"`
	Races uint64 `json:"races"`
	Wins  uint64 `json:"wins"`

	Mood  float64    `json:"mood" gorm:"default:0"`
	Stats SnailStats `json:"stats" gorm:"embedded"`

	racePosition int `json:"-" gorm:"-"`
	lastStep     int `json:"-" gorm:"-"`
}

// Step calculates the next step for the snail, based on the snail's stats and
// mood. This is still in testing stages and will probably be changed depending
// on how the game feels. TODO: Test and Change when needed
func (s *Snail) Step() {
	// Generate Random Bias
	bias := generateMoodBias(s.Mood)

	// Calculate base interval before bias and acceleration
	maxStep := 10.0 + s.Stats.Speed
	minStep := math.Min(s.Stats.Stamina, maxStep-5)
	avgStep := (maxStep + minStep) / 2.0

	// Calculate acceleration factor with weight and prevStep
	acceleration := (s.Stats.Weight-5)/5.0 + (float64(s.lastStep)-avgStep)/5.0
	minStep = math.Max(
		0,
		minStep+(calcTernaryf(-1, 1, s.Stats.Weight < 5.0))*acceleration+bias,
	)
	maxStep = math.Min(
		20,
		maxStep+(calcTernaryf(1, -1, s.Stats.Weight < 5.0))*acceleration+bias,
	)

	// Calculate new position
	s.lastStep = int(rand.Float64()*(maxStep-minStep) + minStep)
	s.racePosition = int(math.Min(float64(s.racePosition+s.lastStep), 100))
}

func generateMoodBias(mood float64) float64 {
	return rand.Float64() + mood
}

// calcTernary is a ternary operator for float64s
func calcTernaryf(a, b float64, condition bool) float64 {
	if condition {
		return a
	}
	return b
}
