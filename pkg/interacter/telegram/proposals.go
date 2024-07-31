package telegram

import (
	"strconv"

	tele "gopkg.in/telebot.v3"
)

func (interacter *Interacter) GetActiveProposalsCommand() Command {
	return Command{
		Name:    "proposals",
		Execute: interacter.HandleActiveProposals,
	}
}

func (interacter *Interacter) HandleActiveProposals(c tele.Context) (string, error) {
	chainBinds, err := interacter.Database.GetAllChainBinds(strconv.FormatInt(c.Chat().ID, 10))
	if err != nil {
		interacter.Logger.Error().Err(err).Msg("Error getting chain binds")
		return "", err
	}

	proposalsInfo := interacter.DataFetcher.GetActiveProposals(chainBinds)
	return interacter.TemplateManager.Render("proposals", proposalsInfo)
}
