package commands

import (
	"fmt"
	"math"
	"strings"

	"github.com/lcox74/snailrace/internal/models"

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
				Name:        "user-option",
				Description: "The user to look up",
				Type:        discordgo.ApplicationCommandOptionUser,
				Required:    false,
			},
		},
	}
}

func GetRequestedUser(s *discordgo.Session, i *discordgo.InteractionCreate) (*discordgo.User, bool, error) {
	if len(i.ApplicationCommandData().Options) > 0 {
		for _, opt := range i.ApplicationCommandData().Options[0].Options {
			switch opt.Name {
			case "user-option":
				usr, err := s.User(opt.Value.(string))
				return usr, false, err
			default:
				// if all else fails
				return i.Member.User, true, nil
			}
		}
	}

	return i.Member.User, true, nil
}

func (c *CommandDisplayProfile) AppHandler(state *models.State) func(s *discordgo.Session, i *discordgo.InteractionCreate) {
	return func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		// get the discord user if requested, otherwise use the current one
		discorduser, personal, err := GetRequestedUser(s, i)
		if err != nil {
			log.WithField("cmd", "/display").WithError(err).Infof("Something went wrong getting the requested user %s", discorduser.Username)
			ResponseEmbedFail(s, i, true,
				fmt.Sprintf("Something went wrong getting the requested user %s", discorduser.Username),
				"Please try again, and report the error if it continues",
			)
			return
		}

		// Check if the user is initialised, if the user isn't initialised then
		// we need to tell them to initialise their account.
		user, err := models.GetUserByDiscordID(state.DB, discorduser.ID)
		if err != nil {
			log.WithField("cmd", "/display").WithError(err).Infof("User %s is not initialised", discorduser.Username)
			if personal {
				ResponseEmbedFail(s, i, true,
					fmt.Sprintf("I'm sorry %s, but you arent initialised", discorduser.Username),
					"You'll need to initialise your account with `/snailrace init` to use this command.",
				)
			} else {
				ResponseEmbedFail(s, i, true,
					fmt.Sprintf("I'm sorry, but %s is not initialised", discorduser.Username),
					"They will need to initialise an account with `/snailrace init`",
				)
			}
			return
		}

		// get some information about the user
		winRate := user.Wins * 100 / user.Races
		levelProgress := models.GetPercentageLevelProgress(state.DB, user)
		progressBar := GenerateProgressBar(levelProgress)

		allSnails, err := models.GetAllSnails(state.DB, *user)
		activeSnail, err2 := models.GetActiveSnail(state.DB, *user)

		if err != nil || err2 != nil {
			log.WithField("cmd", "/display").WithError(err).Infof("Could not find snails for user %s", discorduser.Username)
			ResponseEmbedFail(s, i, true,
				fmt.Sprintf("Error, could not find snails for user %s", discorduser.Username),
				"This shouldn't happen. Please report this error to someone",
			)
			return
		}

		ResponseEmbedSuccess(s, i, personal, "Profile",
			fmt.Sprintf("**Username**: %s\n\n**Level**: %d\n**Progress**: %s\n\n**Win Rate**: %d%%\n**Races**: %d\n**Total Snails**: %d\n\n🐌 %s\n💰 %dg",
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
