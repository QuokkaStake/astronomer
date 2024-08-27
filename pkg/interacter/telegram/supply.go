package telegram

import (
	"main/pkg/constants"

	tele "gopkg.in/telebot.v3"
)

func (interacter *Interacter) GetSupplyCommand() Command {
	return Command{
		Name:    "supply",
		Execute: interacter.HandleSupply,
	}
}

func (interacter *Interacter) HandleSupply(c tele.Context, chainBinds []string) (string, error) {
	valid, usage, args := interacter.BoundChainsNoArgsParser(c.Text(), chainBinds)
	if !valid {
		return usage, constants.ErrWrongInvocation
	}

	supply := interacter.DataFetcher.GetSupply(args.ChainNames)
	return interacter.TemplateManager.Render("supply", supply)
}
