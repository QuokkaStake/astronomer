package telegram

import (
	"fmt"
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
		return interacter.ChainNotFound()
	}

	if err := interacter.DataFetcher.DoesWalletExist(chains[0], args.Value); err != nil {
		return fmt.Sprintf("Error linking wallet: %s", err), err
	}

	walletLink := &types.WalletLink{
		Chain:    args.ChainName,
		Reporter: interacter.Name(),
		UserID:   strconv.FormatInt(c.Sender().ID, 10),
		Address:  args.Value,
		Alias:    null.StringFrom(args.Alias),
	}

	err = interacter.Database.InsertWalletLink(walletLink)
	if err != nil {
		if strings.Contains(err.Error(), "duplicate key value") {
			return "You have already linked this wallet!", err
		}

		interacter.Logger.Error().Err(err).Msg("Error inserting wallet link")
		return "", err
	}

	explorers, err := interacter.Database.GetExplorersByChains([]string{args.ChainName})
	if err != nil {
		return "Error fetching chains!", err
	}

	return interacter.TemplateManager.Render("wallet_link", types.ChainWallet{
		Chain:     chains[0],
		Explorers: explorers,
		Wallet:    walletLink,
	})
}
