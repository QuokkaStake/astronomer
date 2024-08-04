package telegram

import (
	"fmt"
	"main/pkg/constants"
	"main/pkg/types"
	"main/pkg/utils"
	"strings"

	tele "gopkg.in/telebot.v3"
)

func (interacter *Interacter) GetChainAddCommand() Command {
	return Command{
		Name:    "chain_add",
		Execute: interacter.HandleAddChain,
	}
}

func (interacter *Interacter) HandleAddChain(c tele.Context, chainBinds []string) (string, error) {
	args := strings.SplitN(c.Text(), " ", 2)
	if len(args) < 2 {
		return "Usage: /chain_add <params>", constants.ErrWrongInvocation
	}
	query := args[1]

	argsAsMap, valid := utils.ParseArgsAsMap(query)
	if !valid {
		return "Invalid input syntax!", constants.ErrWrongInvocation
	}

	chain := &types.Chain{}

	for key, value := range argsAsMap {
		switch key {
		case "name":
			chain.Name = value
		case "lcd_endpoint":
			chain.LCDEndpoint = value
		case "lcd-endpoint":
			chain.LCDEndpoint = value
		case "pretty_name":
			chain.PrettyName = value
		case "pretty-name":
			chain.PrettyName = value
		}
	}

	if err := chain.Validate(); err != nil {
		return fmt.Sprintf("Invalid data provided: %s", err.Error()), err
	}

	err := interacter.Database.InsertChain(chain)
	if err != nil {
		if strings.Contains(err.Error(), "duplicate key value") {
			return "This chain is already inserted!", err
		}

		return "", err
	}

	return "Successfully added a new chain!", nil
}
