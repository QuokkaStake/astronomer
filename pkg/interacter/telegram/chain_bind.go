package telegram

import (
	"errors"
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

	err = interacter.Database.InsertChainBind(
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

	return interacter.TemplateManager.Render("chain_bind", chain)
}
