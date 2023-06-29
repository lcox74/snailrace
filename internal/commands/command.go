package commands

import (
	"log"

	"github.com/bwmarrin/discordgo"
	"github.com/lcox74/snailrace/internal/models"
)

type DiscordAppCommand interface {
	// The Decleration is the information that will be sent to Discord when
	// registering the command. This will show up as an Application Command in
	// the Discord Client with a supplied name and description.
	Decleration() *discordgo.ApplicationCommandOption

	// The Handler is the function that will be called when this command is
	// triggered.
	Handler(state *models.State) func(s *discordgo.Session, i *discordgo.InteractionCreate)
}

// RegisterCommand registers a command with Discord and adds a handler for the
// corresponding event on the disscord session.
func RegisterCommand(state *models.State, s *discordgo.Session, command DiscordAppCommand) error {

	decleration := command.Decleration()
	
	// Register a handler for the messageCreate events
	s.AddHandler(func (s *discordgo.Session, i *discordgo.InteractionCreate) {
		if i.ApplicationCommandData().Name != "snailrace" {
			return
		}

		log.Printf("%s used %s\n", i.Member.User.ID, i.ApplicationCommandData().Options[0].Name)
		if i.ApplicationCommandData().Options[0].Name == decleration.Name {
			command.Handler(state)(s, i)
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
					Title: title,
					Color: color,
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