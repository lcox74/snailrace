package commands

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
	"github.com/lcox74/snailrace/internal/models"
	log "github.com/sirupsen/logrus"
)

type BetCommand struct{}

func (c *BetCommand) Decleration() *discordgo.ApplicationCommandOption {
	return &discordgo.ApplicationCommandOption{
		Name:        "bet",
		Description: "So you want to put your money where your mouth is?",
		Type:        discordgo.ApplicationCommandOptionSubCommand,
		Options: []*discordgo.ApplicationCommandOption{
			{
				Name:        "race_id",
				Description: "The race id you want to bet on.",
				Type:        discordgo.ApplicationCommandOptionString,
				Required:    true,
			},
			{
				Name:        "snail_index",
				Description: "The index of the snail in the race.",
				Type:        discordgo.ApplicationCommandOptionInteger,
				Required:    true,
			},
			{
				Name:        "amount",
				Description: "The amount of money you want to bet.",
				Type:        discordgo.ApplicationCommandOptionInteger,
				Required:    true,
			},
		},
	}
}

func (c *BetCommand) AppHandler(state *models.State) func(s *discordgo.Session, i *discordgo.InteractionCreate) {
	return func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		log.Printf("[CMD] Bet!\n")

		// Check if the user is initialised, if the user isn't initialised then
		// we need to tell them to initialise their account.
		user, err := models.GetUserByDiscordID(state.DB, i.Member.User.ID)
		if err != nil {
			ResponseEmbedFail(s, i, true,
				fmt.Sprintf("I'm sorry %s, but you arent initialised", i.Member.User.Username),
				"You'll need to initialise your account with `/snailrace init` to use this command.",
			)
			return
		}

		// Pull the options from the interaction
		raceId := ""
		snailIndex := 0
		amount := 0
		for _, option := range i.ApplicationCommandData().Options[0].Options {
			switch option.Name {
			case "race_id":
				raceId = option.StringValue()
			case "snail_index":
				snailIndex = int(option.IntValue())
			case "amount":
				amount = int(option.IntValue())
			}
		}

		// Check if the race exists, if it doesn't then we need to tell the
		// user
		race, ok := state.Races[raceId]
		if !ok {
			ResponseEmbedFail(s, i, true, fmt.Sprintf("Race %s not avaliable", raceId), "There is currently no race with the ID you supplied.")
			return
		}

		// Check if the snail exists, if it doesn't then we need to tell the
		// user
		snail := race.GetSnail(snailIndex)
		if snail == nil {
			ResponseEmbedFail(s, i, true, fmt.Sprintf("Invalid snail to bet for race %s", raceId), "There is currently no snail with the ID you supplied.")
			return
		}

		// Check if the user has enough money to make the bet
		if int(user.Money) < amount {
			log.Warnf("Player %s tried to bet %d g but only has %d g", i.Member.User.Username, amount, user.Money)
			ResponseEmbedFail(s, i, true, fmt.Sprintf("Sorry %s but you can't afford the bet", i.Member.User.Username), fmt.Sprintf("You don't have enough money to place that bet, you only have %d g.", user.Money))
			return
		}

		// Place the bet and remove the money from the user
		switch race.PlaceBet(snailIndex, amount, user.DiscordID) {
		case models.ErrInvalidSnail:
			ResponseEmbedFail(s, i, true, fmt.Sprintf("Sorry %s that snail doesn't exist", i.Member.User.Username), "The snail you have selected to bet is invalid, the snail isn't in the race.")
			return
		case models.ErrBetsClosed:
			ResponseEmbedFail(s, i, true, fmt.Sprintf("Sorry %s Bets are Closed", i.Member.User.Username), "Bet's are closed so we can't accept your bet.")
			return
		}
		ResponseEmbedSuccess(s, i, true, fmt.Sprintf("Bet placed for %s", snail.Name), fmt.Sprintf("You've placed a bet for %s of %d g", snail.Name, amount))
		user.RemoveMoney(state.DB, uint64(amount))

	}
}

func (c *BetCommand) ActionHandler(state *models.State, options ...string) map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate) {
	return map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate){}
}
func (c *BetCommand) ModalHandler(state *models.State, options ...string) map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate) {
	return map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate){}
}
