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

	racePosition   float64 `json:"-" gorm:"-"`
	currentStamina float64 `json:"-" gorm:"-"`
}

func (s *Snail) NewRace() {
	s.racePosition = 0
	s.currentStamina = s.Stats.Stamina
}

// Step calculates the next step for the snail, based on the snail's stats and
// mood. This is still in testing stages and will probably be changed depending
// on how the game feels.
func (s *Snail) Step() {
	// Generate Random Bias
	bias := generateMoodBias(s.Mood)

	// Calculate max next step
	maxStepPotential := (s.Stats.Speed / 20.0 * MaxSnailStep) // ?

	if s.currentStamina > 0.0 {
		if bias >= (1.0 - maxStepPotential) {
			s.racePosition += MaxSnailStep
			s.currentStamina -= rand.Float64() * 2.0
		} else {
			s.racePosition += float64(rand.Intn(int(MaxSnailStep)))
			s.currentStamina -= rand.Float64()
		}

		s.racePosition -= rand.Float64()
	} else {
		s.currentStamina += s.Stats.Recovery / 10.0
	}

	// Make sure the snail doesn't go out of bounds
	s.racePosition = math.Min(s.racePosition, float64(MaxRaceLength))
	s.currentStamina = math.Max(0.0, s.currentStamina)
}

func (s Snail) renderPosition() string {
	trail := int((s.racePosition/float64(MaxRaceLength))*20.0) - 1
	line := strings.Repeat(".", int(math.Max(0.0, float64(trail))))
	line += "🐌"

	return fmt.Sprintf("%-20s", line)
}

func (s Snail) renderName(codeBlock bool) string {
	if s.Level > 0 && !codeBlock {
		return fmt.Sprintf("%s (<@%s>)", s.Name, s.OwnerID)
	}
	return s.Name
}

func CreateSnail(db *gorm.DB, owner User, levelType SnailStatLevel) (*Snail, error) {
	log.Debugf("CreateSnail(owner: %s, levelType: %v)", owner.DiscordID, levelType)

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
	log.Debugf("CreateDummySnail(levelType: %v)", levelType)

	snail := &Snail{Level: 0}
	snail.Stats.GenerateStats(levelType)
	snail.Name = generateSnailName()

	return snail
}

func GetAllSnails(db *gorm.DB, owner User) ([]Snail, error) {
	log.Debugf("GetAllSnails(owner: %s)", owner.DiscordID)

	snails := []Snail{}
	result := db.Where("owner_id = ?", owner.DiscordID).Preload("Owner").Find(&snails)
	return snails, result.Error
}

func GetActiveSnail(db *gorm.DB, owner User) (*Snail, error) {
	log.Debugf("GetActiveSnail(owner: %s)", owner.DiscordID)

	snail := &Snail{}
	result := db.Where("owner_id = ? AND active = ?", owner.DiscordID, true).Preload("Owner").First(snail)
	return snail, result.Error
}

func SetActiveSnail(db *gorm.DB, owner User, snail Snail) error {
	log.Debugf("SetActiveSnail(owner: %s, snail: %s)", owner.DiscordID, snail.Name)

	// Set all other snails to inactive
	db.Model(&Snail{}).Where("owner_id = ?", owner.DiscordID).Update("active", false)

	// Set the new snail to active
	result := db.Model(&snail).Update("active", true)
	return result.Error
}

func (snail *Snail) AddXP(db *gorm.DB, amount uint64) error {
	log.Debugf("GetActiveSnail(snail: %s, amount: %d)", snail.Name, amount)

	// Fetch the updated snail data
	db.Where("id = ?", snail.ID).First(snail)

	snail.Exp += amount
	if snail.Exp >= snail.Level*100 {
		snail.Exp -= snail.Level * 100
		snail.Level++
	}

	result := db.Save(snail)
	return result.Error
}

func (snail *Snail) AddRace(db *gorm.DB, win bool) error {
	log.Debugf("AddRace(snail: %s, win: %v)", snail.Name, win)

	// Fetch the updated snail data
	db.Where("id = ?", snail.ID).First(snail)

	snail.Races++
	if win {
		snail.Wins++
	}

	result := db.Save(snail)
	return result.Error
}

func generateSnailName() string {
	nounsFile, err := os.ReadFile("./res/snail_noun.txt")
	if err != nil {
		log.WithError(err).Warn("Error reading snail_noun.txt")
		return "buggy-snail"
	}
	adjectivesFile, err := os.ReadFile("./res/snail_adj.txt")
	if err != nil {
		log.WithError(err).Warn("Error reading snail_adj.txt")
		return "buggy-snail"
	}

	nouns := strings.Split(string(nounsFile), "\n")
	adjectives := strings.Split(string(adjectivesFile), "\n")

	return adjectives[rand.Intn(len(adjectives))] + "-" + nouns[rand.Intn(len(nouns))]
}

func generateMoodBias(mood float64) float64 {
	return rand.Float64() + mood
}
