package telegram

import (
	"fmt"
	"html"
	"main/pkg/constants"
	"strconv"
	"strings"

	tele "gopkg.in/telebot.v3"
)

func (interacter *Interacter) GetSingleProposalCommand() Command {
	return Command{
		Name:    "proposal",
		Execute: interacter.HandleSingleProposal,
	}
}

func (interacter *Interacter) HandleSingleProposal(c tele.Context) (string, error) {
	chainBinds, err := interacter.Database.GetAllChainBinds(strconv.FormatInt(c.Chat().ID, 10))
	if err != nil {
		interacter.Logger.Error().Err(err).Msg("Error getting chain binds")
		return "", err
	}

	if len(chainBinds) == 1 {
		interacter.Logger.Debug().Msg("Single chain bound to a chat")
	} else {
		interacter.Logger.Debug().
			Strs("chains", chainBinds).
			Msg("Multiple or no chain bound to a chat")
	}

	var chainName string
	var proposalID string

	args := strings.Split(c.Text(), " ")

	if len(args) == 3 {
		// call is like /proposal <chain name> <proposal ID>
		chainName = args[1]
		proposalID = args[2]
	} else if len(chainBinds) == 1 && len(args) == 2 {
		// 1 chain bound to a chat, call is like /proposal <proposal ID>
		proposalID = args[1]
		chainName = chainBinds[0]
	} else {
		// 0 or >=2 chains bound to a chat and there's not enough info from query
		// to understand which chain to query.
		return html.EscapeString(fmt.Sprintf(
			"Usage: %s [chain] <proposal ID>",
			args[0],
		)), constants.ErrWrongInvocation
	}

	proposalsInfo := interacter.DataFetcher.GetSingleProposal(chainName, proposalID)
	return interacter.TemplateManager.Render("proposal", proposalsInfo)
}
