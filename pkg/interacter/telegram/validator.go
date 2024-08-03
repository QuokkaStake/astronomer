package telegram

import (
	"main/pkg/constants"

	tele "gopkg.in/telebot.v3"
)

func (interacter *Interacter) GetValidatorCommand() Command {
	return Command{
		Name:    "validator",
		Execute: interacter.HandleValidator,
	}
}

func (interacter *Interacter) HandleValidator(c tele.Context, chainBinds []string) (string, error) {
	valid, usage, args := interacter.SingleQueryParser(c)
	if !valid {
		return usage, constants.ErrWrongInvocation
	}

	validatorsInfo := interacter.DataFetcher.FindValidator(args.Value, chainBinds)

	return interacter.TemplateManager.Render("validator", validatorsInfo)
}
