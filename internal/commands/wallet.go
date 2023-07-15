package commands

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
	"github.com/lcox74/snailrace/internal/models"
	"github.com/lcox74/snailrace/pkg/styles"
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

			styles.RespondInitialiseErr(s, i.Interaction, i.Member.Mention())
			return
		}

		// Display the users wallet
		styles.RespondOk(s, i.Interaction, true, "Wallet", fmt.Sprintf("💰 %dg", user.Money), nil)
	}
}

func (c *WalletCommand) ActionHandler(state *models.State, options ...string) map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate) {
	return map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate){}
}

func (c *WalletCommand) ModalHandler(state *models.State, options ...string) map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate) {
	return map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate){}
}
