package telegram

import (
	"errors"
	"fmt"
	"html"
	"main/pkg/constants"
	"main/pkg/types"
	"strings"

	tele "gopkg.in/telebot.v3"
)

func (interacter *Interacter) GetChainInfoCommand() Command {
	return Command{
		Name:    "chain",
		Execute: interacter.HandleChainInfo,
	}
}

func (interacter *Interacter) HandleChainInfo(c tele.Context, chainBinds []string) (string, error) {
	args := strings.SplitN(c.Text(), " ", 2)
	if len(args) < 2 {
		return html.EscapeString(fmt.Sprintf("Usage: %s <chain name>", args[0])), constants.ErrWrongInvocation
	}

	chain, err := interacter.Database.GetChainByName(args[1])
	if err != nil && errors.Is(err, constants.ErrChainNotFound) {
		return interacter.ChainNotFound()
	} else if err != nil {
		return "", err
	}

	explorers, err := interacter.Database.GetExplorersByChains([]string{chain.Name})
	if err != nil {
		return "Error getting explorers!", err
	}

	denoms, err := interacter.Database.GetDenomsByChain(chain)
	if err != nil {
		return "Error getting denoms!", err
	}

	lcds, err := interacter.Database.GetLCDHosts(chain)
	if err != nil {
		return "Error getting LCD hosts!", err
	}

	return interacter.TemplateManager.Render("chain", &types.ChainInfo{
		Chain:        chain,
		Explorers:    explorers,
		Denoms:       denoms,
		LCDEndpoints: lcds,
	})
}
