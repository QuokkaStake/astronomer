package telegram

import (
	"fmt"
	"html"
	"strconv"
	"strings"

	tele "gopkg.in/telebot.v3"
)

func (interacter *Interacter) GetChainUnbindCommand() Command {
	return Command{
		Name:    "chain_unbind",
		Execute: interacter.HandleChainUnbind,
		MinArgs: 1,
		Usage:   "Usage: %s <chain>",
	}
}

func (interacter *Interacter) HandleChainUnbind(c tele.Context) (string, error) {
	args := strings.Split(c.Text(), " ")

	chain := interacter.Chains.FindByName(args[1])
	if chain == nil {
		return html.EscapeString(fmt.Sprintf(
			"Could not find a chain with the name '%s'",
			args[1],
		)), fmt.Errorf("could not find chain to unbound")
	}

	deleted, err := interacter.Database.DeleteChainBind(
		interacter.Name(),
		strconv.FormatInt(c.Chat().ID, 10),
		chain.Name,
	)
	if err != nil {
		interacter.Logger.Error().Err(err).Msg("Error inserting chain bind")
		return "", err
	}

	if !deleted {
		interacter.Logger.Error().Err(err).Msg("Chain is not bound to this chat!")
		return "", err
	}

	return "Successfully removed a chain bind from this chat!", nil
}
