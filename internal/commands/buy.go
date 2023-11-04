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

		for _, opt := range i.ApplicationCommandData().Options[0].Options {
			// establish if the user has the capital to purchase this snail (TODO)
			log.Info(opt.IntValue())
			canAfford, snailPrice := models.CanUserAffordSnail(*user, models.SnailStatLevel(opt.IntValue()))
			if !canAfford {
				log.WithField("cmd", "/buy").WithError(err).Infof("Snail costs %d, you only have %d", snailPrice, user.Money)
				ResponseEmbedFail(s, i, true,
					fmt.Sprintf("Snail costs %d, you only have %d", snailPrice, user.Money),
					"You first need to earn more money before buying this snail.\nTry buying a lower tier snail or trying your luck betting for more gold.",
				)
				return
			}

			// create and purchase the snail
			err := user.RemoveMoney(state.DB, uint64(snailPrice))
			if err != nil {
				log.WithField("cmd", "/buy").WithError(err).Infof("Something went wrong whilst trying to remove %dg from user %s wallet", snailPrice, i.Member.User.Username)
				ResponseEmbedFail(s, i, true,
					fmt.Sprintf("Something went wrong whilst trying to remove %dg from your wallet", snailPrice),
					"Please try purchasing the snail again",
				)
				return
			}

			snail, err := models.CreateSnail(state.DB, *user, models.SnailStatLevel(opt.IntValue()))
			if err != nil {
				log.WithField("cmd", "/buy").WithError(err).Infof("Something went wrong in the creation of a %s for %s", opt.Value.(string), i.Member.User.Username)
				ResponseEmbedFail(s, i, true,
					"Something went wrong attempting to purchase a snail",
					fmt.Sprintf("Please report the following, the snail stat level was %d for %s", opt.Value.(int), i.Member.User.Username),
				)
				return
			}

			// make the new snail their active snail
			erractiv := models.SetActiveSnail(state.DB, *user, *snail)
			if erractiv != nil {
				// if this fails, still continue on regardless, but let them know
				log.WithField("cmd", "/buy").WithError(err).Infof("Something went wrong setting a purchased snail %d as active", snail.ID)
				ResponseEmbedFail(s, i, true,
					"Something went wrong attempting to set the new snail as your active snail",
					"Please attempt manually setting your new snail as active using `/snailrace active [snail-id]`\n. This snail has been added to your backpack.",
				)
			}

			// Enable this if we're creating a backpack id once they have bought a snail
			// HandleNewSnail(state, s, i, snail, user)
			ResponseEmbedSuccess(s, i, false, "Congratulations on your New Snail",
				fmt.Sprintf("New Snail: **%s (lvl. %d)** with the following stats:\n```\n%s```\n", snail.Name, snail.Level, snail.Stats.RenderStatBlock(models.SnailStatLevel(snail.Tier))))
		}
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
