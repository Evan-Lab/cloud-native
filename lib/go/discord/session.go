package discord

import (
	"context"
	"crypto/ed25519"
	"encoding/hex"
	"log/slog"
	"os"

	"github.com/EpitechPGE45-2025/G-CLO-910-PAR-9-1-cloudnative-13/lib/go/secrets"
	"github.com/EpitechPGE45-2025/G-CLO-910-PAR-9-1-cloudnative-13/lib/go/utils"
	"github.com/bwmarrin/discordgo"
)

var token string
var pubKey ed25519.PublicKey

func init() {
	ctx := context.Background()
	token = string(utils.Must(secrets.Secret(ctx, "DISCORD_BOT_TOKEN")))
	keyHex := utils.Must(secrets.Secret(ctx, "DISCORD_PUBLIC_KEY"))
	pubKey = utils.Must(hex.DecodeString(string(keyHex)))

	if token == "" || pubKey == nil {
		slog.Error("DISCORD_BOT_TOKEN or DISCORD_PUBLIC_KEY not set in environment")
		os.Exit(2)
	}
}

func Session() (*discordgo.Session, error) {

	s, err := discordgo.New("Bot " + token)
	if err != nil {
		slog.Error("Can not create session", "error", err)
		return nil, err
	}

	if err := s.Open(); err != nil {
		slog.Error("Can not open session", "error", err)
		return nil, err
	}
	return s, nil
}
