package telegram

import (
	"fmt"
	"html"
	"strconv"
	"strings"

	tele "gopkg.in/telebot.v3"
)

func (interacter *Interacter) HandleChainBind(c tele.Context) error {
	interacter.Logger.Info().
		Str("sender", c.Sender().Username).
		Str("text", c.Text()).
		Msg("Got status query")

	args := strings.Split(c.Text(), " ")
	if len(args) != 2 {
		return interacter.BotReply(c, html.EscapeString(fmt.Sprintf(
			"Usage: %s <chain>",
			args[0],
		)))
	}

	chain := interacter.Chains.FindByName(args[1])
	if chain == nil {
		return interacter.BotReply(c, html.EscapeString(fmt.Sprintf(
			"Could not find a chain with the name '%s'",
			args[1],
		)))
	}

	err := interacter.Database.InsertChainBind(
		interacter.Name(),
		strconv.FormatInt(c.Chat().ID, 10),
		c.Chat().Title,
		chain.Name,
	)
	if err != nil {
		interacter.Logger.Error().Err(err).Msg("Error inserting chain bind")
		return interacter.BotReply(c, "Internal error!")
	}

	return interacter.BotReply(c, "Successfully added a chain bind to this chat!")
}
