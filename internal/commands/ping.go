package commands

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
	log "github.com/sirupsen/logrus"
)

// CommandPing is a simple ping command, this is used as a basic test to see if
// the bot it working correctly.
type CommandPing struct{}

func (c *CommandPing) Decleration() discordgo.ApplicationCommand {
	return discordgo.ApplicationCommand{
		Name:        "ping",
		Description: "Ping the bot, is it alive?",
	}
}

func (c *CommandPing) Handler() func(s *discordgo.Session, i *discordgo.InteractionCreate) {
	return func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		response := fmt.Sprintf("Pong <@%s>!", i.Member.User.ID)
		log.Printf("[CMD] Ping! -> %s\n", response)

		// Respond to the interaction with a message
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: response,
			},
		})
	}
}
