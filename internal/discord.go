package internal

import (
	"os"

	"github.com/lcox74/snailrace/internal/commands"
	"github.com/lcox74/snailrace/internal/models"

	"github.com/bwmarrin/discordgo"
	log "github.com/sirupsen/logrus"
)

const (
	DiscordTokenEnv       = "DISCORD_TOKEN"
	DiscordGameStatus     = "Snail Manager"
	DiscordCmdPrefix      = "snailrace"
	DiscordCmdDescription = "Snailrace Commands"
)

func SetupDiscord(state *models.State) *discordgo.Session {
	// Check for Discord Token
	if os.Getenv(DiscordTokenEnv) == "" {
		log.Fatal("Unable to find Discord Token in environment variables")
	}

	// Create Discord Session
	discord, err := discordgo.New("Bot " + os.Getenv(DiscordTokenEnv))
	if err != nil {
		log.WithError(err).Fatal("Failed creating Discord session:", err)
	}

	// Regiser a handler for the ready event
	discord.AddHandler(func(s *discordgo.Session, event *discordgo.Ready) {
		// Set the playing status.
		s.UpdateGameStatus(0, DiscordGameStatus)

		// Log ready
		log.Infof("Bot is ready! (User: %s)\n", event.User.Username)
	})

	// Open a websocket connection to Discord and begin listening.
	err = discord.Open()
	if err != nil {
		log.WithError(err).Fatal("Failed opening connection to Discord:", err)
	}

	// Register Commands
	err = RegisterCommands(state, discord)
	if err != nil {
		log.WithError(err).Fatal("Failed registering commands:", err)
	}

	return discord
}

func RegisterCommands(state *models.State, s *discordgo.Session) error {
	// Commands to register
	cmds := []commands.DiscordAppCommand{
		&commands.CommandPing{},
		&commands.CommandInitialise{},
		&commands.CommandHostRace{},
		&commands.CommandJoinRace{},
		&commands.BetCommand{},
	}

	// Create Full decleration
	decleration := &discordgo.ApplicationCommand{
		Name:        DiscordCmdPrefix,
		Description: DiscordCmdDescription,
		Options:     []*discordgo.ApplicationCommandOption{},
	}

	// Register all the commands as handlers
	for _, cmd := range cmds {
		err := commands.RegisterCommand(state, s, cmd)
		if err != nil {
			log.WithError(err).Warnf("Failed registering command: %s", err)
			return err
		}

		// Add the command to the decleration as a subcommand
		decleration.Options = append(decleration.Options, cmd.Decleration())
	}

	// Register the Application commands decleration with Discord
	_, err := s.ApplicationCommandCreate(s.State.User.ID, "", decleration)
	return err
}
