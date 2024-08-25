package telegram

import (
	"strconv"

	tele "gopkg.in/telebot.v3"
)

func (interacter *Interacter) GetBalanceCommand() Command {
	return Command{
		Name:    "balance",
		Execute: interacter.HandleBalanceCommand,
	}
}

func (interacter *Interacter) HandleBalanceCommand(c tele.Context, chainBinds []string) (string, error) {
	balances := interacter.DataFetcher.GetBalances(strconv.FormatInt(c.Sender().ID, 10), interacter.Name())
	return interacter.TemplateManager.Render("balance", balances)
}
