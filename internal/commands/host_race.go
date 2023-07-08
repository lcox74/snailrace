package commands

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/lcox74/snailrace/internal/models"

	"github.com/bwmarrin/discordgo"
	log "github.com/sirupsen/logrus"
)

type CommandHostRace struct{}

func (c *CommandHostRace) Decleration() *discordgo.ApplicationCommandOption {
	return &discordgo.ApplicationCommandOption{
		Name:        "host",
		Description: "Let's host a race",
		Type:        discordgo.ApplicationCommandOptionSubCommand,
		Options: []*discordgo.ApplicationCommandOption{
			{
				Name:        "no-bets",
				Description: "This flag skips the ability to place bets.",
				Type:        discordgo.ApplicationCommandOptionBoolean,
			},
			{
				Name:        "dont-fill",
				Description: "If this is set, then there wont be any additional snails added if the race has less than 4 snails",
				Type:        discordgo.ApplicationCommandOptionBoolean,
			},
			{
				Name:        "only-one",
				Description: "Continue racing until there is only one snail left. No Ties.",
				Type:        discordgo.ApplicationCommandOptionBoolean,
			},
		},
	}
}

func (c *CommandHostRace) AppHandler(state *models.State) func(s *discordgo.Session, i *discordgo.InteractionCreate) {
	return func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		// Check if the user is initialised, if the user isn't initialised then
		// we need to tell them to initialise their account.
		user, err := models.GetUserByDiscordID(state.DB, i.Member.User.ID)
		if err != nil {
			log.WithField("cmd", "/host").WithError(err).Infof("No record for user %s", i.Member.User.Username)

			ResponseEmbedFail(s, i, true,
				fmt.Sprintf("I'm sorry %s, but you arent initialised", i.Member.User.Username),
				"You'll need to initialise your account with `/snailrace init` to use this command.",
			)
			return
		}

		// We need to get the active snail of the host to automatically add them
		// to the race
		snail, err := models.GetActiveSnail(state.DB, *user)
		if err != nil {
			log.WithField("cmd", "/host").WithError(err).Warnf("Error getting active snail for %s", i.Member.User.Username)

			ResponseEmbedFail(s, i, true,
				fmt.Sprintf("I'm sorry %s, but we couldn't get your active snail", i.Member.User.Username),
				"There has been an issue with the action you sent, please try again.",
			)
			return
		}

		// Generate the race and add the host as the first snail
		race := state.NewRace(s, i.ChannelID, i.Member.User)
		race.AddSnail(snail)

		// Add flags to the Race
		if len(i.ApplicationCommandData().Options) > 0 {
			for _, opt := range i.ApplicationCommandData().Options[0].Options {
				switch opt.Name {
				case "no-bets":
					race.SetNoBets()
				case "dont-fill":
					race.SetDontFill()
				case "only-one":
					race.SetOnlyOne()
				}
			}
		}

		// Start the race as a seperate process
		go models.StartRace(s, race)

		// Respond to the interaction with a message
		ResponseEmbedSuccess(s, i, true,
			fmt.Sprintf("You just hosted a race %s!", i.Member.User.Username),
			"Your snail is officially waiting at the starting line for other snails to join.",
		)
	}
}

