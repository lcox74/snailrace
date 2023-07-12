package commands

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
	"github.com/lcox74/snailrace/internal/models"
	log "github.com/sirupsen/logrus"
)

// WalletCommand is a simple command that displays the users wallet.
type WalletCommand struct{}

func (c *WalletCommand) Decleration() *discordgo.ApplicationCommandOption {
	return &discordgo.ApplicationCommandOption{
		Name:        "wallet",
		Description: "Display you wallet.",
		Type:        discordgo.ApplicationCommandOptionSubCommand,
	}
}

func (c *WalletCommand) AppHandler(state *models.State) func(s *discordgo.Session, i *discordgo.InteractionCreate) {
	return func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		// Check if the user is initialised, if the user isn't initialised then
		// we need to tell them to initialise their account.
		user, err := models.GetUserByDiscordID(state.DB, i.Member.User.ID)
		if err != nil {
			log.WithField("cmd", "/wallet").WithError(err).Infof("User %s is not initialised", i.Member.User.Username)
			ResponseEmbedFail(s, i, true,
				fmt.Sprintf("I'm sorry %s, but you arent initialised", i.Member.User.Username),
				"You'll need to initialise your account with `/snailrace init` to use this command.",
			)
			return
		}

		// Display the users wallet
		ResponseEmbedSuccess(s, i, true, "Wallet", fmt.Sprintf("💰 %dg", user.Money))

	}
}

func (c *WalletCommand) ActionHandler(state *models.State, options ...string) map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate) {
	return map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate){}
}

func (c *WalletCommand) ModalHandler(state *models.State, options ...string) map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate) {
	return map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate){}
}
