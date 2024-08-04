package telegram

import (
	"main/pkg/constants"

	tele "gopkg.in/telebot.v3"
)

func (interacter *Interacter) GetParamsCommand() Command {
	return Command{
		Name:    "params",
		Execute: interacter.HandleParams,
	}
}

func (interacter *Interacter) HandleParams(c tele.Context, chainBinds []string) (string, error) {
	valid, usage, args := interacter.BoundChainsNoArgsParser(c, chainBinds)
	if !valid {
		return usage, constants.ErrWrongInvocation
	}

	params := interacter.DataFetcher.GetChainsParams(args.ChainNames)
	return interacter.TemplateManager.Render("params", params)
}
