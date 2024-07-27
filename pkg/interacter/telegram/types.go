package telegram

import (
	"main/pkg/types"

	tele "gopkg.in/telebot.v3"
)

type Command struct {
	Name    string
	Execute func(c tele.Context) (string, error)
	MinArgs int
	Usage   string
}

type ChainsInfo struct {
	Chains     []*types.Chain
	ChainBinds []string
}
