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
		log.WithError(err).Fatal("Error loading .env file")
		return
	}

	// Bring up the database
	db, err := internal.SetupDatabase()
	if err != nil {
		log.WithError(err).Fatal("Error setting up database")
		return
	}
	log.Infoln("Database initialised and connected.")

	// Create State
	state := models.NewState(db)

	// Initiliase Discord
	discord := internal.SetupDiscord(state)
	if discord == nil {
		log.Fatal("Unable to setup Discord")
		return
	}

	// Wait until CTRL-C or other term signal is received.
	log.Infoln("Snail racer is now running.")
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
