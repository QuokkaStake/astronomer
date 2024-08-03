package telegram

import (
	tele "gopkg.in/telebot.v3"
)

func (interacter *Interacter) GetParamsCommand() Command {
	return Command{
		Name:    "params",
		Execute: interacter.HandleParams,
	}
}

func (interacter *Interacter) HandleParams(c tele.Context, chainBinds []string) (string, error) {
	params := interacter.DataFetcher.GetChainsParams(chainBinds)
	return interacter.TemplateManager.Render("params", params)
}
