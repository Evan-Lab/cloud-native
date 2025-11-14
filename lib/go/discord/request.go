package discord

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/bwmarrin/discordgo"
)

func ParseRequest(w http.ResponseWriter, r *http.Request) (discordgo.Interaction, error) {
	if !discordgo.VerifyInteraction(r, pubKey) {
		return discordgo.Interaction{}, fmt.Errorf("invalid request signature")
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		return discordgo.Interaction{}, err
	}

	var interaction discordgo.Interaction
	if err := json.Unmarshal(body, &interaction); err != nil {
		return discordgo.Interaction{}, err
	}

	return interaction, nil
}
