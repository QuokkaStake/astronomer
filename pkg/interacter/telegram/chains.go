package telegram

import (
	"strconv"

	tele "gopkg.in/telebot.v3"
)

func (interacter *Interacter) GetChainsListCommand() Command {
	return Command{
		Name:    "chains",
		Execute: interacter.HandleChainsList,
	}
}

func (interacter *Interacter) HandleChainsList(c tele.Context) (string, error) {
	chainBinds, err := interacter.Database.GetAllChainBinds(strconv.FormatInt(c.Chat().ID, 10))
	if err != nil {
		interacter.Logger.Error().Err(err).Msg("Error getting chain binds")
		return "", err
	}

	return interacter.TemplateManager.Render("chains", ChainsInfo{
		Chains:     interacter.Chains,
		ChainBinds: chainBinds,
	})
}
