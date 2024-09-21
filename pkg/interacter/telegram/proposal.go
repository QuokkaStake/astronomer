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
	valid, usage, args := interacter.SingleChainItemParser(c.Text(), chainBinds, "proposal ID")
	if !valid {
		return usage, constants.ErrWrongInvocation
	}

	chains, err := interacter.Database.GetChainsByNames([]string{args.ChainName})
	if err != nil {
		return "", err
	} else if len(chains) < 1 {
		return interacter.ChainNotFound()
	}

	proposalsInfo := interacter.DataFetcher.GetSingleProposal(chains[0], args.ItemID)
	return interacter.TemplateManager.Render("proposal", proposalsInfo)
}
