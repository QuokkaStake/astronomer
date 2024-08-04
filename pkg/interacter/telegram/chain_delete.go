package telegram

import (
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
		return "Usage: /chain_delete <chain name>", constants.ErrWrongInvocation
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
