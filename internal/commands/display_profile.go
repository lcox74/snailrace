package commands

import (
	"fmt"
	"math"
	"strings"

	"github.com/lcox74/snailrace/internal/models"
	"golang.org/x/text/language"
	"golang.org/x/text/message"

	"github.com/bwmarrin/discordgo"
	log "github.com/sirupsen/logrus"
)

type CommandDisplayProfile struct{}

func (c *CommandDisplayProfile) Decleration() *discordgo.ApplicationCommandOption {
	return &discordgo.ApplicationCommandOption{
		Name:        "profile",
		Description: "Check our your user profile",
		Type:        discordgo.ApplicationCommandOptionSubCommand,
		Options: []*discordgo.ApplicationCommandOption{
			{
				Name:        "discor",
				Description: "The user to look up",
				Type:        discordgo.ApplicationCommandOptionString,
				Required:    false,
			},
		},
	}
}

func (c *CommandDisplayProfile) AppHandler(state *models.State) func(s *discordgo.Session, i *discordgo.InteractionCreate) {
	return func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		// Check if the user is initialised, if the user isn't initialised then
		// we need to tell them to initialise their account.
		user, err := models.GetUserByDiscordID(state.DB, i.Member.User.ID)
		if err != nil {
			log.WithField("cmd", "/display").WithError(err).Infof("User %s is not initialised", i.Member.User.Username)
			ResponseEmbedFail(s, i, true,
				fmt.Sprintf("I'm sorry %s, but you arent initialised", i.Member.User.Username),
				"You'll need to initialise your account with `/snailrace init` to use this command.",
			)
			return
		}

		// Display the user profile
		p := message.NewPrinter(language.English)

		// get some information about the user
		winRate := user.Wins * 100 / user.Races
		levelProgress := models.GetPercentageLevelProgress(state.DB, user)
		progressBar := GenerateProgressBar(levelProgress)

		allSnails, err := models.GetAllSnails(state.DB, *user)
		activeSnail, err2 := models.GetActiveSnail(state.DB, *user)

		if err != nil || err2 != nil {
			log.WithField("cmd", "/display").WithError(err).Infof("Could not find snails for user %s", i.Member.User.Username)
			ResponseEmbedFail(s, i, true,
				fmt.Sprintf("Error, could not find snails for user %s", i.Member.User.Username),
				"This shouldn't happen. Please report this error to someone",
			)
			return
		}

		ResponseEmbedSuccess(s, i, true, "Profile",
			p.Sprintf("**Username**: %s\n\n**Level**: %d\n**Progress**: %s\n\n**Win Rate**: %d%%\n**Races**: %d\n**Total Snails**: %d\n\n🐌 %s\n💰 %dg",
				i.Member.User.Username, user.Level, progressBar, winRate, user.Races, len(allSnails), activeSnail.Name, user.Money))
	}
}

func GenerateProgressBar(percentage float64) string {
	numSquares := int(math.Floor(percentage / 10))
	progress := strings.Repeat("🟩", numSquares)
	progress += strings.Repeat("⬛", 10-numSquares)
	return progress
}

func (c *CommandDisplayProfile) ActionHandler(state *models.State, options ...string) map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate) {
	return map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate){}
}

func (c *CommandDisplayProfile) ModalHandler(state *models.State, options ...string) map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate) {
	return map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate){}
}
