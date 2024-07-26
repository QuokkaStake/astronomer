package telegram

import (
	"fmt"
	"html"
	"strconv"
	"strings"

	tele "gopkg.in/telebot.v3"
)

func (interacter *Interacter) GetChainBindCommand() Command {
	return Command{
		Name:    "chain_bind",
		Execute: interacter.HandleChainBind,
	}
}

func (interacter *Interacter) HandleChainBind(c tele.Context) (string, error) {
	interacter.Logger.Info().
		Str("sender", c.Sender().Username).
		Str("text", c.Text()).
		Msg("Got status query")

	args := strings.Split(c.Text(), " ")
	if len(args) != 2 {
		return html.EscapeString(fmt.Sprintf(
			"Usage: %s <chain>",
			args[0],
		)), fmt.Errorf("invalid command invocation")
	}

	chain := interacter.Chains.FindByName(args[1])
	if chain == nil {
		return html.EscapeString(fmt.Sprintf(
			"Could not find a chain with the name '%s'",
			args[1],
		)), fmt.Errorf("could not find chain to bind")
	}

	err := interacter.Database.InsertChainBind(
		interacter.Name(),
		strconv.FormatInt(c.Chat().ID, 10),
		c.Chat().Title,
		chain.Name,
	)
	if err != nil {
		interacter.Logger.Error().Err(err).Msg("Error inserting chain bind")
		return "", err
	}

	return "Successfully added a chain bind to this chat!", nil
}
