package telegram

import (
	"fmt"
	"html"
	"main/pkg/constants"
	"strings"

	tele "gopkg.in/telebot.v3"
)

func (interacter *Interacter) GetDenomDeleteCommand() Command {
	return Command{
		Name:    "denom_delete",
		Execute: interacter.HandleDeleteDenom,
	}
}

func (interacter *Interacter) HandleDeleteDenom(c tele.Context, chainBinds []string) (string, error) {
	args := strings.Split(c.Text(), " ")
	if len(args) < 3 {
		return html.EscapeString(fmt.Sprintf("Usage: %s <chain name> <denom name>", args[0])), constants.ErrWrongInvocation
	}

	deleted, err := interacter.Database.DeleteDenom(args[1], args[2])
	if err != nil {
		return "", err
	}

	if !deleted {
		return "Denom was not found!", err
	}

	return "Successfully deleted denom!", nil
}
