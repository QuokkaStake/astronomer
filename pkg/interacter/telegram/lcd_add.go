package telegram

import (
	"fmt"
	"html"
	"main/pkg/constants"
	"main/pkg/types"
	"strings"

	tele "gopkg.in/telebot.v3"
)

func (interacter *Interacter) GetLCDAddCommand() Command {
	return Command{
		Name:    "lcd_add",
		Execute: interacter.HandleAddLCD,
	}
}

func (interacter *Interacter) HandleAddLCD(c tele.Context, chainBinds []string) (string, error) {
	args := strings.SplitN(c.Text(), " ", 3)
	if len(args) < 3 {
		return html.EscapeString(fmt.Sprintf("Usage: %s <chain name> <host>", args[0])), constants.ErrWrongInvocation
	}

	chainName, host := args[1], args[2]
	chain, err := interacter.Database.GetChainByName(chainName)
	if err != nil {
		return "Error finding chain!", err
	}

	if insertErr := interacter.Database.InsertLCDHost(chain, host); insertErr != nil {
		return "Error inserting LCD host!", insertErr
	}

	return interacter.TemplateManager.Render("lcd_add", types.ChainWithLCD{Chain: *chain, LCDEndpoint: host})
}
