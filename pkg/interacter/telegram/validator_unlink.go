package telegram

import (
	"fmt"
	"html"
	"main/pkg/constants"
	"strconv"
	"strings"

	tele "gopkg.in/telebot.v3"
)

func (interacter *Interacter) GetValidatorUnlinkCommand() Command {
	return Command{
		Name:    "validator_unlink",
		Execute: interacter.HandleValidatorUnlink,
	}
}

func (interacter *Interacter) HandleValidatorUnlink(c tele.Context, chainBinds []string) (string, error) {
	args := strings.Split(c.Text(), " ")
	if len(args) < 3 {
		return html.EscapeString(fmt.Sprintf("Usage: %s <chain name> <address>", args[0])), constants.ErrWrongInvocation
	}

	deleted, err := interacter.Database.DeleteValidatorLink(args[1], interacter.Name(), args[2], strconv.FormatInt(c.Sender().ID, 10))
	if err != nil {
		return "", err
	}

	if !deleted {
		return "Validator was not linked!", err
	}

	return "Successfully unlinked a validator!", nil
}
