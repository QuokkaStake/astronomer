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

func (interacter *Interacter) GetExplorerAddCommand() Command {
	return Command{
		Name:    "explorer_add",
		Execute: interacter.HandleAddExplorer,
	}
}

func (interacter *Interacter) HandleAddExplorer(c tele.Context, chainBinds []string) (string, error) {
	args := strings.SplitN(c.Text(), " ", 2)
	if len(args) < 2 {
		return html.EscapeString(fmt.Sprintf("Usage: %s <params>", args[0])), constants.ErrWrongInvocation
	}
	query := args[1]

	argsAsMap, valid := utils.ParseArgsAsMap(query)
	if !valid {
		return "Invalid input syntax!", constants.ErrWrongInvocation
	}

	explorer := types.ExplorerFromArgs(argsAsMap)

	if err := explorer.Validate(); err != nil {
		return fmt.Sprintf("Invalid data provided: %s", err.Error()), err
	}

	err := interacter.Database.InsertExplorer(explorer)
	if err != nil {
		if strings.Contains(err.Error(), "duplicate key value") {
			return "This explorer is already inserted!", err
		}

		return "", err
	}

	return interacter.TemplateManager.Render("explorer_add", explorer)
}
