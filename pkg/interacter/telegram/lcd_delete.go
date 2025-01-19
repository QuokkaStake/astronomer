package telegram

import (
	"fmt"
	"html"
	"main/pkg/constants"
	"main/pkg/types"
	"strings"

	tele "gopkg.in/telebot.v3"
)

func (interacter *Interacter) GetLCDDeleteCommand() Command {
	return Command{
		Name:    "lcd_delete",
		Execute: interacter.HandleDeleteLCD,
	}
}

func (interacter *Interacter) HandleDeleteLCD(c tele.Context, chainBinds []string) (string, error) {
	args := strings.SplitN(c.Text(), " ", 3)
	if len(args) < 3 {
		return html.EscapeString(fmt.Sprintf("Usage: %s <chain name> <host>", args[0])), constants.ErrWrongInvocation
	}

	chainName, host := args[1], args[2]
	chain, err := interacter.Database.GetChainByName(chainName)
	if err != nil {
		return "Error finding chain!", err
	}

	allLCDs, err := interacter.Database.GetLCDHosts(chain)
	if err != nil {
		return "Error finding LCD hosts!", err
	}

	if len(allLCDs) <= 1 && host == allLCDs[0] {
		return "Cannot remove the only chain LCD!", constants.ErrWrongInvocation
	}

	deleted, deleteErr := interacter.Database.DeleteLCDHost(chain, host)
	if deleteErr != nil {
		return "Error deleting LCD host!", deleteErr
	}

	if !deleted {
		return "Chain LCD host was not found!", constants.ErrLCDNotFound
	}

	return interacter.TemplateManager.Render("lcd_delete", types.ChainWithLCD{
		Chain:       *chain,
		LCDEndpoint: host,
	})
}
