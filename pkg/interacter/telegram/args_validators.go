package telegram

import tele "gopkg.in/telebot.v3"

type ArgsValidator func(c tele.Context) (bool, string)

func NoArgs() ArgsValidator {
	return func(c tele.Context) (bool, string) {
		return true, ""
	}
}

func MinArgs(args int, usage string) ArgsValidator {
	return func(c tele.Context) (bool, string) {
		if len(c.Args()) < args {
			return false, usage
		}

		return true, ""
	}
}
