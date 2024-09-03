package telegram

import (
	"strconv"

	tele "gopkg.in/telebot.v3"
)

func (interacter *Interacter) GetWalletsCommand() Command {
	return Command{
		Name:    "wallets",
		Execute: interacter.HandleWalletsCommand,
	}
}

func (interacter *Interacter) HandleWalletsCommand(c tele.Context, chainBinds []string) (string, error) {
	wallets := interacter.DataFetcher.GetWallets(strconv.FormatInt(c.Sender().ID, 10), interacter.Name())
	return interacter.TemplateManager.Render("wallets", wallets)
}
