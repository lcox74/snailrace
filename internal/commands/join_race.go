package commands

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
	"github.com/lcox74/snailrace/internal/models"
	log "github.com/sirupsen/logrus"
)

type CommandJoinRace struct{}

func (c *CommandJoinRace) Decleration() *discordgo.ApplicationCommandOption {
	return &discordgo.ApplicationCommandOption{
		Name:        "join",
		Description: "Let's join a race",
		Type:        discordgo.ApplicationCommandOptionSubCommand,
		Options: []*discordgo.ApplicationCommandOption{
			{
				Name:        "race_id",
				Description: "The race to join",
				Type:        discordgo.ApplicationCommandOptionString,
				Required:    true,
			},
		},
	}
}

func (c *CommandJoinRace) AppHandler(state *models.State) func(s *discordgo.Session, i *discordgo.InteractionCreate) {
	return func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		log.Printf("[CMD] Join!\n")

		raceId := ""

		// The Join Action acts as the command /snailrace join <race_id>
		// If the caller doesn't supply the `race_id` then we need to
		// through and error, theoretically this should nevery error
		if len(i.ApplicationCommandData().Options) > 0 {
			for _, opt := range i.ApplicationCommandData().Options[0].Options {
				if opt.Name == "race_id" {
					raceId = opt.Value.(string)
				}
			}
		}

		// Check if there is a raceId supplied, if there isn't then we need to
		// tell the user that they need to supply a raceId
		if raceId == "" {
			ResponseEmbedFail(s, i, true,
				"There is no RaceId supplied",
				"Please try again by supplying a race RaceId.",
			)
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

		// Fetch the race from the supplied raceId, if there is no race with the
		// RaceId then warn the user.
		race, ok := state.Races[raceId]
		if !ok {
			ResponseEmbedFail(s, i, true, fmt.Sprintf("Race %s not avaliable", raceId), "There is currently no race with the ID you supplied.")
			return
		}

		// Add the snail to the race and
		err = race.AddSnail(snail)
		if err != nil {
			ResponseEmbedInfo(s, i, true, fmt.Sprintf("You're already in the race %s", i.Member.User.Username), "You can't join the race twice, good luck with the race!")
			return
		}

		race.Render(s)
		ResponseEmbedSuccess(s, i, true, fmt.Sprintf("You've joined the race #%s", raceId), "We've just got your snail lined up at the starting line, good luck!")
	}
}

func (c *CommandJoinRace) ActionHandler(state *models.State, options ...string) map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate) {
	return map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate){}
}
