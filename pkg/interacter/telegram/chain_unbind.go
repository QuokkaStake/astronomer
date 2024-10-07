package telegram

import (
	"errors"
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
	valid, usage, args := interacter.SingleArgParser(c.Text(), "chain")
	if !valid {
		return usage, constants.ErrWrongInvocation
	}

	chain, err := interacter.Database.GetChainByName(args.Value)
	if err != nil && errors.Is(err, constants.ErrChainNotFound) {
		return interacter.ChainNotFound()
	} else if err != nil {
		return "", err
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
		return "Chain is not bound to this chat!", constants.ErrChainNotBound
	}

	return interacter.TemplateManager.Render("chain_unbind", chain)
}
