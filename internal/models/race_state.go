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
	ErrRaceClosed    = fmt.Errorf("race is closed")
	ErrAlreadyJoined = fmt.Errorf("snail already joined")
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

	NoBets  bool
	NoFill  bool
	OnlyOne bool

	Snails []*Snail
	Bets   []RaceBet
}

func (r *Race) SetupNewRace(id string, channelId string, host *discordgo.User, endRace func()) {
	r.Id = id
	r.ChannelId = channelId
	r.Host = host
	r.EndRace = endRace
	r.Stage = RaceStageOpen
	r.Snails = make([]*Snail, 0)
}

func (r *Race) SetNoBets() {
	r.NoBets = true
}
func (r *Race) SetNoFill() {
	r.NoFill = true
}
func (r *Race) SetOnlyOne() {
	r.OnlyOne = true
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
func (r *Race) PlaceBet(index int, amount int, userDiscordId string) {
	if index < 0 || index >= len(r.Snails) {
		return
	}

	r.Bets = append(r.Bets, RaceBet{
		UserDiscordId: userDiscordId,
		Amount:        amount,
		SnailIndex:    index,
	})
}

func StartRace(s *discordgo.Session, race *Race) {
	log.Printf("Starting a race %+v\n", *race)
	race.Stage = RaceStageOpen
	if race.setupMessage(s) != nil {
		race.EndRace()
		return
	}

	// Race Open Stage
	race.Render(s)
	time.Sleep(10 * time.Second)
	race.Stage = RaceStageBetting

	// Race Betting Stage
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
	// Build the Embed Message
	title := "Race: Bets are Open"
	body := fmt.Sprintf(
		"Bets are now open to everyone, do you feel lucky? Here are the entrants:\n\nRace ID: `%s`\n\n**Entrants: (%d/12)**\n",
		r.Id,
		len(r.Snails),
	)

	select_options := make([]discordgo.SelectMenuOption, 0)

	// Add the snails to the body as entrants `index - <snail_name>(<@owner_id>)`
	for index, snail := range r.Snails {
		body += fmt.Sprintf("%2d - %s(<@%s>)\n", index, snail.Name, snail.Owner.DiscordID)

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
func (r *Race) renderRunning(s *discordgo.Session) {

}
func (r *Race) renderFinished(s *discordgo.Session) {

}
