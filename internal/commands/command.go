package commands

import (
	"github.com/bwmarrin/discordgo"
)

type DiscordAppCommand interface {
	// The Decleration is the information that will be sent to Discord when
	// registering the command. This will show up as an Application Command in
	// the Discord Client with a supplied name and description.
	Decleration() discordgo.ApplicationCommand

	// The Handler is the function that will be called when this command is
	// triggered.
	Handler() func(s *discordgo.Session, i *discordgo.InteractionCreate)
}

// RegisterCommand registers a command with Discord and adds a handler for the
// corresponding event on the disscord session.
func RegisterCommand(s *discordgo.Session, command DiscordAppCommand) error {
	// Register a handler for the messageCreate events
	s.AddHandler(command.Handler())

	// Register the command as an application command
	decleration := command.Decleration()
	_, err := s.ApplicationCommandCreate(s.State.User.ID, "", &decleration)
	return err
}
