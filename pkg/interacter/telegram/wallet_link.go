package telegram

import (
	"fmt"
	"html"
	"main/pkg/constants"
	"main/pkg/types"
	"strconv"
	"strings"

	"github.com/guregu/null/v5"

	tele "gopkg.in/telebot.v3"
)

func (interacter *Interacter) GetWalletLinkCommand() Command {
	return Command{
		Name:    "wallet_link",
		Execute: interacter.HandleWalletLinkCommand,
	}
}

func (interacter *Interacter) HandleWalletLinkCommand(c tele.Context, chainBinds []string) (string, error) {
	valid, usage, args := interacter.BoundChainAliasParser(c.Text(), chainBinds)
	if !valid {
		return usage, constants.ErrWrongInvocation
	}

	chains, err := interacter.Database.GetChainsByNames([]string{args.ChainName})
	if err != nil {
		return "", err
	} else if len(chains) < 1 {
		return html.EscapeString(fmt.Sprintf(
			"Could not find a chain with the name '%s'",
			args.ChainName,
		)), constants.ErrChainNotFound
	}

	err = interacter.Database.InsertWalletLink(&types.WalletLink{
		Chain:    args.ChainName,
		Reporter: interacter.Name(),
		UserID:   strconv.FormatInt(c.Sender().ID, 10),
		Address:  args.Value,
		Alias:    null.StringFrom(args.Alias),
	})
	if err != nil {
		if strings.Contains(err.Error(), "duplicate key value") {
			return "You have already linked this wallet!", err
		}

		interacter.Logger.Error().Err(err).Msg("Error inserting wallet link")
		return "", err
	}

	return "Successfully linked a wallet!", nil
}
