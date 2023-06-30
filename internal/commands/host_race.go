package commands

import (
	"fmt"

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

		// Start the race as a seperate process
		go models.StartRace(s, race)

		// Respond to the interaction with a message
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Flags: discordgo.MessageFlagsEphemeral,
				Embeds: []*discordgo.MessageEmbed{
					{
						Title:       fmt.Sprintf("You just hosted a race %s!", i.Member.User.Username),
						Description: "Your snail is officially waiting at the starting line for other snails to join.",
					},
				},
			},
		})
	}
}

func (c *CommandHostRace) ActionHandler(state *models.State, options ...string) map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate) {
	return map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate){
		models.RaceActionJoin: func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			log.Printf("[CMD] Host Form Interaction!\n")

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
			}

			race.Render(s)

			// Respond to the interaction with a message
			ResponseEmbedSuccess(s, i, true, "Not Implemented", "This feature is not implemented yet. Sorry!")
		},
	}
}
