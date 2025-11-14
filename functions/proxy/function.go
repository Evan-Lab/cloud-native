package proxy

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/Evan-Lab/cloud-native/lib/go/discord"
	"github.com/GoogleCloudPlatform/functions-framework-go/functions"
	"github.com/bwmarrin/discordgo"
)

func init() {
	functions.HTTP("DiscordProxy", DiscordProxy)
}

func DiscordProxy(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// Setup logging
	setupLogging()

	// Setup metrics, tracing, and context propagation
	shutdown, err := setupOpenTelemetry(ctx)
	if err != nil {
		slog.ErrorContext(ctx, "error setting up OpenTelemetry", slog.Any("error", err))
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	defer shutdown(ctx)

	req, err := discord.ParseRequest(w, r)
	if err != nil {
		slog.WarnContext(ctx, "Failed to parse request", "error", err)
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	slog.InfoContext(ctx, "Handling interaction", "type", req.Type, "interaction", req)
	var resp *discordgo.InteractionResponse
	switch req.Type {
	case discordgo.InteractionApplicationCommand:
		resp, err = cmdProxy(ctx, req)
	// case discordgo.InteractionApplicationCommandAutocomplete:
	// 	panic("Not implemented")
	case discordgo.InteractionPing:
		resp, err = ping(ctx)
	default:
		slog.WarnContext(ctx, "Unknown interaction type", "type", req.Type)
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}
	if err != nil {
		slog.ErrorContext(ctx, "Failed to handle interaction", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	slog.InfoContext(ctx, "Sending response", "response", resp)
	discord.SetHeaders(w.Header())
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		slog.ErrorContext(ctx, "Failed to encode response", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}
