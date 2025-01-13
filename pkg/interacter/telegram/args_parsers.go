package telegram

import (
	"fmt"
	"html"
	"strings"
)

type SingleArg struct {
	Value string
}

func (interacter *Interacter) SingleArgParser(
	query string,
	argumentName string,
) (bool, string, SingleArg) {
	args := strings.Split(query, " ")

	if len(args) < 2 {
		return false, html.EscapeString(fmt.Sprintf(
			"Usage: %s <%s>",
			args[0],
			argumentName,
		)), SingleArg{}
	}

	return true, "", SingleArg{Value: args[1]}
}

type BoundChainSingleQuery struct {
	ChainNames []string
	Query      string
}

// Args parser when the command is called with 1 argument (like ID, but with chains)
// on multiple chains (like validator search).
// How it can be called:
// - /command query params - if there is 1 chain bound to a chat
// - /command chain1,chain2,chain3 query params - if there are 0 or 2+ chains bound to a chat.

func (interacter *Interacter) BoundChainSingleQueryParser(
	query string,
	chainBinds []string,
) (bool, string, BoundChainSingleQuery) {
	if len(chainBinds) > 0 {
		interacter.Logger.Debug().
			Strs("chains", chainBinds).
			Msg("There are chains bound to this chat")
	} else {
		interacter.Logger.Debug().
			Msg("No chain bound to this chat")
	}

	args := strings.SplitN(query, " ", 3)

	if len(args) == 3 {
		return true, "", BoundChainSingleQuery{ChainNames: strings.Split(args[1], ","), Query: args[2]}
	} else if len(chainBinds) > 0 && len(args) == 2 {
		return true, "", BoundChainSingleQuery{ChainNames: chainBinds, Query: args[1]}
	} else {
		return false, html.EscapeString(fmt.Sprintf(
			"Usage: %s [chain] <query>",
			args[0],
		)), BoundChainSingleQuery{}
	}
}

// Args parser when the command is called with 2 arguments, like a wallet address and an alias.
// How it can be called:
// - /command wallet alias - if there is 1 chain bound to a chat
// - /command chain wallet alias - if there are 0 or 2+ chains bound to a chat.

type BoundChainAlias struct {
	ChainName string
	Value     string
	Alias     string
}

func (interacter *Interacter) BoundChainAliasParser(
	query string,
	chainBinds []string,
) (bool, string, BoundChainAlias) {
	if len(chainBinds) == 1 { //nolint:nestif
		args := strings.SplitN(query, " ", 3)

		interacter.Logger.Debug().Msg("Single chain bound to a chat")

		if len(args) == 3 { // /command address alias
			return true, "", BoundChainAlias{
				ChainName: chainBinds[0],
				Value:     args[1],
				Alias:     args[2],
			}
		} else { // invalid
			return false, html.EscapeString(fmt.Sprintf(
				"Usage: %s <address> <alias>",
				args[0],
			)), BoundChainAlias{}
		}
	} else {
		args := strings.SplitN(query, " ", 4)

		interacter.Logger.Debug().
			Strs("chains", chainBinds).
			Msg("Zero or multiple or no chain bound to a chat")

		if len(args) == 4 { // /command chain address alias
			return true, "", BoundChainAlias{
				ChainName: args[1],
				Value:     args[2],
				Alias:     args[3],
			}
		} else { // invalid
			return false, html.EscapeString(fmt.Sprintf(
				"Usage: %s <chain> <address> <alias>",
				args[0],
			)), BoundChainAlias{}
		}
	}
}

type SingleChainItemArgs struct {
	ChainName string
	ItemID    string
}

// Args parser when the command is called with 1 argument (like ID, but with chains)
// on 1 chain (like, a proposal on a specific chain).
// How it can be called:
// - /command ID - if there's exactly 1 chain bound to a chat
// - /command chain_name ID - if there's 0 or 2+ more chains bound to a chat

func (interacter *Interacter) SingleChainItemParser(
	query string,
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

	args := strings.Split(query, " ")

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
			"Usage: %s <chain> <%s>",
			args[0],
			argumentName,
		)), SingleChainItemArgs{}
	}
}

type BoundChainsNoArgs struct {
	ChainNames []string
}

// Args parser when the command is called without arguments, but with chains.
// How it can be called:
// - /command - if there are chains bound to a chat
// - /command chain1,chain2 - in any case

func (interacter *Interacter) BoundChainsNoArgsParser(
	query string,
	chainBinds []string,
) (bool, string, BoundChainsNoArgs) {
	if len(chainBinds) == 1 {
		interacter.Logger.Debug().Msg("Single chain bound to a chat")
	} else {
		interacter.Logger.Debug().
			Strs("chains", chainBinds).
			Msg("Multiple or no chain bound to a chat")
	}

	args := strings.SplitN(query, " ", 2)

	if len(args) == 2 {
		// call is like /command <chain name>
		return true, "", BoundChainsNoArgs{ChainNames: strings.Split(args[1], ",")}
	} else if len(chainBinds) > 0 && len(args) == 1 {
		// 1 chain bound to a chat, call is like /command
		return true, "", BoundChainsNoArgs{ChainNames: chainBinds}
	} else {
		// No chains bound to a chat and there's not enough info from query
		// to understand which chain to query.
		return false, html.EscapeString(fmt.Sprintf(
			"Usage: %s [chain]",
			args[0],
		)), BoundChainsNoArgs{}
	}
}