func (c *CommandHostRace) ActionHandler(state *models.State, options ...string) map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate) {
	return map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate){
		models.RaceActionJoin: func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			// The Join Action acts as the command /snailrace join <race_id>
			// If the caller doesn't supply the `race_id` then we need to
			// through and error, theoretically this should nevery error
			if len(options) != 1 {
				log.WithField("interaction", models.RaceActionJoin).WithError(errors.New("invalid options")).Errorf("Not enough arguments/options from user %s", i.Member.User.Username)
				ResponseEmbedFail(s, i, true,
					fmt.Sprintf("I'm sorry %s, but there has been an issue", i.Member.User.Username),
					"There has been an issue with the action you sent, please try again.",
				)
				return
			}

			// Check if the user is initialised, if the user isn't initialised then
			// we need to tell them to initialise their account.
			user, err := models.GetUserByDiscordID(state.DB, i.Member.User.ID)
			if err != nil {
				log.WithField("interaction", models.RaceActionJoin).WithError(err).Infof("Error getting record for user %s", i.Member.User.Username)
				ResponseEmbedFail(s, i, true,
					fmt.Sprintf("I'm sorry %s, but you arent initialised", i.Member.User.Username),
					"You'll need to initialise your account with `/snailrace init` to use this command.",
				)
				return
			}

			// We neet to get the user's active snail to add to the race
			snail, err := models.GetActiveSnail(state.DB, *user)
			if err != nil {
				log.WithField("interaction", models.RaceActionJoin).WithError(err).Warnf("Error getting active snail for user %s", i.Member.User.Username)
				ResponseEmbedFail(s, i, true,
					fmt.Sprintf("I'm sorry %s, but we couldn't get your active snail", i.Member.User.Username),
					"There has been an issue with the action you sent, please try again.",
				)
				return
			}

			// Check if the race exists, if it doesn't then we need to tell the
			// user
			raceId := options[0]
			race, ok := state.Races[raceId]
			if !ok {
				log.WithField("interaction", models.RaceActionJoin).WithError(errors.New("no existing race")).Infof("The raceid %s is not active, requested by user %s", raceId, i.Member.User.Username)
				ResponseEmbedFail(s, i, true, fmt.Sprintf("Race %s not avaliable", raceId), "There is currently no race with the ID you supplied.")
				return
			}

			err = race.AddSnail(snail)
			if err != nil {
				log.WithField("interaction", models.RaceActionJoin).WithError(err).Infof("The user %s is already in the race", i.Member.User.Username)
				ResponseEmbedInfo(s, i, true, fmt.Sprintf("You're already in the race %s", i.Member.User.Username), "You can't join the race twice, good luck with the race!")
				return
			}

			// Respond to the interaction with a message
			race.Render(s)
			ResponseEmbedSuccess(s, i, true, fmt.Sprintf("You've joined the race #%s", raceId), "We've just got your snail lined up at the starting line, good luck!")
		},
		models.RaceActionBet: func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			if len(options) != 1 {
				log.WithField("interaction", models.RaceActionBet).WithError(errors.New("invalid options")).Errorf("Not enough arguments/options from user %s", i.Member.User.Username)
				ResponseEmbedFail(s, i, true,
					fmt.Sprintf("I'm sorry %s, but there has been an issue", i.Member.User.Username),
					"There has been an issue with the action you sent, please try again.",
				)
				return
			}

			// Check if the user is initialised, if the user isn't initialised then
			// we need to tell them to initialise their account.
			_, err := models.GetUserByDiscordID(state.DB, i.Member.User.ID)
			if err != nil {
				log.WithField("interaction", models.RaceActionBet).WithError(err).Infof("No record for user %s", i.Member.User.Username)

				ResponseEmbedFail(s, i, true,
					fmt.Sprintf("I'm sorry %s, but you arent initialised", i.Member.User.Username),
					"You'll need to initialise your account with `/snailrace init` to use this command.",
				)
				return
			}

			// Check if the race exists, if it doesn't then we need to tell the
			// user
			raceId := options[0]
			race, ok := state.Races[raceId]
			if !ok {
				log.WithField("interaction", models.RaceActionJoin).WithError(errors.New("no existing race")).Infof("The raceid %s is not active, requested by user %s", raceId, i.Member.User.Username)
				ResponseEmbedFail(s, i, true, fmt.Sprintf("Race %s not avaliable", raceId), "There is currently no race with the ID you supplied.")
				return
			}

			// Check if the snail exists, if it doesn't then we need to tell the
			// user
			data := i.MessageComponentData()
			snailIndex, _ := strconv.Atoi(data.Values[0])
			snail := race.GetSnail(snailIndex)
			if snail == nil {
				log.WithField("interaction", models.RaceActionBet).WithError(err).Infof("User %s betting invalid snail", i.Member.User.Username)
				ResponseEmbedFail(s, i, true, fmt.Sprintf("Invalid snail to bet for race %s", raceId), "There is currently no snail with the ID you supplied.")
				return
			}

			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Flags: discordgo.MessageFlagsEphemeral,
					Embeds: []*discordgo.MessageEmbed{
						{
							Title:       "Looks like you want to make a bet",
							Color:       0x2ecc71,
							Description: fmt.Sprintf("So you want to make a bet on %s. Well select one of the following predetermined amounts, or use the following command for a custom amount: \n```\n/snailrace bet race_id: %s snail_index: %d amount: \n```\n", snail.Name, raceId, snailIndex),
						},
					},
					Components: []discordgo.MessageComponent{
						discordgo.ActionsRow{
							Components: []discordgo.MessageComponent{
								discordgo.Button{
									Label:    "5g",
									Style:    discordgo.SuccessButton,
									CustomID: fmt.Sprintf("%s:%s:%d:%d", models.RaceActionBetAmount, raceId, snailIndex, 5),
								},
								discordgo.Button{
									Label:    "10g",
									Style:    discordgo.SuccessButton,
									CustomID: fmt.Sprintf("%s:%s:%d:%d", models.RaceActionBetAmount, raceId, snailIndex, 10),
								},
								discordgo.Button{
									Label:    "20g",
									Style:    discordgo.SuccessButton,
									CustomID: fmt.Sprintf("%s:%s:%d:%d", models.RaceActionBetAmount, raceId, snailIndex, 20),
								},
							},
						},
					},
				},
			})
		},
		models.RaceActionBetAmount: func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			if len(options) != 3 {
				log.WithField("interaction", models.RaceActionBetAmount).WithError(errors.New("invalid options")).Errorf("Not enough arguments/options from user %s", i.Member.User.Username)
				ResponseEmbedFail(s, i, true,
					fmt.Sprintf("I'm sorry %s, but there has been an issue", i.Member.User.Username),
					"There has been an issue with the action you sent, please try again.",
				)
				return
			}

			// Check if the user is initialised, if the user isn't initialised then
			// we need to tell them to initialise their account.
			user, err := models.GetUserByDiscordID(state.DB, i.Member.User.ID)
			if err != nil {
				log.WithField("interaction", models.RaceActionBetAmount).WithError(err).Infof("No record for user %s", i.Member.User.Username)
				ResponseEmbedFail(s, i, true,
					fmt.Sprintf("I'm sorry %s, but you arent initialised", i.Member.User.Username),
					"You'll need to initialise your account with `/snailrace init` to use this command.",
				)
				return
			}

			// Check if the race exists, if it doesn't then we need to tell the
			// user
			raceId := options[0]
			race, ok := state.Races[raceId]
			if !ok {
				log.WithField("interaction", models.RaceActionBetAmount).WithError(errors.New("no existing race")).Warnf("The raceid %s is not active, requested by user %s", raceId, i.Member.User.Username)
				ResponseEmbedFail(s, i, true, fmt.Sprintf("Race %s not avaliable", raceId), "There is currently no race with the ID you supplied.")
				return
			}

			// Check if the snail exists, if it doesn't then we need to tell the
			// user
			snailIndex, _ := strconv.Atoi(options[1])
			snail := race.GetSnail(snailIndex)
			if snail == nil {
				log.WithField("interaction", models.RaceActionBetAmount).WithError(err).Warnf("User %s betting invalid snail", i.Member.User.Username)
				ResponseEmbedFail(s, i, true, fmt.Sprintf("Invalid snail to bet for race %s", raceId), "There is currently no snail with the ID you supplied.")
				return
			}

			// Check if the user has enough money to make the bet
			amount, _ := strconv.Atoi(options[2])
			if int(user.Money) < amount {
				log.WithField("interaction", models.RaceActionBetAmount).WithError(err).Infof("User %s doesn't have the funds to place a bet", i.Member.User.Username)
				ResponseEmbedFail(s, i, true, fmt.Sprintf("Sorry %s but you can't afford the bet", i.Member.User.Username), fmt.Sprintf("You don't have enough money to place that bet, you only have %d g.", user.Money))
				return
			}

			// Place the bet and remove the money from the user
			switch race.PlaceBet(snailIndex, amount, user.DiscordID) {
			case models.ErrInvalidSnail:
				log.WithField("interaction", models.RaceActionBetAmount).WithError(models.ErrInvalidSnail).Warnf("User %s failed to place bet on snail", i.Member.User.Username)
				ResponseEmbedFail(s, i, true, fmt.Sprintf("Sorry %s that snail doesn't exist", i.Member.User.Username), "The snail you have selected to bet is invalid, the snail isn't in the race.")
				return
			case models.ErrBetsClosed:
				log.WithField("interaction", models.RaceActionBetAmount).WithError(models.ErrBetsClosed).Warnf("User %s failed to place bet on snail as bets are closed", i.Member.User.Username)
				ResponseEmbedFail(s, i, true, fmt.Sprintf("Sorry %s Bets are Closed", i.Member.User.Username), "Bet's are closed so we can't accept your bet.")
				return
			case models.ErrNotEnough:
				log.WithField("interaction", models.RaceActionBetAmount).WithError(models.ErrNotEnough).Warnf("User %s failed to place bet on snail as there aren't enough racers in the race", i.Member.User.Username)
				ResponseEmbedFail(s, i, true, fmt.Sprintf("Sorry %s Not Enough Racers", i.Member.User.Username), "We need at least 2 racers to enable bets.")
				return
			}

			ResponseEmbedSuccess(s, i, true, fmt.Sprintf("Bet placed for %s", snail.Name), fmt.Sprintf("You've placed a bet for %s of %d g", snail.Name, amount))
			user.RemoveMoney(state.DB, uint64(amount))
		},
	}
}

func (c *CommandHostRace) ModalHandler(state *models.State, options ...string) map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate) {
	return map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate){}
}
