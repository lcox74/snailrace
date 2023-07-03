package models

import (
	"fmt"
	"time"

	"github.com/bwmarrin/discordgo"
	log "github.com/sirupsen/logrus"
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
	RaceTimeout          = 10 * time.Minute

	// Action Ids
	RaceActionJoin      = "host_join"
	RaceActionBet       = "host_bet"
	RaceActionBetAmount = "host_bet_amout"
)

var (
	ErrInvalidSnail  = fmt.Errorf("invalid snail")
	ErrRaceClosed    = fmt.Errorf("race is closed")
	ErrAlreadyJoined = fmt.Errorf("snail already joined")
	ErrBetsClosed    = fmt.Errorf("bets are closed")
)

type RaceBet struct {
	UserDiscordId string
	Amount        int
	SnailIndex    int
}

type Race struct {
	Id        string
	ChannelId string
	Stage     RaceStage
	EndRace   func()

	Host    *discordgo.User
	Message *discordgo.Message

	NoBets   bool
	DontFill bool
	OnlyOne  bool

	Snails []*Snail
	Bets   []RaceBet
	Odds   []float64
}

func (r *Race) SetupNewRace(id string, channelId string, host *discordgo.User, endRace func()) {
	r.Id = id
	r.ChannelId = channelId
	r.Host = host
	r.EndRace = endRace
	r.Stage = RaceStageOpen
	r.Snails = make([]*Snail, 0)
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

	r.Bets = append(r.Bets, RaceBet{
		UserDiscordId: userDiscordId,
		Amount:        amount,
		SnailIndex:    index,
	})
	return nil
}

func StartRace(s *discordgo.Session, race *Race) {
	log.Printf("Starting a race %+v\n", *race)
	race.Stage = RaceStageOpen
	if race.setupMessage(s) != nil {
		race.EndRace()
		return
	}

	for _, snail := range race.Snails {
		snail.NewRace()
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

	// Betting Stage
	race.Render(s)
	if race.NoBets {
		time.Sleep(RaceNoBettingTimeout)
	} else {
		time.Sleep(RaceBettingTimeout)
	}
	race.Stage = RaceStageRunning

	// Race Stage
	race.Render(s)
	snailsFinished := 0
	for snailsFinished != len(race.Snails) {
		snailsFinished = 0
		for _, snail := range race.Snails {
			snail.Step()
			if snail.racePosition >= MaxRaceLength {
				snailsFinished++
			}
		}
		race.Render(s)
		log.Printf("Odds: %+v", race.Odds)
		log.Println("-------------------------")
		time.Sleep(1 * time.Second)
	}
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
		log.Warnf("failed to send race setup message: %s", err)
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
		body += fmt.Sprintf("- %s(<@%s>)\n", snail.Name, snail.Owner.DiscordID)
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
		"Bets are now open to everyone, do you feel lucky? Here are the entrants:\n\nRace ID: `%s`\n\n**Entrants: (%d/12)**\n",
		r.Id,
		len(r.Snails),
	)

	select_options := make([]discordgo.SelectMenuOption, 0)

	// Add the snails to the body as entrants `index - <oods> <snail_name>(<@owner_id>)`
	for index, snail := range r.Snails {
		if snail.Level == 0 {
			body += fmt.Sprintf("%2d - `%.02f` %s\n", index, r.Odds[index], snail.Name)
		} else {
			body += fmt.Sprintf("%2d - `%.02f` %s(<@%s>)\n", index, r.Odds[index], snail.Name, snail.Owner.DiscordID)
		}

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
		if snail.Level == 0 {
			body += fmt.Sprintf("%2d - `%.02f` %s\n", index, r.Odds[index], snail.Name)
		} else {
			body += fmt.Sprintf("%2d - `%.02f` %s(<@%s>)\n", index, r.Odds[index], snail.Name, snail.Owner.DiscordID)
		}
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
	track += "                        🏁\n"
	track += "  |-----------------------|\n"

	// Build snails
	for index, snail := range r.Snails {
		line := snail.renderPosition()
		track += fmt.Sprintf("%2d| %s | %s\n", index, line, snail.Name)

		if snail.Level == 0 {
			entrants += fmt.Sprintf("%2d - `%02f` %s\n", index, r.Odds[index], snail.Name)
		} else {
			entrants += fmt.Sprintf("%2d - `%02f` %s(<@%s>)\n", index, r.Odds[index], snail.Name, snail.Owner.DiscordID)
		}
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

}

func (r *Race) generateOdds() {
	r.Odds = make([]float64, len(r.Snails))

	sum_speed, sum_stamina := 0.0, 0.0
	for _, snail := range r.Snails {
		sum_speed += snail.Stats.Speed
		sum_stamina += snail.Stats.Stamina
	}

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
