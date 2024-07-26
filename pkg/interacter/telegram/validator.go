package telegram

import (
	"fmt"
	"strconv"
	"strings"

	tele "gopkg.in/telebot.v3"
)

func (interacter *Interacter) GetValidatorCommand() Command {
	return Command{
		Name:    "validator",
		Execute: interacter.HandleValidator,
	}
}

func (interacter *Interacter) HandleValidator(c tele.Context) (string, error) {
	chainBinds, err := interacter.Database.GetAllChainBinds(strconv.FormatInt(c.Chat().ID, 10))
	if err != nil {
		interacter.Logger.Error().Err(err).Msg("Error getting chain binds")
		return "", err
	}

	args := strings.SplitN(c.Text(), " ", 2)
	if len(args) < 2 {
		return fmt.Sprintf("Usage: %s <query>", args[0]), fmt.Errorf("invalid command invocation")
	}

	validatorsInfo := interacter.DataFetcher.FindValidator(args[1], chainBinds)

	template, err := interacter.TemplateManager.Render("validator", validatorsInfo)
	if err != nil {
		interacter.Logger.Error().Err(err).Msg("Error rendering template")
		return "", err
	}

	return template, nil
}
