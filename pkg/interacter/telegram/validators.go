package telegram

import (
	"main/pkg/constants"
	"strconv"

	tele "gopkg.in/telebot.v3"
)

func (interacter *Interacter) GetValidatorsCommand() Command {
	return Command{
		Name:    "validators",
		Execute: interacter.HandleValidators,
	}
}

func (interacter *Interacter) HandleValidators(c tele.Context, chainBinds []string) (string, error) {
	valid, usage, args := interacter.BoundChainsNoArgsParser(c.Text(), chainBinds)
	if !valid {
		return usage, constants.ErrWrongInvocation
	}

	validatorsInfo := interacter.DataFetcher.FindMyValidators(
		args.ChainNames,
		strconv.FormatInt(c.Sender().ID, 10),
		interacter.Name(),
	)
	return interacter.TemplateManager.Render("validators", validatorsInfo)
}
