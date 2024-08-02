package telegram

import (
	"strconv"

	tele "gopkg.in/telebot.v3"
)

func (interacter *Interacter) GetParamsCommand() Command {
	return Command{
		Name:         "params",
		Execute:      interacter.HandleParams,
		ValidateArgs: NoArgs(),
	}
}

func (interacter *Interacter) HandleParams(c tele.Context) (string, error) {
	chainBinds, err := interacter.Database.GetAllChainBinds(strconv.FormatInt(c.Chat().ID, 10))
	if err != nil {
		interacter.Logger.Error().Err(err).Msg("Error getting chain binds")
		return "", err
	}

	params := interacter.DataFetcher.GetChainsParams(chainBinds)
	return interacter.TemplateManager.Render("params", params)
}
