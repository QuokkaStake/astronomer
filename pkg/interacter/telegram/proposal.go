package telegram

import (
	"errors"
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

	chain, err := interacter.Database.GetChainByName(args.ChainName)
	if err != nil && errors.Is(err, constants.ErrChainNotFound) {
		return interacter.ChainNotFound()
	} else if err != nil {
		return "", err
	}

	proposalsInfo := interacter.DataFetcher.GetSingleProposal(chain, args.ItemID)
	return interacter.TemplateManager.Render("proposal", proposalsInfo)
}
