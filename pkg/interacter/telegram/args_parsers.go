package telegram

import (
	"fmt"
	"html"
	"strings"

	tele "gopkg.in/telebot.v3"
)

type SingleArg struct {
	Value string
}

func (interacter *Interacter) SingleArgParser(
	c tele.Context,
	argumentName string,
) (bool, string, SingleArg) {
	args := strings.Split(c.Text(), " ")

	if len(args) < 2 {
		return false, html.EscapeString(fmt.Sprintf(
			"Usage: %s <%s>",
			args[0],
			argumentName,
		)), SingleArg{}
	}

	return true, "", SingleArg{Value: args[1]}
}

type SingleQuery struct {
	Value string
}

func (interacter *Interacter) SingleQueryParser(
	c tele.Context,
) (bool, string, SingleArg) {
	args := strings.SplitN(c.Text(), " ", 2)

	if len(args) < 2 {
		return false, html.EscapeString(fmt.Sprintf(
			"Usage: %s <query>",
			args[0],
		)), SingleArg{}
	}

	return true, "", SingleArg{Value: args[1]}
}

type SingleChainItemArgs struct {
	ChainName string
	ItemID    string
}

// Args parser when the command is called with 1 argument (like ID, but with chains).
// How it can be called:
// - /command ID - if there's exactly 1 chain bound to a chat
// - /command chain_name ID - if there's 0 or 2+ more chains bound to a chat

func (interacter *Interacter) SingleChainItemParser(
	c tele.Context,
	chainBinds []string,
	argumentName string,
) (bool, string, SingleChainItemArgs) {
	if len(chainBinds) == 1 {
		interacter.Logger.Debug().Msg("Single chain bound to a chat")
	} else {
		interacter.Logger.Debug().
			Strs("chains", chainBinds).
			Msg("Multiple or no chain bound to a chat")
	}

	args := strings.Split(c.Text(), " ")

	if len(args) == 3 {
		// call is like /command <chain name> <proposal ID>
		return true, "", SingleChainItemArgs{ChainName: args[1], ItemID: args[2]}
	} else if len(chainBinds) == 1 && len(args) == 2 {
		// 1 chain bound to a chat, call is like /command <proposal ID>
		return true, "", SingleChainItemArgs{ChainName: chainBinds[0], ItemID: args[1]}
	} else {
		// 0 or >=2 chains bound to a chat and there's not enough info from query
		// to understand which chain to query.
		return false, html.EscapeString(fmt.Sprintf(
			"Usage: %s [chain] <%s>",
			args[0],
			argumentName,
		)), SingleChainItemArgs{}
	}
}

type SingleChainNoArgs struct {
	ChainNames []string
}

// Args parser when the command is called without arguments, but with chains.
// How it can be called:
// - /command - if there's exactly 1 chain bound to a chat
// - /command chain1,chain2 - if there's 0 or 2+ more chains bound to a chat

func (interacter *Interacter) SingleChainNoArgsParser(
	c tele.Context,
	chainBinds []string,
) (bool, string, SingleChainNoArgs) {
	if len(chainBinds) == 1 {
		interacter.Logger.Debug().Msg("Single chain bound to a chat")
	} else {
		interacter.Logger.Debug().
			Strs("chains", chainBinds).
			Msg("Multiple or no chain bound to a chat")
	}

	args := strings.SplitN(c.Text(), " ", 2)

	if len(args) == 2 {
		// call is like /command <chain name>
		return true, "", SingleChainNoArgs{ChainNames: strings.Split(args[1], ",")}
	} else if len(chainBinds) == 1 && len(args) == 1 {
		// 1 chain bound to a chat, call is like /command
		return true, "", SingleChainNoArgs{ChainNames: chainBinds}
	} else {
		// 0 or >=2 chains bound to a chat and there's not enough info from query
		// to understand which chain to query.
		return false, html.EscapeString(fmt.Sprintf(
			"Usage: %s [chain]",
			args[0],
		)), SingleChainNoArgs{}
	}
}
