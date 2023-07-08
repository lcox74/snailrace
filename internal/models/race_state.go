package models

import (
	"fmt"
	"time"

	"github.com/bwmarrin/discordgo"
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type RaceStage uint8

const (
	// Race State Constants
	RaceStageOpen RaceStage = iota
	RaceStageBetting
	RaceStageRunning
	RaceStageFinished

	// State Timeout Constants
	RaceOpenTimeout      = 10 * time.Second
	RaceBettingTimeout   = 30 * time.Second
	RaceNoBettingTimeout = 10 * time.Second
	RaceStepInterval     = 1 * time.Second
	RaceTimeout          = 10 * time.Minute

	// Action Ids
	RaceActionJoin      = "host_join"
	RaceActionBet       = "host_bet"
	RaceActionBetAmount = "host_bet_amout"

	// Reward Constants
	BaseMoney = 10
	BaseXP    = 5
	WinPos1XP = 15
	WinPos2XP = 10
	WinPos3XP = 5
)

var (
	ErrInvalidSnail  = fmt.Errorf("invalid snail")
	ErrRaceClosed    = fmt.Errorf("race is closed")
	ErrAlreadyJoined = fmt.Errorf("snail already joined")
	ErrNotEnough     = fmt.Errorf("not enough racers")
	ErrBetsClosed    = fmt.Errorf("bets are closed")
)

type RaceBet struct {
	UserDiscordId string
	Amount        int
	SnailIndex    int
}

type RaceSnailPos struct {
	Position int
	Frame    int
	Snail    *Snail
}
type Race struct {
	Id        string
	ChannelId string
	Stage     RaceStage
	EndRace   func()
	DB        *gorm.DB

	Host    *discordgo.User
	Message *discordgo.Message

	NoBets   bool
	DontFill bool
	OnlyOne  bool

	Snails  []*Snail
	Bets    []RaceBet
	Odds    []float64
	Winners []RaceSnailPos
}

func (r *Race) SetupNewRace(id string, channelId string, db *gorm.DB, host *discordgo.User, endRace func()) {
	r.Id = id
	r.ChannelId = channelId
	r.Host = host
	r.EndRace = endRace
	r.Stage = RaceStageOpen
	r.Snails = make([]*Snail, 0)
	r.Bets = make([]RaceBet, 0)
	r.Odds = make([]float64, 0)
	r.Winners = make([]RaceSnailPos, 0)
	r.DB = db
}

// Flag setters
func (r *Race) SetNoBets() {
	r.NoBets = true
}
func (r *Race) SetDontFill() {
	r.DontFill = true
}
func (r *Race) SetOnlyOne() {
	r.OnlyOne = true
}

// If the race doesn't have the dont-fill flag, and the race has less than 4
// racers, then generate random snails to meet the 4 racer requirement.
func (r *Race) autoFillRace() {
	if r.DontFill {
		return
	}

	for len(r.Snails) < 4 {
		snail := CreateDummySnail(StartingSnail)
		r.Snails = append(r.Snails, snail)
	}
}

func (r *Race) AddSnail(snail *Snail) error {
	if r.Stage != RaceStageOpen {
		return ErrRaceClosed
	}

	for _, s := range r.Snails {
		if s.ID == snail.ID {
			return ErrAlreadyJoined
		}
	}
	r.Snails = append(r.Snails, snail)
	return nil
}

func (r *Race) GetSnail(index int) *Snail {
	if index < 0 || index >= len(r.Snails) {
		return nil
	}

	return r.Snails[index]
}
func (r *Race) PlaceBet(index int, amount int, userDiscordId string) error {
	if r.Stage != RaceStageBetting || r.NoBets {
		return ErrBetsClosed
	}

	if index < 0 || index >= len(r.Snails) {
		return ErrInvalidSnail
	}

	if len(r.Snails) < 2 {
		return ErrNotEnough
	}

	r.Bets = append(r.Bets, RaceBet{
		UserDiscordId: userDiscordId,
		Amount:        amount,
		SnailIndex:    index,
	})
	return nil
}

func StartRace(s *discordgo.Session, race *Race) {
	defer race.EndRace()

	raceStart := time.Now()
	defer func() {
		log.WithField("race", race.Id).Infof("Race took %s", time.Since(raceStart))
	}()

	log.WithField("race", race.Id).Infoln("Starting a race")
	race.Stage = RaceStageOpen
	if race.setupMessage(s) != nil {
		return
	}

	// Open Stage
	race.Render(s)
	time.Sleep(RaceOpenTimeout)
	race.Stage = RaceStageBetting

	// Autofill the Race
	if !race.DontFill {
		race.autoFillRace()
	}
	race.generateOdds()

	for _, snail := range race.Snails {
		log.WithFields(log.Fields{
			"race":     race.Id,
			"snail":    snail.Name,
			"speed":    snail.Stats.Speed,
			"stamina":  snail.Stats.Stamina,
			"recovery": snail.Stats.Recovery,
		}).Debugln("Entrant stats")
	}

	// Betting Stage
	race.Render(s)
	if race.NoBets {
		time.Sleep(RaceNoBettingTimeout)
	} else {
		time.Sleep(RaceBettingTimeout)
	}
	race.Stage = RaceStageRunning

	// Race Stage
	firstRace, raceAttempt := true, 0
	for firstRace || (race.racePosTie() && race.OnlyOne && raceAttempt < 5) {
		race.Render(s)
		firstRace = false
		raceAttempt++
		race.Winners = make([]RaceSnailPos, 0)

		// Reset the snails to start at the beginning
		for _, snail := range race.Snails {
			snail.NewRace()
		}

		snailsFinished, frame := 0, 0

		// Race until all snails have finished
		requiredFinished := len(race.Snails)
		for snailsFinished < requiredFinished {
			snailsFinished = 0
			for _, snail := range race.Snails {
				snail.Step()
				if snail.racePosition >= float64(MaxRaceLength) {
					snailsFinished++
					race.racePosAdd(snail, frame)
				}
			}
			race.Render(s)
			frame++
			time.Sleep(RaceStepInterval)
		}
	}

	// Finished Stage
	race.Stage = RaceStageFinished
	race.Payout(s)
	race.Render(s)
}

func (r *Race) Render(s *discordgo.Session) {
	switch r.Stage {
	case RaceStageOpen:
		r.renderOpenRace(s)
	case RaceStageBetting:
		r.renderBetting(s)
	case RaceStageRunning:
		r.renderRunning(s)
	case RaceStageFinished:
		r.renderFinished(s)
	}
}

func (r *Race) setupMessage(s *discordgo.Session) (err error) {
	r.Message, err = s.ChannelMessageSendComplex(r.ChannelId, &discordgo.MessageSend{
		Embeds: []*discordgo.MessageEmbed{
			{
				Title:       "Here Comes a New Race!",
				Description: "Loading...",
				Color:       0x2ecc71,
			},
		},
	})

	if err != nil {
		log.WithField("race", r.Id).WithError(err).Warnln("failed to send race setup message")
	}

	return err
}

func (r *Race) renderOpenRace(s *discordgo.Session) {
	// Build the Embed Message
	title := "Race: Open"
	body := fmt.Sprintf(
		"A new race has been hosted by %s\n\nRace ID: `%s`\n\nTo join via command, enter the following:\n```\n/snailrace join race_id: %s\n```\n**Entrants: (%d/12)**\n",
		r.Host.Username,
		r.Id,
		r.Id,
		len(r.Snails),
	)

	// Add the snails to the body as entrants `- <snail_name>(<@owner_id>)`
	for _, snail := range r.Snails {
		body += fmt.Sprintf("- %s\n", snail.renderName())
	}

	// Edit the message to reflect the current state of the race, in this
	// sense it will mainly update the entrants
	edit := discordgo.NewMessageEdit(r.ChannelId, r.Message.ID)
	edit.Embeds = []*discordgo.MessageEmbed{
		{
			Title:       title,
			Description: body,
			Color:       0x2ecc71,
		},
	}
	edit.Components = []discordgo.MessageComponent{
		discordgo.ActionsRow{
			Components: []discordgo.MessageComponent{
				discordgo.Button{
					Label:    "Join",
					Style:    discordgo.SuccessButton,
					CustomID: fmt.Sprintf("%s:%s", RaceActionJoin, r.Id),
				},
			},
		},
	}
	s.ChannelMessageEditComplex(edit)
}

func (r *Race) renderBetting(s *discordgo.Session) {

	if r.NoBets {
		r.renderNoBetting(s)
		return
	}

	// Build the Embed Message
	title := "Race: Bets are Open"
	body := fmt.Sprintf(
		"Bets are now open to everyone, do you feel lucky? To place a bet you can select the snail via the drop down. Here are the entrants:\n\nRace ID: `%s`\n\n**Entrants: (%d/12)**\n",
		r.Id,
		len(r.Snails),
	)

	select_options := make([]discordgo.SelectMenuOption, 0)

	// Add the snails to the body as entrants `index - <oods> <snail_name>(<@owner_id>)`
	for index, snail := range r.Snails {
		body += fmt.Sprintf("%2d - `%.02f` %s\n", index, r.Odds[index], snail.renderName())
		select_options = append(
			select_options,
			discordgo.SelectMenuOption{
				Label: snail.Name,
				Value: fmt.Sprintf("%d", index),
			},
		)
	}

	// Edit the message to reflect the current state of the race, in this
	// sense it will mainly update the entrants
	edit := discordgo.NewMessageEdit(r.ChannelId, r.Message.ID)
	edit.Embeds = []*discordgo.MessageEmbed{
		{
			Title:       title,
			Description: body,
			Color:       0x2ecc71,
		},
	}
	edit.Components = []discordgo.MessageComponent{
		discordgo.ActionsRow{
			Components: []discordgo.MessageComponent{
				discordgo.SelectMenu{
					MenuType: discordgo.StringSelectMenu,
					CustomID: fmt.Sprintf("%s:%s", RaceActionBet, r.Id),
					Options:  select_options,
				},
			},
		},
	}

	s.ChannelMessageEditComplex(edit)
}

func (r *Race) renderNoBetting(s *discordgo.Session) {
	// Build the Embed Message
	title := "Race: Ready to Race"
	body := fmt.Sprintf(
		"We are ready to race `%s`, here are the entrants:\n\n**Entrants: (%d/12)**\n",
		r.Id,
		len(r.Snails),
	)

	// Add the snails to the body as entrants `index - <odds> <snail_name>(<@owner_id>)`
	for index, snail := range r.Snails {
		body += fmt.Sprintf("%2d - `%.02f` %s\n", index, r.Odds[index], snail.renderName())
	}

	// Edit the message to reflect the current state of the race, in this
	// sense it will mainly update the entrants
	edit := discordgo.NewMessageEdit(r.ChannelId, r.Message.ID)
	edit.Embeds = []*discordgo.MessageEmbed{
		{
			Title:       title,
			Description: body,
			Color:       0x2ecc71,
		},
	}
	edit.Components = []discordgo.MessageComponent{}

	s.ChannelMessageEditComplex(edit)
}

func (r *Race) renderRunning(s *discordgo.Session) {
	title := "Race: Ready to Race"
	body := ""

	entrants := fmt.Sprintf("**Entrants: (%d/12):**\n", len(r.Snails))

	track := fmt.Sprintf("```\nRace ID: %s\n\n", r.Id)
	track += "                          🏁\n"
	track += "  |-----------------------|\n"

	// Build snails
	for index, snail := range r.Snails {
		line := snail.renderPosition()

		// Render the snail on the track
		row := fmt.Sprintf("%2d| %s | %s\n", index, line, snail.Name)
		if pos := r.racePosPosition(snail); pos > 0 {
			row = fmt.Sprintf("%2d| %s %d %s\n", index, line, pos, snail.Name)
		}
		track += row

		// Render the entrant in the list
		row = fmt.Sprintf("%2d - `%.02f` %s\n", index, r.Odds[index], snail.renderName())
		entrants += row
	}
	track += "  |-----------------------|\n```"

	body += track + entrants

	// Edit the message to reflect the current state of the race, in this
	// sense it will mainly update the entrants
	edit := discordgo.NewMessageEdit(r.ChannelId, r.Message.ID)
	edit.Embeds = []*discordgo.MessageEmbed{
		{
			Title:       title,
			Description: body,
			Color:       0x2ecc71,
		},
	}
	edit.Components = []discordgo.MessageComponent{}

	s.ChannelMessageEditComplex(edit)
}
func (r *Race) renderFinished(s *discordgo.Session) {
	title := "Race: Complete"
	body := r.getWinnersStr() + "\n\n"

	entrants := fmt.Sprintf("**Entrants: (%d/12):**\n", len(r.Snails))

	track := fmt.Sprintf("```\nRace ID: %s\n\n", r.Id)
	track += "                          🏁\n"
	track += "  |-----------------------|\n"

	// Build snails
	for index, snail := range r.Snails {
		line := snail.renderPosition()

		// Render the snail on the track
		row := fmt.Sprintf("%2d| %s | %s\n", index, line, snail.Name)
		if pos := r.racePosPosition(snail); pos > 0 {
			row = fmt.Sprintf("%2d| %s %d %s\n", index, line, pos, snail.Name)
		}
		track += row

		// Render the entrant in the list
		row = fmt.Sprintf("%2d - `%.02f` %s\n", index, r.Odds[index], snail.renderName())
		entrants += row
	}
	track += "  |-----------------------|\n```"

	body += track + entrants

	// Edit the message to reflect the current state of the race, in this
	// sense it will mainly update the entrants
	edit := discordgo.NewMessageEdit(r.ChannelId, r.Message.ID)
	edit.Embeds = []*discordgo.MessageEmbed{
		{
			Title:       title,
			Description: body,
			Color:       0x2ecc71,
		},
	}
	edit.Components = []discordgo.MessageComponent{}

	s.ChannelMessageEditComplex(edit)

}

func (r Race) getWinnersStr() string {
	winners := make([]*Snail, 0)
	for _, racePos := range r.Winners {
		if racePos.Position == 1 {
			winners = append(winners, racePos.Snail)
		}
	}

	switch len(winners) {
	case 1:
		return fmt.Sprintf("Wow congratulations %s for winning the race!", winners[0].renderName())
	case 2:
		return fmt.Sprintf("It's a Tie! Good job %s and %s!", winners[0].renderName(), winners[1].renderName())
	case 3:
		return fmt.Sprintf("It's a three way tie! Good job %s, %s, and %s!", winners[0].renderName(), winners[1].renderName(), winners[2].renderName())
	}

	winStr := "I don't believe it... "
	for _, winner := range winners {
		winStr += winner.renderName()
	}
	winStr += " have all won..."
	return winStr
}

func (r *Race) Payout(s *discordgo.Session) {

	// Give Snails Base XP
	for _, snail := range r.Snails {
		if snail.Level == 0 {
			continue
		}

		switch r.racePosPosition(snail) {
		case 1:
			snail.AddXP(r.DB, uint64(BaseXP+(WinPos1XP*len(r.Snails))))
			snail.Owner.AddXP(r.DB, uint64(BaseXP+(WinPos1XP*len(r.Snails))))
			snail.Owner.AddMoney(r.DB, uint64(BaseMoney*len(r.Snails)))
			snail.AddRace(r.DB, true)
			snail.Owner.AddRace(r.DB, true)
		case 2:
			snail.AddXP(r.DB, uint64(BaseXP+(WinPos2XP*len(r.Snails))))
			snail.Owner.AddXP(r.DB, uint64(BaseXP+(WinPos2XP*len(r.Snails))))
			snail.AddRace(r.DB, false)
			snail.Owner.AddRace(r.DB, false)
		case 3:
			snail.AddXP(r.DB, uint64(BaseXP+(WinPos3XP*len(r.Snails))))
			snail.Owner.AddXP(r.DB, uint64(BaseXP+(WinPos3XP*len(r.Snails))))
			snail.AddRace(r.DB, false)
			snail.Owner.AddRace(r.DB, false)
		default:
			snail.AddXP(r.DB, uint64(BaseXP))
			snail.Owner.AddXP(r.DB, uint64(BaseXP))
			snail.AddRace(r.DB, false)
			snail.Owner.AddRace(r.DB, false)
		}
	}

	// Calculate the payout for each bet
	for _, bet := range r.Bets {
		if r.racePosPosition(r.Snails[bet.SnailIndex]) == 1 {
			// Get the user who placed the bet
			user, err := GetUserByDiscordID(r.DB, bet.UserDiscordId)
			if err != nil {
				log.WithField("race", r.Id).WithError(err).Warnln("Failed to get user for payout")
				continue
			}

			// Calculate the payout
			payout := uint64(float64(bet.Amount) * r.Odds[bet.SnailIndex])
			user.AddMoney(r.DB, payout)
		}
	}
}

// Uses the each snails stats, create the odds of the each snail winning. The
// lower the number the more likely the snail is to win. The odds are based on
// the normalized stats of the snail, with a modifier based on the snails win
// history. The Odds will be used to calculate the payout for each bet.
func (r *Race) generateOdds() {
	r.Odds = make([]float64, len(r.Snails))

	// Pre-calculate the sum of the speed, and stamina stats to normalize the
	// stats for each snail later.
	sum_speed, sum_stamina := 0.0, 0.0
	for _, snail := range r.Snails {
		sum_speed += snail.Stats.Speed
		sum_stamina += snail.Stats.Stamina
	}

	// Generate for each snail
	for index, snail := range r.Snails {
		// Calculate modifier from normalized stats
		norm_speed := snail.Stats.Speed / sum_speed
		norm_stamina := snail.Stats.Stamina / sum_stamina
		modifier := 1.0 - (norm_speed + norm_stamina)

		// Check the snail's win history
		win_rate := 1.0
		if snail.Wins > 0 {
			win_rate = 1.0 - (float64(snail.Wins) / float64(snail.Races))
			if win_rate == 0.0 {
				win_rate = 1.0
			}
		}

		// Limit the odd
		odd := 10.0 * modifier * win_rate
		if odd < 1.0 {
			odd = 1.0
		}
		r.Odds[index] = odd
	}
}

// Checks if the snail is already in the winners list
func (r Race) racePosContains(snail *Snail) bool {
	for _, p := range r.Winners {
		if p.Snail == snail {
			return true
		}
	}
	return false
}

func (r Race) racePosPosition(snail *Snail) int {
	for _, p := range r.Winners {
		if p.Snail == snail {
			return p.Position
		}
	}
	return 0
}

// When a snail crosses the line, add it to the winners list. If the snail is
// already in the list, do nothing. If there is already a snail with the same
// frame, then the snail is tied with the other snail.
func (r *Race) racePosAdd(snail *Snail, frame int) {
	if r.racePosContains(snail) {
		return
	}

	pos := 1
	for _, p := range r.Winners {
		if p.Position >= pos {
			pos = p.Position + 1
		}

		// Check for a possible tie
		if p.Frame == frame {
			pos = p.Position
			break
		}
	}

	// Add the snail to the winners list
	r.Winners = append(r.Winners, RaceSnailPos{
		Position: pos,
		Snail:    snail,
		Frame:    frame,
	})
}

// Check the race for a tie, it doesn't matter how many are in the tie, just
// that there is a tie.
func (r Race) racePosTie() bool {
	for _, a := range r.Winners {
		for _, b := range r.Winners {
			if a.Position == b.Position && a.Snail != b.Snail {
				return true
			}
		}
	}
	return false
}
