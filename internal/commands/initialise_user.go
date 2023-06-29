package commands

import (
	"fmt"

	"github.com/lcox74/snailrace/internal/models"

	"github.com/bwmarrin/discordgo"
	log "github.com/sirupsen/logrus"
)

// CommandPing is a simple ping command, this is used as a basic test to see if
// the bot it working correctly.
type CommandInitialise struct{}

func (c *CommandInitialise) Decleration() *discordgo.ApplicationCommandOption {
	return &discordgo.ApplicationCommandOption {
		Name:        "init",
		Description: "Initialise your account if you don't already have one",
		Type: discordgo.ApplicationCommandOptionSubCommand,
	}
}

func (c *CommandInitialise) Handler(state *models.State) func(s *discordgo.Session, i *discordgo.InteractionCreate) {
	return func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		log.Printf("[CMD] Init!\n")

		// Check if the user already has an account
		user, err := models.GetUserByDiscordID(state.DB, i.Member.User.ID)
		if err != nil {
			log.Printf("[CMD] Error getting user: %s\n", err)

			// Respond to the interaction with a message
			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: fmt.Sprintf("Error: %s", err),
				},
			})

			return
		}

		// Respond to the interaction with a message
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: fmt.Sprintf("User: %v", user),
			},
		})
	}
}
