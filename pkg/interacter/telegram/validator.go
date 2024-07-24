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

	validatorsInfo := interacter.DataFetcher.FindValidator(args[1])

	template, err := interacter.TemplateManager.Render("validator", validatorsInfo)
	if err != nil {
		interacter.Logger.Error().Err(err).Msg("Error rendering template")
		return interacter.BotReply(c, "Error rendering template")
	}
	//
	// fmt.Printf("template: %s\n", template)

	return interacter.BotReply(c, template)
}
