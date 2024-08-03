package telegram

import (
	"main/pkg/constants"

	tele "gopkg.in/telebot.v3"
)

func (interacter *Interacter) GetSingleProposalCommand() Command {
	return Command{
		Name:    "proposal",
		Execute: interacter.HandleSingleProposal,
	}
}

func (interacter *Interacter) HandleSingleProposal(
	c tele.Context,
	chainBinds []string,
) (string, error) {
	valid, usage, args := interacter.SingleChainItemParser(c, chainBinds, "proposal ID")
	if !valid {
		return usage, constants.ErrWrongInvocation
	}

	proposalsInfo := interacter.DataFetcher.GetSingleProposal(args.ChainName, args.ItemID)
	return interacter.TemplateManager.Render("proposal", proposalsInfo)
}
