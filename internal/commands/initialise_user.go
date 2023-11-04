package commands

import (
	"fmt"

	"github.com/lcox74/snailrace/internal/models"
	"gorm.io/gorm"

	"github.com/bwmarrin/discordgo"
	log "github.com/sirupsen/logrus"
)

// CommandInitialise initialise a user, this is used to create a new user and
// snail if they don't already have one.
type CommandInitialise struct{}

func (c *CommandInitialise) Decleration() *discordgo.ApplicationCommandOption {
	return &discordgo.ApplicationCommandOption{
		Name:        "init",
		Description: "Initialise your account if you don't already have one",
		Type:        discordgo.ApplicationCommandOptionSubCommand,
	}
}

func (c *CommandInitialise) AppHandler(state *models.State) func(s *discordgo.Session, i *discordgo.InteractionCreate) {
	return func(s *discordgo.Session, i *discordgo.InteractionCreate) {

		// Check if the user already has an account
		user, err := models.GetUserByDiscordID(state.DB, i.Member.User.ID)
		if err != nil && err != gorm.ErrRecordNotFound {
			log.WithField("cmd", "/init").WithError(err).Warnf("Error getting user %s", i.Member.User.Username)
			c.respondWithFail(s, i)
			return
		}

		// Check if the user doesn't exist, if it doesn't exist we want to
		// create it and then create a snail for them.
		if err == gorm.ErrRecordNotFound {
			log.WithField("cmd", "/init").Infof("Creating record for user %s", i.Member.User.Username)
			c.respondCreateNew(s, i, state.DB)
			return
		}

		// User already exists, lets just remind them of their snail
		log.WithField("cmd", "/init").Infof("Existing record for user %s", i.Member.User.Username)
		c.respondExisting(s, i, state.DB, user)
	}
}

func (c *CommandInitialise) ActionHandler(state *models.State, options ...string) map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate) {
	return map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate){}
}

func (c *CommandInitialise) ModalHandler(state *models.State, options ...string) map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate) {
	return map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate){}
}

func (c CommandInitialise) respondCreateNew(s *discordgo.Session, i *discordgo.InteractionCreate, db *gorm.DB) {
	// Create a new user
	user, err := models.CreateUser(db, i.Member.User.ID)
	if err != nil {
		log.WithField("cmd", "/init").WithError(err).Warnf("Error creating user %s", i.Member.User.Username)
		c.respondWithFail(s, i)
		return
	}

	// Create a new snail
	snail, err := models.CreateSnail(db, *user, models.StartingSnail)
	if err != nil {
		log.WithField("cmd", "/init").WithError(err).Warnf("Error creating snail for user %s", i.Member.User.Username)
		c.respondWithFail(s, i)
		return
	}
	models.SetActiveSnail(db, *user, *snail)

	// Notify the user that they have been created
	ResponseEmbedSuccess(s, i, false,
		fmt.Sprintf("Welcome to Snailrace %s!", i.Member.User.Username),
		fmt.Sprintf("Your snail is called **%s (lvl. %d)** and has the following stats:\n```\n%s```\n", snail.Name, snail.Level, snail.Stats.RenderStatBlock(models.SnailStatLevel(snail.Tier))),
	)
}

func (c CommandInitialise) respondExisting(s *discordgo.Session, i *discordgo.InteractionCreate, db *gorm.DB, user *models.User) {
	// Get the user's active snail
	snail, err := models.GetActiveSnail(db, *user)
	if err != nil && err != gorm.ErrRecordNotFound {
		log.WithField("cmd", "/init").WithError(err).Warnf("Error getting active snail for user %s", i.Member.User.Username)
		c.respondWithFail(s, i)
		return
	}

	// If the user doesn't have an active snail, check if they have any snails
	if err != gorm.ErrRecordNotFound {
		snails, err := models.GetAllSnails(db, *user)
		if err != nil && err != gorm.ErrRecordNotFound {
			log.WithField("cmd", "/init").WithError(err).Warnf("Error getting all snails for user %s", i.Member.User.Username)
			c.respondWithFail(s, i)
			return
		}

		// If the user has no snails, we create a new one
		if err == gorm.ErrRecordNotFound || len(snails) == 0 {
			// We create a new snail for the user
			snail, err := models.CreateSnail(db, *user, models.StartingSnail)
			if err != nil {
				log.WithField("cmd", "/init").WithError(err).Warnf("Error creating snail for user %s", i.Member.User.Username)
				c.respondWithFail(s, i)
				return
			}
			models.SetActiveSnail(db, *user, *snail)

			// Notify the user that they have been created
			ResponseEmbedSuccess(s, i, false,
				fmt.Sprintf("Welcome to Snailrace %s!", i.Member.User.Username),
				fmt.Sprintf("For some reason you had no snails, your snail is called **%s (lvl. %d)** and has the following stats:\n```\n%s```\n", snail.Name, snail.Level, snail.Stats.RenderStatBlock(models.SnailStatLevel(snail.Tier))),
			)
			return
		}

		// We set the first snail as the active snail
		models.SetActiveSnail(db, *user, snails[0])
	}

	// Respond to the interaction with a message
	ResponseEmbedInfo(s, i, false,
		fmt.Sprintf("You are already initialised  %s!", i.Member.User.Username),
		fmt.Sprintf("Your snail currently active snail is **%s (lvl. %d)** with the following stats:\n```\n%s```\n", snail.Name, snail.Level, snail.Stats.RenderStatBlock(models.SnailStatLevel(snail.Tier))),
	)
}

func (c CommandInitialise) respondWithFail(s *discordgo.Session, i *discordgo.InteractionCreate) {
	ResponseEmbedFail(s, i, false,
		fmt.Sprintf("I'm sorry %s, but there has been an issue", i.Member.User.Username),
		"There has been an issue with initialising your account. Please try again later.",
	)
}
