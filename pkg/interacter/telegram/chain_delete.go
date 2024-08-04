package telegram

import (
	"fmt"
	"html"
	"main/pkg/constants"
	"strings"

	tele "gopkg.in/telebot.v3"
)

func (interacter *Interacter) GetChainDeleteCommand() Command {
	return Command{
		Name:    "chain_delete",
		Execute: interacter.HandleDeleteChain,
	}
}

func (interacter *Interacter) HandleDeleteChain(c tele.Context, chainBinds []string) (string, error) {
	args := strings.Split(c.Text(), " ")
	if len(args) < 2 {
		return html.EscapeString(fmt.Sprintf("Usage: %s <chain name>", args[0])), constants.ErrWrongInvocation
	}

	deleted, err := interacter.Database.DeleteChain(args[1])
	if err != nil {
		return "", err
	}

	if !deleted {
		return "Chain was not found!", err
	}

	return "Successfully deleted chain!", nil
}
