package telegram

import (
	"fmt"
	"html"
	"main/pkg/constants"
	"main/pkg/types"
	"main/pkg/utils"
	"strings"

	tele "gopkg.in/telebot.v3"
)

func (interacter *Interacter) GetChainUpdateCommand() Command {
	return Command{
		Name:    "chain_update",
		Execute: interacter.HandleUpdateChain,
	}
}

func (interacter *Interacter) HandleUpdateChain(c tele.Context, chainBinds []string) (string, error) {
	args := strings.SplitN(c.Text(), " ", 2)
	if len(args) < 2 {
		return html.EscapeString(fmt.Sprintf("Usage: %s <params>", args[0])), constants.ErrWrongInvocation
	}
	query := args[1]

	argsAsMap, valid := utils.ParseArgsAsMap(query)
	if !valid {
		return "Invalid input syntax!", constants.ErrWrongInvocation
	}

	chain := types.ChainFromArgs(argsAsMap)
	if err := chain.Validate(); err != nil {
		return fmt.Sprintf("Invalid data provided: %s", err.Error()), err
	}

	updated, err := interacter.Database.UpdateChain(chain)
	if err != nil {
		return "", err
	}

	if !updated {
		return "Chain was not found!", err
	}

	return "Successfully updated chain!", nil
}
