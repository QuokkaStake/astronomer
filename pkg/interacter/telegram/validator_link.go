package telegram

import (
	"fmt"
	"main/pkg/constants"
	"main/pkg/types"
	"strconv"
	"strings"

	tele "gopkg.in/telebot.v3"
)

func (interacter *Interacter) GetValidatorLinkCommand() Command {
	return Command{
		Name:    "validator_link",
		Execute: interacter.HandleValidatorLinkCommand,
	}
}

func (interacter *Interacter) HandleValidatorLinkCommand(c tele.Context, chainBinds []string) (string, error) {
	valid, usage, args := interacter.SingleChainItemParser(c.Text(), chainBinds, "address")
	if !valid {
		return usage, constants.ErrWrongInvocation
	}

	chains, err := interacter.Database.GetChainsByNames([]string{args.ChainName})
	if err != nil {
		return "", err
	} else if len(chains) < 1 {
		return interacter.ChainNotFound()
	}

	validator, err := interacter.DataFetcher.DoesValidatorExist(chains[0], args.ItemID)
	if err != nil {
		return fmt.Sprintf("Error linking validator: %s", err), err
	}

	err = interacter.Database.InsertValidatorLink(&types.ValidatorLink{
		Chain:    args.ChainName,
		Reporter: interacter.Name(),
		UserID:   strconv.FormatInt(c.Sender().ID, 10),
		Address:  args.ItemID,
	})
	if err != nil {
		if strings.Contains(err.Error(), "duplicate key value") {
			return "You have already linked this validator!", err
		}

		interacter.Logger.Error().Err(err).Msg("Error inserting validator link")
		return "", err
	}

	return interacter.TemplateManager.Render("validator_link", validator)
}
