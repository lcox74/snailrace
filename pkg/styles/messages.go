// Suplies consistent Discord styling for messages and interactions
package styles

import "github.com/bwmarrin/discordgo"

const (
	ColourOk    int = 0x2ECC71 // Green
	ColourInfo  int = 0x3498DB // Blue
	ColourError int = 0xE74C3C // Red
)

// Converts the ephemeral bool to a discordgo.MessageFlags to be used in the
// discordgo.InteractionResponseData struct
func ephemeralFlag(ephemeral bool) discordgo.MessageFlags {
	if ephemeral {
		return discordgo.MessageFlagsEphemeral
	}
	return discordgo.MessageFlags(0)
}

// A consistent way to respond to an User Interaction. It will create a embed
// with the title and content provided, and set the ephemeral flag if required.
func createResponse(title, content string, ephemeral bool, colour int, components []discordgo.MessageComponent) *discordgo.InteractionResponse {
	return &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Flags:   ephemeralFlag(ephemeral),
			Content: content,
			Embeds: []*discordgo.MessageEmbed{
				{
					Title:       title,
					Description: content,
					Color:       colour,
				},
			},
			Components: components,
		},
	}
}

// If sending a message that isn't a response to an interaction, this should be
// used. It will create a embed with the title and content provided and a
// reference to the message to allow editing.
func createMessage(title, content string, colour int, components []discordgo.MessageComponent) *discordgo.MessageSend {
	return &discordgo.MessageSend{
		Content: content,
		Embeds: []*discordgo.MessageEmbed{
			{
				Title:       title,
				Description: content,
				Color:       colour,
			},
		},
		Components: components,
	}
}

// Supply a consistent way to respond to an User Interaction with an Ok message.
// This will create a embed with the supplied title and content provided, and
// setthe ephemeral flag if required as well as the colour to green
func RespondOk(session *discordgo.Session, interaction *discordgo.Interaction, ephemeral bool, title, msg string, components []discordgo.MessageComponent) {
	session.InteractionRespond(
		interaction,
		createResponse(title, msg, ephemeral, ColourOk, components),
	)
}

// If sending a message that isn't a response to an interaction, use the
// following. It will create the same design as RespondOk but without the
// ephemeral flag.
func SendOk(session *discordgo.Session, channelId string, title, msg string, components []discordgo.MessageComponent) (*discordgo.Message, error) {
	return session.ChannelMessageSendComplex(
		channelId,
		createMessage(title, msg, ColourOk, components),
	)
}

// A consistent way to edit a given message already sent in a channel. It will
// update the message to display the same design as the SendOk function.
func EditOk(session *discordgo.Session, channelId string, messageId string, title, msg string, components []discordgo.MessageComponent) (*discordgo.Message, error) {
	// Create the new message
	message := createMessage(title, msg, ColourOk, components)

	// Fetch the message to get the components and then apply the edit
	edit := discordgo.NewMessageEdit(channelId, messageId)
	edit.Content = &message.Content
	edit.Embed = message.Embeds[0]

	// Send the edit
	return session.ChannelMessageEditComplex(edit)
}

// Supply a consistent way to respond to an User Interaction with an Ok message.
// This will create a embed with the supplied title and content provided, and
// setthe ephemeral flag if required as well as the colour to green
func RespondInfo(session *discordgo.Session, interaction *discordgo.Interaction, ephemeral bool, title, msg string, components []discordgo.MessageComponent) {
	session.InteractionRespond(
		interaction,
		createResponse(title, msg, ephemeral, ColourInfo, components),
	)
}

// If sending a message that isn't a response to an interaction, use the
// following. It will create the same design as RespondOk but without the
// ephemeral flag.
func SendInfo(session *discordgo.Session, channelId string, title, msg string, components []discordgo.MessageComponent) (*discordgo.Message, error) {
	return session.ChannelMessageSendComplex(
		channelId,
		createMessage(title, msg, ColourInfo, components),
	)
}

// A consistent way to edit a given message already sent in a channel. It will
// update the message to display the same design as the SendOk function.
func EditInfo(session *discordgo.Session, channelId string, messageId string, title, msg string, components []discordgo.MessageComponent) (*discordgo.Message, error) {
	// Create the new message
	message := createMessage(title, msg, ColourInfo, components)

	// Fetch the message to get the components and then apply the edit
	edit := discordgo.NewMessageEdit(channelId, messageId)
	edit.Content = &message.Content
	edit.Embed = message.Embeds[0]
	edit.Components = message.Components

	// Send the edit
	return session.ChannelMessageEditComplex(edit)
}
