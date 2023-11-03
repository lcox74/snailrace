package commands

import (
	"fmt"
	"math"
	"strings"

	"github.com/lcox74/snailrace/internal/models"
	log "github.com/sirupsen/logrus"

	"github.com/bwmarrin/discordgo"
)

var STARTING_INDEX = 0

// CommandBackpack shows the user their backpack of snails
type CommandBackpack struct{}

const (
	// Action Ids
	BackpackActionNextPage = "next_page"
	BackpackActionPrevPage = "prev_page"
)

func (c *CommandBackpack) Decleration() *discordgo.ApplicationCommandOption {
	return &discordgo.ApplicationCommandOption{
		Name:        "backpack",
		Description: "Check your snail backpack",
		Type:        discordgo.ApplicationCommandOptionSubCommand,
	}
}

func (c *CommandBackpack) AppHandler(state *models.State) func(s *discordgo.Session, i *discordgo.InteractionCreate) {
	return func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		log.WithField("cmd", "/backpack").Infof("Getting backpack for user %s", i.Member.User.Username)
		user, snails, active := GetBackpackState(state, s, i) // error handling done within
		RenderBackpack(state, s, i, snails, active, user)
	}
}

// this is in display_profile, remove later when it gets merged
func GenerateProgressBar(percentage float64) string {
	numSquares := int(math.Floor(percentage / 10))
	progress := strings.Repeat("🟩", numSquares)
	progress += strings.Repeat("⬛", 10-numSquares)
	return progress
}

func (c *CommandBackpack) ActionHandler(state *models.State, options ...string) map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate) {
	return map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate){
		BackpackActionPrevPage: func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			log.WithField("interaction", BackpackActionPrevPage).Info("PREV page for snails has been requested")
			// reduce the pagination
			STARTING_INDEX -= 10
			// re render the modal
			user, snails, active := GetBackpackState(state, s, i) // error handling done within
			RenderBackpack(state, s, i, snails, active, user)
		},
		BackpackActionNextPage: func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			log.WithField("interaction", BackpackActionPrevPage).Info("NEXT page for snails has been requested")
			// reduce the pagination
			STARTING_INDEX += 10
			// re render the modal
			user, snails, active := GetBackpackState(state, s, i) // error handling done within
			RenderBackpack(state, s, i, snails, active, user)
		},
	}
}

func (c *CommandBackpack) ModalHandler(state *models.State, options ...string) map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate) {
	return map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate){}
}

func GetBackpackState(state *models.State, s *discordgo.Session, i *discordgo.InteractionCreate) (*models.User, []models.Snail, *models.Snail) {
	user, err := models.GetUserByDiscordID(state.DB, i.Member.User.ID)
	if err != nil {
		log.WithField("cmd", "/backpack").WithError(err).Infof("User %s is not initialised", i.Member.User.Username)
		ResponseEmbedFail(s, i, true,
			fmt.Sprintf("I'm sorry %s, but you arent initialised", i.Member.User.Username),
			"You'll need to initialise your account with `/snailrace init` to use this command.",
		)
		return user, nil, nil
	}

	snails, err := models.GetAllSnails(state.DB, *user)
	if err != nil {
		log.WithField("cmd", "/backpack").WithError(err).Infof("There was an error retrieving all the snails for user %s", i.Member.User.Username)
		ResponseEmbedFail(s, i, true,
			fmt.Sprintf("Error retrieving snails for user %s", i.Member.User.Username),
			"Something went wrong attempting to retrieve snails for this user. Please try again.",
		)
		return user, snails, nil
	}

	active, err := models.GetActiveSnail(state.DB, *user)
	if err != nil {
		log.WithField("cmd", "/backpack").WithError(err).Infof("There was an error retrieving the active snail for the user %s", i.Member.User.Username)
		ResponseEmbedFail(s, i, true,
			fmt.Sprintf("Error retrieving the active snail for user %s", i.Member.User.Username),
			"Something went wrong attempting to retrieve the active snail for this user. Please try again.",
		)
		return user, snails, active
	}

	return user, snails, active
}

func RenderBackpack(state *models.State, s *discordgo.Session, i *discordgo.InteractionCreate,
	snails []models.Snail, active *models.Snail, user *models.User) {
	backpack := fmt.Sprintf("Active Snail: **%s**\n\nSnails: **(%d - %d)**\n", active.Name, STARTING_INDEX+1, STARTING_INDEX+10)
	for i := STARTING_INDEX; i < len(snails); i++ {
		// Get information about snail and pretty it up
		snailProgress := models.GetLevelProgress(&snails[i])
		progressBar := GenerateProgressBar(snailProgress)
		backpack += fmt.Sprintf("%d - `%s`  lvl.%d  %s\n", snails[i].ID, snails[i].Name, snails[i].Level, progressBar)
	}

	paginationDisabled := len(snails) <= 10

	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Flags: discordgo.MessageFlagsEphemeral,
			Embeds: []*discordgo.MessageEmbed{
				{
					Title:       "Your Backpack",
					Color:       0x2ecc71,
					Description: backpack,
				},
			},
			Components: []discordgo.MessageComponent{
				discordgo.ActionsRow{
					Components: []discordgo.MessageComponent{
						discordgo.Button{
							Label:    "Prev",
							Style:    discordgo.SuccessButton,
							CustomID: fmt.Sprintf("%s:%d", BackpackActionPrevPage, user.ID),
							Disabled: paginationDisabled || STARTING_INDEX == 0,
						},
						discordgo.Button{
							Label:    "Next",
							Style:    discordgo.SuccessButton,
							CustomID: fmt.Sprintf("%s:%d", BackpackActionNextPage, user.ID),
							Disabled: paginationDisabled || len(snails)-STARTING_INDEX < 10,
						},
					},
				},
			},
		},
	})
}
