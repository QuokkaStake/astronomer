package telegram

import (
	"fmt"
	"html"
	"main/pkg/constants"
	"strings"

	tele "gopkg.in/telebot.v3"
)

func (interacter *Interacter) GetExplorerDeleteCommand() Command {
	return Command{
		Name:    "chain_delete",
		Execute: interacter.HandleDeleteExplorer,
	}
}

func (interacter *Interacter) HandleDeleteExplorer(c tele.Context, chainBinds []string) (string, error) {
	args := strings.Split(c.Text(), " ")
	if len(args) < 3 {
		return html.EscapeString(fmt.Sprintf("Usage: %s <chain name> <explorer name>", args[0])), constants.ErrWrongInvocation
	}

	deleted, err := interacter.Database.DeleteExplorer(args[1], args[2])
	if err != nil {
		return "", err
	}

	if !deleted {
		return "Explorer was not found!", err
	}

	return "Successfully deleted explorer!", nil
}
