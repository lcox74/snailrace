package commands

import (
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/lcox74/snailrace/internal/models"
	"github.com/lcox74/snailrace/pkg/styles"

	log "github.com/sirupsen/logrus"
)

type DiscordAppCommand interface {
	// The Decleration is the information that will be sent to Discord when
	// registering the command. This will show up as an Application Command in
	// the Discord Client with a supplied name and description.
	Decleration() *discordgo.ApplicationCommandOption

	// The Discord Application Handler is the function that will be called when
	// this command is triggered.
	AppHandler(state *models.State) func(s *discordgo.Session, i *discordgo.InteractionCreate)

	// The Discord Message Handler for component reactions
	ActionHandler(state *models.State, options ...string) map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate)

	// The Discord Modal Handler for modal sumbits
	ModalHandler(state *models.State, options ...string) map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate)
}

// RegisterCommand registers a command with Discord and adds a handler for the
// corresponding event on the disscord session.
func RegisterCommand(state *models.State, s *discordgo.Session, command DiscordAppCommand) error {

	decleration := command.Decleration()

	// Register a handler for the messageCreate events
	s.AddHandler(func(s *discordgo.Session, i *discordgo.InteractionCreate) {

		// Check if the interaction is a DM
		if i.Member == nil {
			styles.ErrDm(s, i.Interaction)
			return
		}

		switch i.Type {
		case discordgo.InteractionApplicationCommand:
			if i.ApplicationCommandData().Name != "snailrace" {
				return
			}

			if i.ApplicationCommandData().Options[0].Name == decleration.Name {
				log.WithField("cmd", i.ApplicationCommandData().Options[0].Name).Infof("User %s sent command", i.Member.User.Username)
				command.AppHandler(state)(s, i)
			}
		case discordgo.InteractionMessageComponent:

			breakDown := strings.Split(i.MessageComponentData().CustomID, ":")
			if len(breakDown) == 0 {
				if handler, ok := command.ActionHandler(state)[breakDown[0]]; ok {
					log.WithField("interaction", breakDown[0]).Infof("User %s sent interaction", i.Member.User.Username)
					handler(s, i)
				}
			} else {
				if handler, ok := command.ActionHandler(state, breakDown[1:]...)[breakDown[0]]; ok {
					log.WithField("interaction", breakDown[0]).Infof("User %s sent interaction", i.Member.User.Username)
					handler(s, i)
				}
			}

		case discordgo.InteractionModalSubmit:
			breakDown := strings.Split(i.ModalSubmitData().CustomID, ":")
			if len(breakDown) == 0 {
				if handler, ok := command.ModalHandler(state)[breakDown[0]]; ok {
					log.WithField("modal", breakDown[0]).Infof("User %s sent modal response", i.Member.User.Username)
					handler(s, i)
				}
			} else {
				if handler, ok := command.ModalHandler(state, breakDown[1:]...)[breakDown[0]]; ok {
					log.WithField("modal", breakDown[0]).Infof("User %s sent modal response", i.Member.User.Username)
					handler(s, i)
				}
			}
		}

	})

	return nil
}

func ResponseEmbed(s *discordgo.Session, i *discordgo.InteractionCreate, ephemeral bool, title string, color int, msg string) {
	flag := discordgo.MessageFlags(0)
	if ephemeral {
		flag = discordgo.MessageFlagsEphemeral
	}

	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Embeds: []*discordgo.MessageEmbed{
				{
					Title:       title,
					Color:       color,
					Description: msg,
				},
			},
			Flags: flag,
		},
	})
}

func ResponseEmbedSuccess(s *discordgo.Session, i *discordgo.InteractionCreate, ephemeral bool, title string, msg string) {
	ResponseEmbed(s, i, ephemeral, title, 0x2ecc71, msg)
}
func ResponseEmbedInfo(s *discordgo.Session, i *discordgo.InteractionCreate, ephemeral bool, title string, msg string) {
	ResponseEmbed(s, i, ephemeral, title, 0x3498db, msg)
}
func ResponseEmbedFail(s *discordgo.Session, i *discordgo.InteractionCreate, ephemeral bool, title string, msg string) {
	ResponseEmbed(s, i, ephemeral, title, 0xe74c3c, msg)
}
