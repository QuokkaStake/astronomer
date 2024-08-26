package telegram

import (
	tele "gopkg.in/telebot.v3"
)

type HelpData struct {
	Version string
	Chains  []string
}

func (h HelpData) HasOneChain() bool {
	return len(h.Chains) == 1
}

func (interacter *Interacter) GetHelpCommand() Command {
	return Command{
		Name:    "help",
		Execute: interacter.HandleHelpCommand,
	}
}

func (interacter *Interacter) HandleHelpCommand(_ tele.Context, chainBinds []string) (string, error) {
	return interacter.TemplateManager.Render("help", HelpData{
		Version: interacter.Version,
		Chains:  chainBinds,
	})
}
