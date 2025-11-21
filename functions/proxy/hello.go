package functions

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/bwmarrin/discordgo"
)

func init() {
	RegisterCommand("hello", helloCmd)
}

func helloCmd(ctx context.Context, interaction discordgo.Interaction, data discordgo.ApplicationCommandInteractionData) (*discordgo.InteractionResponse, error) {
	name := "World"
	if opt := data.GetOption("name"); opt != nil {
		if s, ok := opt.Value.(string); ok && s != "" {
			name = s
		}
	}
	slog.Info("Received hello command", "name", name)

	return &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: fmt.Sprintf("Hello, %s!", name),
		},
	}, nil
}
