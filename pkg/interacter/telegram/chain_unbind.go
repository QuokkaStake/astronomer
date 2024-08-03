package telegram

import (
	"fmt"
	"html"
	"main/pkg/constants"
	"strconv"

	tele "gopkg.in/telebot.v3"
)

func (interacter *Interacter) GetChainUnbindCommand() Command {
	return Command{
		Name:    "chain_unbind",
		Execute: interacter.HandleChainUnbind,
	}
}

func (interacter *Interacter) HandleChainUnbind(c tele.Context, chainBinds []string) (string, error) {
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
