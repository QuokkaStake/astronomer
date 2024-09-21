package telegram

import (
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

	chains, err := interacter.Database.GetChainsByNames([]string{args.Value})
	if err != nil {
		return "", err
	} else if len(chains) < 1 {
		return interacter.ChainNotFound()
	}

	err = interacter.Database.InsertChainBind(
		interacter.Name(),
		strconv.FormatInt(c.Chat().ID, 10),
		c.Chat().Title,
		chains[0].Name,
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
