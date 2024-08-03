package telegram

import (
	tele "gopkg.in/telebot.v3"
)

func (interacter *Interacter) GetActiveProposalsCommand() Command {
	return Command{
		Name:    "proposals",
		Execute: interacter.HandleActiveProposals,
	}
}

func (interacter *Interacter) HandleActiveProposals(c tele.Context, chainBinds []string) (string, error) {
	proposalsInfo := interacter.DataFetcher.GetActiveProposals(chainBinds)
	return interacter.TemplateManager.Render("proposals", proposalsInfo)
}
