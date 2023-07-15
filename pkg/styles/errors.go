package styles

import "github.com/bwmarrin/discordgo"

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
