package proxy

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/Evan-Lab/cloud-native/lib/go/discord"
	"github.com/GoogleCloudPlatform/functions-framework-go/functions"
	"github.com/bwmarrin/discordgo"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/propagation"
)

var (
	tracer = otel.Tracer("github.com/Evan-Lab/cloud-native/functions/proxy")
)

func init() {
	functions.HTTP("DiscordProxy", DiscordProxy)
}

func DiscordProxy(w http.ResponseWriter, r *http.Request) {
	ctx := otel.GetTextMapPropagator().Extract(r.Context(), propagation.HeaderCarrier(r.Header))
	ctx, span := tracer.Start(ctx, "DiscordProxy")
	defer span.End()
	r = r.WithContext(ctx)

	req, err := discord.ParseRequest(w, r)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "failed to parse interaction request")
		slog.WarnContext(ctx, "Failed to parse request", "error", err)
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	span.SetAttributes(attribute.Int("discord.interaction.type", int(req.Type)))
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
		span.SetStatus(codes.Error, "unknown interaction type")
		slog.WarnContext(ctx, "Unknown interaction type", "type", req.Type)
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "Proxy handler returned error")
		slog.ErrorContext(ctx, "Failed to handle interaction", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	slog.InfoContext(ctx, "Proxy response", "response", resp)
	discord.SetHeaders(w.Header())
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "failed to encode response")
		slog.ErrorContext(ctx, "Failed to encode response", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}
