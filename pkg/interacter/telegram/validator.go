package telegram

import (
	"strconv"
	"strings"

	tele "gopkg.in/telebot.v3"
)

func (interacter *Interacter) GetValidatorCommand() Command {
	return Command{
		Name:         "validator",
		Execute:      interacter.HandleValidator,
		ValidateArgs: MinArgs(1, "Usage: %s <query>"),
	}
}

func (interacter *Interacter) HandleValidator(c tele.Context) (string, error) {
	chainBinds, err := interacter.Database.GetAllChainBinds(strconv.FormatInt(c.Chat().ID, 10))
	if err != nil {
		interacter.Logger.Error().Err(err).Msg("Error getting chain binds")
		return "", err
	}

	args := strings.SplitN(c.Text(), " ", 2)
	validatorsInfo := interacter.DataFetcher.FindValidator(args[1], chainBinds)

	return interacter.TemplateManager.Render("validator", validatorsInfo)
}
