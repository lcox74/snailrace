package main

import (
	"os"
	"os/signal"

	"github.com/lcox74/snailrace/internal"
	"github.com/lcox74/snailrace/internal/models"

	"github.com/joho/godotenv"
	log "github.com/sirupsen/logrus"
)

func main() {
	// Setup Logging
	setupLogging()

	// Load .env file
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// Bring up the database
	db, err := internal.SetupDatabase()
	if err != nil {
		log.Fatal("Error setting up database")
	}
	log.Printf("Database initialised and connected: %s\n", db.Name())

	// Create State
	state := models.NewState(db)

	// Initiliase Discord
	discord := internal.SetupDiscord(state)
	log.Printf("Discord Bot connected as %s\n", discord.State.User.Username)


	// Wait until CTRL-C or other term signal is received.
	log.Println("Snail Manager is now running.  Press CTRL-C to exit.")
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c
}

func setupLogging() {
	// Log as JSON instead of the default ASCII formatter.
	log.SetFormatter(&log.TextFormatter{})

	// Output to stdout instead of the default stderr
	log.SetOutput(os.Stdout)

	// Only log the info severity or above.
	log.SetLevel(log.InfoLevel)
}
