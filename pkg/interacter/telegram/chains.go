package telegram

import (
	"main/pkg/types"
	"main/pkg/utils"

	tele "gopkg.in/telebot.v3"
)

func (interacter *Interacter) GetChainsListCommand() Command {
	return Command{
		Name:    "chains",
		Execute: interacter.HandleChainsList,
	}
}

func (interacter *Interacter) HandleChainsList(c tele.Context, chainBinds []string) (string, error) {
	chains, err := interacter.Database.GetAllChains()
	if err != nil {
		return "Error fetching chains!", err
	}

	chainNames := utils.Map(chains, func(c *types.Chain) string {
		return c.Name
	})

	explorers, err := interacter.Database.GetExplorersByChains(chainNames)
	if err != nil {
		return "Error fetching chains!", err
	}

	return interacter.TemplateManager.Render("chains", ChainsInfo{
		Chains:     chains,
		ChainBinds: chainBinds,
		Explorers:  explorers,
	})
}
