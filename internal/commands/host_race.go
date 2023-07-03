package commands

import (
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
		log.Printf("[CMD] Host!\n")

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

		// We need to get the active snail of the host to automatically add them
		// to the race
		snail, err := models.GetActiveSnail(state.DB, *user)
		if err != nil {
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
					race.SetNoFill()
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
			log.Printf("[CMD] Host Join Interaction!\n")

			// The Join Action acts as the command /snailrace join <race_id>
			// If the caller doesn't supply the `race_id` then we need to
			// through and error, theoretically this should nevery error
			if len(options) != 1 {
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
				ResponseEmbedFail(s, i, true,
					fmt.Sprintf("I'm sorry %s, but you arent initialised", i.Member.User.Username),
					"You'll need to initialise your account with `/snailrace init` to use this command.",
				)
				return
			}

			// We neet to get the user's active snail to add to the race
			snail, err := models.GetActiveSnail(state.DB, *user)
			if err != nil {
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
				ResponseEmbedFail(s, i, true, fmt.Sprintf("Race %s not avaliable", raceId), "There is currently no race with the ID you supplied.")
				return
			}

			err = race.AddSnail(snail)
			if err != nil {
				ResponseEmbedInfo(s, i, true, fmt.Sprintf("You're already in the race %s", i.Member.User.Username), "You can't join the race twice, good luck with the race!")
				return
			}

			// Respond to the interaction with a message
			race.Render(s)
			ResponseEmbedSuccess(s, i, true, fmt.Sprintf("You've joined the race #%s", raceId), "We've just got your snail lined up at the starting line, good luck!")
		},
		models.RaceActionBet: func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			log.Printf("[CMD] Host Bet Interaction!\n")
			if len(options) != 1 {
				log.Errorf("Invalid number of options for bet amount: %d -> %s\n", len(options), options[0])
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
				ResponseEmbedFail(s, i, true, fmt.Sprintf("Race %s not avaliable", raceId), "There is currently no race with the ID you supplied.")
				return
			}

			// Check if the snail exists, if it doesn't then we need to tell the
			// user
			data := i.MessageComponentData()
			snailIndex, _ := strconv.Atoi(data.Values[0])
			snail := race.GetSnail(snailIndex)
			if snail == nil {
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
							Description: fmt.Sprintf("So you want to make a bet on %s. Well how much? Enter the amount as a number.", snail.Name),
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
				ResponseEmbedFail(s, i, true, fmt.Sprintf("Race %s not avaliable", raceId), "There is currently no race with the ID you supplied.")
				return
			}

			// Check if the snail exists, if it doesn't then we need to tell the
			// user
			snailIndex, _ := strconv.Atoi(options[1])
			snail := race.GetSnail(snailIndex)
			if snail == nil {
				ResponseEmbedFail(s, i, true, fmt.Sprintf("Invalid snail to bet for race %s", raceId), "There is currently no snail with the ID you supplied.")
				return
			}

			// Check if the user has enough money to make the bet
			ammount, _ := strconv.Atoi(options[2])
			if int(user.Money) < ammount {
				log.Warnf("Player %s tried to bet %d g but only has %d g", i.Member.User.Username, ammount, user.Money)
				ResponseEmbedFail(s, i, true, fmt.Sprintf("Sorry %s but you can't afford the bet", i.Member.User.Username), fmt.Sprintf("You don't have enough money to place that bet, you only have %d g.", user.Money))
				return
			}

			// Place the bet and remove the money from the user
			race.PlaceBet(snailIndex, ammount, user.DiscordID)
			user.RemoveMoney(state.DB, uint64(ammount))
			ResponseEmbedSuccess(s, i, true, fmt.Sprintf("Bet placed for %s", snail.Name), fmt.Sprintf("You've placed a bet for %s of %d g", snail.Name, ammount))
		},
	}
}

func (c *CommandHostRace) ModalHandler(state *models.State, options ...string) map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate) {
	return map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate){}
}
