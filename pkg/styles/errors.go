package styles

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
)

// Supply a consistent way to respond to an User Interaction with an error.
// This will create a embed with the "Error" title and content provided, and set
// the ephemeral flag if required as well as the colour to red
func RespondErr(session *discordgo.Session, interaction *discordgo.Interaction, ephemeral bool, msg string) {
	session.InteractionRespond(
		interaction,
		createResponse("Error", msg, ephemeral, ColourError, nil),
	)
}

// If sending a message that isn't a response to an interaction, use the
// following. It will create the same design as RespondErr but without the
// ephemeral flag.
func SendErr(session *discordgo.Session, channelId string, msg string) {
	session.ChannelMessageSendComplex(
		channelId,
		createMessage("Error", msg, ColourError, nil),
	)
}

// =============================
//       Standard Errs
// =============================

func ErrDm(session *discordgo.Session, interaction *discordgo.Interaction) {
	RespondErr(session, interaction, false,
		"I'm sorry, but we don't support DMs. Please use this command in a server.",
	)
}
func ErrInitialise(session *discordgo.Session, interaction *discordgo.Interaction, userMention string) {
	RespondErr(session, interaction, true,
		fmt.Sprintf("I'm sorry %s, but you arent initialised. You'll need to initialise your account with `/snailrace init` to use this command.", userMention),
	)
}
func ErrInvalidRace(session *discordgo.Session, interaction *discordgo.Interaction, userMention string) {
	RespondErr(session, interaction, true,
		fmt.Sprintf("I'm sorry %s, the race id you've supplied is not an active race.", userMention),
	)
}
func ErrNoActiveSnail(session *discordgo.Session, interaction *discordgo.Interaction, userMention string) {
	RespondErr(session, interaction, true,
		fmt.Sprintf("I'm sorry %s, but you don't have an active snail set on your account.", userMention),
	)
}
func ErrAlreadyInRace(session *discordgo.Session, interaction *discordgo.Interaction) {
	RespondErr(session, interaction, true,
		"You can't join the race twice, good luck with the race!",
	)
}
func ErrRaceNotExist(session *discordgo.Session, interaction *discordgo.Interaction) {
	RespondErr(session, interaction, true,
		"There is currently no race with the ID you supplied.",
	)
}
func ErrRaceClosed(session *discordgo.Session, interaction *discordgo.Interaction) {
	RespondErr(session, interaction, true,
		"The race you have just tried to join is currently closed.",
	)
}
func ErrRaceFull(session *discordgo.Session, interaction *discordgo.Interaction) {
	RespondErr(session, interaction, true,
		"The race you have just tried to join is currently full. **MAX 10 Snails**.",
	)
}
func ErrInvalidSnailBet(session *discordgo.Session, interaction *discordgo.Interaction) {
	RespondErr(session, interaction, true,
		"There is currently no snail with the ID you supplied.",
	)
}
func ErrCantAfford(session *discordgo.Session, interaction *discordgo.Interaction) {
	RespondErr(session, interaction, true,
		"You don't have a enough money to afford that action.",
	)
}
