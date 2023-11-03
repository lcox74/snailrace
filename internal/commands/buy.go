package commands

import (
	"fmt"

	"github.com/lcox74/snailrace/internal/models"
	log "github.com/sirupsen/logrus"

	"github.com/bwmarrin/discordgo"
)

// CommandBuy shows the user their backpack of snails
type CommandBuy struct{}

func (c *CommandBuy) Decleration() *discordgo.ApplicationCommandOption {
	return &discordgo.ApplicationCommandOption{
		Name:        "buy",
		Description: "Buy a snail",
		Type:        discordgo.ApplicationCommandOptionSubCommand,
		Options: []*discordgo.ApplicationCommandOption{
			{
				Name:        "snail_type",
				Description: "The type of snail you want to buy",
				Type:        discordgo.ApplicationCommandOptionInteger,
				Required:    true,
				Choices: []*discordgo.ApplicationCommandOptionChoice{
					{
						Name:  "Starting",
						Value: 0,
					},
					{
						Name:  "Amateur",
						Value: 1,
					},
					{
						Name:  "Professional",
						Value: 2,
					},
					{
						Name:  "Expert",
						Value: 3,
					},
				},
			},
		},
	}
}

func (c *CommandBuy) AppHandler(state *models.State) func(s *discordgo.Session, i *discordgo.InteractionCreate) {
	return func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		user, err := models.GetUserByDiscordID(state.DB, i.Member.User.ID)
		if err != nil {
			log.WithField("cmd", "/buy").WithError(err).Infof("User %s is not initialised", i.Member.User.Username)
			ResponseEmbedFail(s, i, true,
				fmt.Sprintf("I'm sorry %s, but you arent initialised", user.DiscordID),
				"You'll need to initialise your account with `/snailrace init` to use this command.",
			)
			return
		}

		// establish if the user has the capital to purchase this snail (TODO)

		for _, opt := range i.ApplicationCommandData().Options[0].Options {
			snail, err := models.CreateSnail(state.DB, *user, opt.Value.(models.SnailStatLevel))
			if err != nil {
				log.WithField("cmd", "/buy").WithError(err).Infof("Something went wrong in the creation of a %s for %s", opt.Value.(string), i.Member.User.Username)
				ResponseEmbedFail(s, i, true,
					"Something went very wrong attempting to purchase a snail",
					fmt.Sprintf("Please report the following, the snail stat level was %d for %s", opt.Value.(int), i.Member.User.Username),
				)
				return
			}
		}

		// take the money away from the user (TODO)

		// Enable this if we're creating a backpack id once they have bought a snail
		// HandleNewSnail(state, s, i, snail, user)
		ResponseEmbedSuccess(s, i, true, "Snailing away", fmt.Sprintf("New Snail Name: %s", user.DiscordID))
	}
}

func HandleNewSnail(state *models.State, s *discordgo.Session, i *discordgo.InteractionCreate, snail *models.Snail, user *models.User) error {
	// is this their first new snail? If so, initialise them a new backpack id and all.
	allsnails, err := models.GetAllSnails(state.DB, *user)
	if err != nil {
		log.WithField("cmd", "/buy").WithError(err).Infof("Couldn't retrieve all snails for %s", user.DiscordID)
		ResponseEmbedFail(s, i, true,
			fmt.Sprintf("Snails could not be retrieved for this user %s", user.DiscordID),
			"Please try again and report the error",
		)
	}
	// numberSnails := len(allsnails)
	// if (numberSnails )
	log.Info(allsnails)

	return nil
}

func (c *CommandBuy) ActionHandler(state *models.State, options ...string) map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate) {
	return map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate){}
}

func (c *CommandBuy) ModalHandler(state *models.State, options ...string) map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate) {
	return map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate){}
}
