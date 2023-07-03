package models

import (
	"fmt"
	"math"
	"math/rand"
	"os"
	"strings"

	"gorm.io/gorm"

	log "github.com/sirupsen/logrus"
)

type SnailMood int

const (
	MoodSad     SnailMood = -1
	MoodHappy   SnailMood = 1
	MoodFocused SnailMood = 0

	MaxSnailStep  float64 = 5.0
	MaxRaceLength int     = 100
)

type Snail struct {
	gorm.Model

	Name    string `json:"name"`
	OwnerID string `json:"-"`
	Owner   User   `json:"owner"  gorm:"references:DiscordID"`
	Active  bool   `json:"active" gorm:"default:false"`

	Level uint64 `json:"level" gorm:"default:1"`
	Exp   uint64 `json:"exp" gorm:"default:0"`
	Races uint64 `json:"races" gorm:"default:0"`
	Wins  uint64 `json:"wins" gorm:"default:0"`

	Mood  float64    `json:"mood" gorm:"default:0"`
	Stats SnailStats `json:"stats" gorm:"embedded"`

	racePosition   int     `json:"-" gorm:"-"`
	currentStamina float64 `json:"-" gorm:"-"`
}

func (s *Snail) NewRace() {

}

// Step calculates the next step for the snail, based on the snail's stats and
// mood. This is still in testing stages and will probably be changed depending
// on how the game feels. TODO: Test and Change when needed
func (s *Snail) Step() {
	// Generate Random Bias
	bias := generateMoodBias(s.Mood)

	if s.currentStamina > 0 {
		// Calculate max step the snail can take
		maxStep := s.Stats.Speed * (s.currentStamina / s.Stats.Stamina)

		// Calculate the next step
		step := maxStep * (s.Stats.Weight / 10.0)
		step += bias

		// Set the new position
		s.racePosition += int(math.Round(step))
		s.currentStamina--

	} else {
		s.currentStamina += 2.0 * bias
	}

	if bias > 0.5 {
		s.currentStamina += 5 * (s.Stats.Stamina / 20.0)
	}

	// Make sure the snail doesn't go out of bounds
	s.racePosition = int(math.Min(float64(s.racePosition), float64(MaxRaceLength)))
}

func (s Snail) renderPosition() string {
	trail := int((float64(s.racePosition)/float64(MaxRaceLength))*20.0) - 1
	line := strings.Repeat(".", int(math.Max(0.0, float64(trail))))
	line += "🐌"

	log.Printf("Snail: %s, pos: %d, cst: %f, sp: %f, st: %f, wt: %f", s.Name, s.racePosition, s.currentStamina, s.Stats.Speed, s.Stats.Stamina, s.Stats.Weight)

	return fmt.Sprintf("%-20s", line)
}

func CreateSnail(db *gorm.DB, owner User, levelType SnailStatLevel) (*Snail, error) {
	snail := &Snail{
		Owner: owner,
		Level: 1,
	}
	snail.Stats.GenerateStats(levelType)
	snail.Name = generateSnailName()

	result := db.Create(snail)
	return snail, result.Error
}

func CreateDummySnail(levelType SnailStatLevel) *Snail {
	snail := &Snail{Level: 0}
	snail.Stats.GenerateStats(levelType)
	snail.Name = generateSnailName()

	return snail
}

func GetAllSnails(db *gorm.DB, owner User) ([]Snail, error) {
	snails := []Snail{}
	result := db.Where("owner_id = ?", owner.DiscordID).Preload("Owner").Find(&snails)
	return snails, result.Error
}

func GetActiveSnail(db *gorm.DB, owner User) (*Snail, error) {
	snail := &Snail{}
	result := db.Where("owner_id = ? AND active = ?", owner.DiscordID, true).Preload("Owner").First(snail)
	return snail, result.Error
}

func SetActiveSnail(db *gorm.DB, owner User, snail Snail) error {
	// Set all other snails to inactive
	db.Model(&Snail{}).Where("owner_id = ?", owner.DiscordID).Update("active", false)

	// Set the new snail to active
	result := db.Model(&snail).Update("active", true)
	return result.Error
}

func generateSnailName() string {
	nounsFile, err := os.ReadFile("./res/snail_noun.txt")
	if err != nil {
		log.Warnf("Error reading snail_noun.txt: %v", err)
		return "buggy-snail"
	}
	adjectivesFile, err := os.ReadFile("./res/snail_adj.txt")
	if err != nil {
		log.Warnf("Error reading snail_adj.txt: %v", err)
		return "buggy-snail"
	}

	nouns := strings.Split(string(nounsFile), "\n")
	adjectives := strings.Split(string(adjectivesFile), "\n")

	return adjectives[rand.Intn(len(adjectives))] + "-" + nouns[rand.Intn(len(nouns))]
}

func generateMoodBias(mood float64) float64 {
	return rand.Float64() + mood
}
