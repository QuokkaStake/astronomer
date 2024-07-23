package telegram

import (
	"fmt"
	"html"
	"strings"

	tele "gopkg.in/telebot.v3"
)

func (interacter *Interacter) HandleValidator(c tele.Context) error {
	interacter.Logger.Info().
		Str("sender", c.Sender().Username).
		Str("text", c.Text()).
		Msg("Got status query")

	args := strings.SplitN(c.Text(), " ", 2)
	if len(args) < 2 {
		return interacter.BotReply(c, html.EscapeString(fmt.Sprintf(
			"Usage: %s <query>",
			args[0],
		)))
	}

	validators := interacter.DataFetcher.FindValidator(args[1])

	var sb strings.Builder

	for _, info := range validators {
		sb.WriteString(fmt.Sprintf("chain: %+v\n", info.Chain))
		sb.WriteString(fmt.Sprintf("error: %+v\n", info.Error))
		sb.WriteString(fmt.Sprintf("validator: %+v\n", info.Validator))
		sb.WriteString("\n------------\n")
	}

	return interacter.BotReply(c, html.EscapeString(sb.String()))
}
