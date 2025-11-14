package functions

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/EpitechPGE45-2025/G-CLO-910-PAR-9-1-cloudnative-13/lib/go/discord"

	"github.com/GoogleCloudPlatform/functions-framework-go/functions"
	"github.com/bwmarrin/discordgo"
)

func init() {
	functions.HTTP("DiscordProxy", DiscordProxy)
}

func ping() (*discordgo.InteractionResponse, error) {
	slog.Info("Received Ping interaction")
	return &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponsePong,
	}, nil
}

type CommandHandler func(ctx context.Context, interaction discordgo.Interaction, data discordgo.ApplicationCommandInteractionData) (*discordgo.InteractionResponse, error)

var cmds = make(map[string]CommandHandler)

func RegisterCommand(cmd string, handler CommandHandler) {
	cmds[cmd] = handler
}

func cmdProxy(ctx context.Context, interaction discordgo.Interaction) (*discordgo.InteractionResponse, error) {
	data := interaction.ApplicationCommandData()
	handler, ok := cmds[data.Name]
	if !ok {
		slog.Warn("No handler for command", "name", data.Name)
		return &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "Unknown command",
			},
		}, nil
	}

	return handler(ctx, interaction, data)
}

// HelloWorld writes "Hello, World!" to the HTTP response.
func DiscordProxy(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	req, err := discord.ParseRequest(w, r)
	if err != nil {
		slog.Warn("Failed to parse request", "error", err)
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	slog.Info("Handling interaction", "type", req.Type, "interaction", req)
	var resp *discordgo.InteractionResponse
	switch req.Type {
	case discordgo.InteractionApplicationCommand:
		resp, err = cmdProxy(ctx, req)
	// case discordgo.InteractionApplicationCommandAutocomplete:
	// 	panic("Not implemented")
	case discordgo.InteractionPing:
		resp, err = ping()
	default:
		slog.Warn("Unknown interaction type", "type", req.Type)
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}
	if err != nil {
		slog.Error("Failed to handle interaction", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	slog.Info("Sending response", "response", resp)
	discord.SetHeaders(w.Header())
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		slog.Error("Failed to encode response", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}
