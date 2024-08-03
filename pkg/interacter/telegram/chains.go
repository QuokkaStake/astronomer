package telegram

import (
	tele "gopkg.in/telebot.v3"
)

func (interacter *Interacter) GetChainsListCommand() Command {
	return Command{
		Name:    "chains",
		Execute: interacter.HandleChainsList,
	}
}

func (interacter *Interacter) HandleChainsList(c tele.Context, chainBinds []string) (string, error) {
	return interacter.TemplateManager.Render("chains", ChainsInfo{
		Chains:     interacter.Chains,
		ChainBinds: chainBinds,
	})
}
