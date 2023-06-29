package models

import (
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
)

type Snail struct {
	gorm.Model

	Name  string `json:"name"`
	OwnerID string  `json:"-"`
	Owner User   `json:"owner"  gorm:"references:DiscordID"` 
	Active bool   `json:"active" gorm:"default:false"`

	Level uint64 `json:"level" gorm:"default:1"`
	Exp   uint64 `json:"exp" gorm:"default:0"`
	Races uint64 `json:"races" gorm:"default:0"`
	Wins  uint64 `json:"wins" gorm:"default:0"`

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

func GetAllSnails(db *gorm.DB, owner User) ([]Snail, error) {
	snails := []Snail{}
	result := db.Where("owner_id = ?", owner.DiscordID).Find(&snails)
	return snails, result.Error
}

func GetActiveSnail(db *gorm.DB, owner User) (*Snail, error) {
	snail := &Snail{}
	result := db.Where("owner_id = ? AND active = ?", owner.DiscordID, true).First(snail)
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

// calcTernary is a ternary operator for float64s
func calcTernaryf(a, b float64, condition bool) float64 {
	if condition {
		return a
	}
	return b
}
