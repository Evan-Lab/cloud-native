package snap

import (
	"context"
	"log/slog"

	"github.com/Evan-Lab/cloud-native/lib/go/discord"
	"github.com/Evan-Lab/cloud-native/lib/go/utils"
	"github.com/bwmarrin/discordgo"
)

func RespondToInteraction(ctx context.Context, interaction_token, pngUrl string) error {
	ctx, span := tracer.Start(ctx, "RespondToInteraction")
	defer span.End()

	s, err := discord.Session()
	if err != nil {
		slog.ErrorContext(ctx, "discord.Session", "error", err)
		span.RecordError(err)
		return err
	}
	appID := s.State.Application.ID

	embed := &discordgo.MessageEmbed{
		Image: &discordgo.MessageEmbedImage{
			URL: pngUrl,
		},
	}

	responseData := discordgo.WebhookEdit{
		Content: utils.Ptr(""),
		Embeds:  &[]*discordgo.MessageEmbed{embed},
	}

	st, err := s.WebhookMessageEdit(appID, interaction_token, "@original", &responseData)
	if err != nil {
		slog.ErrorContext(ctx, "WebhookMessageEdit", "error", err)
		span.RecordError(err)
		return err
	}

	slog.InfoContext(ctx, "Edited interaction response", "message_id", st.ID, "channel_id", st.ChannelID)

	return nil
}
