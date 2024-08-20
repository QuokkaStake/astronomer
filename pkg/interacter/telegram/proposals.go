package telegram

import (
	"main/pkg/constants"

	tele "gopkg.in/telebot.v3"
)

func (interacter *Interacter) GetActiveProposalsCommand() Command {
	return Command{
		Name:    "proposals",
		Execute: interacter.HandleActiveProposals,
	}
}

func (interacter *Interacter) HandleActiveProposals(c tele.Context, chainBinds []string) (string, error) {
	valid, usage, args := interacter.BoundChainsNoArgsParser(c.Text(), chainBinds)
	if !valid {
		return usage, constants.ErrWrongInvocation
	}

	proposalsInfo := interacter.DataFetcher.GetActiveProposals(args.ChainNames)
	return interacter.TemplateManager.Render("proposals", proposalsInfo)
}
