package discord

import (
	"net/http"

	"github.com/bwmarrin/discordgo"
)

func SetHeaders(h http.Header) {
	h.Set("Content-Type", "application/json")
	h.Set("User-Agent", "DiscordBot (https://github.com/bwmarrin/discordgo, v"+discordgo.VERSION+")")
}
