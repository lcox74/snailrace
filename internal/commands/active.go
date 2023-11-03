package commands

import (
	"fmt"

	"github.com/lcox74/snailrace/internal/models"
	log "github.com/sirupsen/logrus"

	"github.com/bwmarrin/discordgo"
)

// CommandActive makes a snail you own one your active snails
type CommandActive struct{}

func (c *CommandActive) Decleration() *discordgo.ApplicationCommandOption {
	return &discordgo.ApplicationCommandOption{
		Name:        "active",
		Description: "Sets your active snail",
		Type:        discordgo.ApplicationCommandOptionSubCommand,
		Options: []*discordgo.ApplicationCommandOption{
			{
				Name:        "active-snail",
				Description: "The ID of your own snail you want to use",
				Type:        discordgo.ApplicationCommandOptionString,
				Required:    true,
			},
		},
	}
}

func (c *CommandActive) AppHandler(state *models.State) func(s *discordgo.Session, i *discordgo.InteractionCreate) {
	return func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		user, err := models.GetUserByDiscordID(state.DB, i.Member.User.ID)
		if err != nil {
			log.WithField("cmd", "/active").WithError(err).Infof("No record for user %s", i.Member.User.Username)

			ResponseEmbedFail(s, i, true,
				fmt.Sprintf("I'm sorry %s, but you arent initialised", i.Member.User.Username),
				"You'll need to initialise your account with `/snailrace init` to use this command.",
			)
			return
		}

		for _, opt := range i.ApplicationCommandData().Options[0].Options {
			switch opt.Name {
			case "active-snail":
				snail, err := models.GetSnailFromID(state.DB, *user, opt.Value.(string))
				if err != nil {
					log.WithField("cmd", "/active").WithError(err).Infof("Could not find snail %s for owner %s", opt.Value.(string), user.DiscordID)
					ResponseEmbedFail(s, i, true,
						fmt.Sprintf("Invalid / could not find snail: ID %s", opt.Value.(string)),
						"Could not find this snail. Are you sure it's yours?\nCheck `/snailrace backpack` to double check.",
					)
					return
				}

				if activeErr := models.SetActiveSnail(state.DB, *user, *snail); activeErr != nil {
					log.WithField("cmd", "/active").WithError(err).Infof("There was a problem setting the active snail: %s", opt.Value.(string))
					ResponseEmbedFail(s, i, true,
						fmt.Sprintf("There was a problem setting the active snail %s", opt.Value.(string)),
						"Something went wrong, try again and report the error",
					)
					return
				}

				ResponseEmbedSuccess(s, i, true, "New Active Snail Set",
					fmt.Sprintf("Your active snail was set successfully.\nYour new active snail is: `%s`", snail.Name),
				)

			default:
				return
			}
		}

	}
}

func (c *CommandActive) ActionHandler(state *models.State, options ...string) map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate) {
	return map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate){}
}

func (c *CommandActive) ModalHandler(state *models.State, options ...string) map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate) {
	return map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate){}
}
