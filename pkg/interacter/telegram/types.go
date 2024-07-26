package telegram

import tele "gopkg.in/telebot.v3"

type Command struct {
	Name    string
	Execute func(c tele.Context) (string, error)
}
