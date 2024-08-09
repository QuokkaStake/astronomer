package telegram

import (
	"main/pkg/types"

	tele "gopkg.in/telebot.v3"
)

type Command struct {
	Name    string
	Execute func(c tele.Context, chainBinds []string) (string, error)
}

type ChainsInfo struct {
	Chains     []*types.Chain
	Explorers  types.Explorers
	ChainBinds []string
}
