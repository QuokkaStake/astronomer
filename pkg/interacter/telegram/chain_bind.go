package telegram

import (
	"fmt"
	"html"
	"main/pkg/constants"
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

func (interacter *Interacter) HandleChainBind(c tele.Context, chainBinds []string) (string, error) {
	valid, usage, args := interacter.SingleArgParser(c, "chain")
	if !valid {
		return usage, constants.ErrWrongInvocation
	}

	chain := interacter.Chains.FindByName(args.Value)
	if chain == nil {
		return html.EscapeString(fmt.Sprintf(
			"Could not find a chain with the name '%s'",
			args.Value,
		)), constants.ErrChainNotFound
	}

	err := interacter.Database.InsertChainBind(
		interacter.Name(),
		strconv.FormatInt(c.Chat().ID, 10),
		c.Chat().Title,
		chain.Name,
	)
	if err != nil {
		if strings.Contains(err.Error(), "duplicate key value") {
			return "This chain is already bound to this chat!", err
		}

		interacter.Logger.Error().Err(err).Msg("Error inserting chain bind")
		return "", err
	}

	return "Successfully added a chain bind to this chat!", nil
}
