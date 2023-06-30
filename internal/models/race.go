package models

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
	log "github.com/sirupsen/logrus"
)

const (
	RaceActionJoin = "host_join"
)

type Race struct {
	Id        string
	ChannelId string

	Host    *discordgo.User
	Message *discordgo.Message

	NoBets  bool
	NoFill  bool
	OnlyOne bool

	Snails []*Snail
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
	for _, s := range r.Snails {
		if s.ID == snail.ID {
			return fmt.Errorf("snail already exists")
		}
	}
	r.Snails = append(r.Snails, snail)
	return nil
}

func (r *Race) Render(s *discordgo.Session) {
	r.renderOpenRace(s)
}

func (r *Race) renderOpenRace(s *discordgo.Session) {
	var err error

	// Build the Embed Message
	title := "Race: Open"
	body := fmt.Sprintf(
		"A new race has been hosted by %s\n\nRace ID: `%s`\n\nTo join via command, enter the following:\n```\n/snailrace join %s\n```\n**Entrants: (%d/12)**\n",
		r.Host.Username,
		r.Id,
		r.Id,
		len(r.Snails),
	)

	// Add the snails to the body as entrants `- <snail_name>(<@owner_id>)`
	for _, snail := range r.Snails {
		body += fmt.Sprintf("- %s(<@%s>)\n", snail.Name, snail.Owner.DiscordID)
	}

	// Check if this is the first message in the Race state
	first_send := true
	if r.Message != nil {
		first_send = false
	}

	if first_send {
		// Send the first message in the Race state and store the message so we
		// can edit it later in the race state
		r.Message, err = s.ChannelMessageSendComplex(r.ChannelId, &discordgo.MessageSend{
			Embeds: []*discordgo.MessageEmbed{
				{
					Title:       title,
					Description: body,
					Color:       0x2ecc71,
				},
			},
			Components: []discordgo.MessageComponent{
				discordgo.ActionsRow{
					Components: []discordgo.MessageComponent{
						discordgo.Button{
							Label:    "Join",
							Style:    discordgo.SuccessButton,
							CustomID: fmt.Sprintf("%s:%s", RaceActionJoin, r.Id),
						},
					},
				},
			},
		})

		if err != nil {
			log.Println(err)
		}

	} else {
		// Edit the message to reflect the current state of the race, in this
		// sense it will mainly update the entrants
		edit := discordgo.NewMessageEdit(r.ChannelId, r.Message.ID)
		edit.Embeds = []*discordgo.MessageEmbed{
			{
				Title:       title,
				Description: body,
			},
		}
		s.ChannelMessageEditComplex(edit)
	}
}

func StartRace(s *discordgo.Session, race *Race) {
	race.Render(s)
}
