package internal

import (
	"os"

	"github.com/lcox74/snailrace/internal/commands"

	"github.com/bwmarrin/discordgo"
	log "github.com/sirupsen/logrus"
)

const (
	DiscordTokenEnv   = "DISCORD_TOKEN"
	DiscordGameStatus = "Snail Manager"
)

func SetupDiscord() *discordgo.Session {
	// Check for Discord Token
	if os.Getenv(DiscordTokenEnv) == "" {
		log.Fatal("Unable to find Discord Token in environment variables")
	}

	// Create Discord Session
	discord, err := discordgo.New("Bot " + os.Getenv(DiscordTokenEnv))
	if err != nil {
		log.Fatal("Failed creating Discord session:", err)
	}

	// Regiser a handler for the ready event
	discord.AddHandler(func(s *discordgo.Session, event *discordgo.Ready) {
		// Set the playing status.
		s.UpdateGameStatus(0, DiscordGameStatus)

		// Log ready
		log.Printf("Bot is ready! (User: %s)\n", event.User.Username)
	})

	// Open a websocket connection to Discord and begin listening.
	err = discord.Open()
	if err != nil {
		log.Fatal("Failed opening connection to Discord:", err)
	}

	// Regiser Commands
	err = commands.RegisterCommand(discord, &commands.CommandPing{})
	if err != nil {
		log.Fatal("Failed registering command:", err)
	}

	return discord
}
