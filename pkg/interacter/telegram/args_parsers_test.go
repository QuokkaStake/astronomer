package telegram

import (
	"html"
	"testing"

	"github.com/rs/zerolog"
	"github.com/stretchr/testify/require"
)

func TestSingleArgParser(t *testing.T) {
	t.Parallel()

	interacter := &Interacter{}

	valid1, usage1, args1 := interacter.SingleArgParser("/command", "chain")
	require.False(t, valid1)
	require.Equal(t, html.EscapeString("Usage: /command <chain>"), usage1)
	require.Empty(t, args1)

	valid2, usage2, args2 := interacter.SingleArgParser("/command chain", "chain")
	require.True(t, valid2)
	require.Empty(t, usage2)
	require.Equal(t, SingleArg{Value: "chain"}, args2)
}

func TestBoundChainSingleQueryParser(t *testing.T) {
	t.Parallel()

	interacter := &Interacter{
		Logger: zerolog.Nop(),
	}

	valid1, usage1, args1 := interacter.BoundChainSingleQueryParser("/command", []string{})
	require.False(t, valid1)
	require.Equal(t, html.EscapeString("Usage: /command [chain] <query>"), usage1)
	require.Empty(t, args1)

	valid2, usage2, args2 := interacter.BoundChainSingleQueryParser("/command chain query path", []string{})
	require.True(t, valid2)
	require.Empty(t, usage2)
	require.Equal(t, BoundChainSingleQuery{
		ChainNames: []string{"chain"},
		Query:      "query path",
	}, args2)

	valid3, usage3, args3 := interacter.BoundChainSingleQueryParser("/command query", []string{"chain1", "chain2"})
	require.True(t, valid3)
	require.Empty(t, usage3)
	require.Equal(t, BoundChainSingleQuery{
		ChainNames: []string{"chain1", "chain2"},
		Query:      "query",
	}, args3)

	valid4, usage4, args4 := interacter.BoundChainSingleQueryParser("/command", []string{"chain1"})
	require.False(t, valid4)
	require.Equal(t, html.EscapeString("Usage: /command [chain] <query>"), usage4)
	require.Empty(t, args4)
}

func TestBoundChainAlias(t *testing.T) {
	t.Parallel()

	interacter := &Interacter{
		Logger: zerolog.Nop(),
	}

	valid1, usage1, args1 := interacter.BoundChainAliasParser("/command", []string{"chain"})
	require.False(t, valid1)
	require.Equal(t, html.EscapeString("Usage: /command <address> <alias>"), usage1)
	require.Empty(t, args1)

	valid2, usage2, args2 := interacter.BoundChainAliasParser("/command address long alias", []string{"chain"})
	require.True(t, valid2)
	require.Empty(t, usage2)
	require.Equal(t, BoundChainAlias{
		ChainName: "chain",
		Value:     "address",
		Alias:     "long alias",
	}, args2)

	valid3, usage3, args3 := interacter.BoundChainAliasParser("/command chain1 address long alias", []string{"chain1", "chain2"})
	require.True(t, valid3)
	require.Empty(t, usage3)
	require.Equal(t, BoundChainAlias{
		ChainName: "chain1",
		Value:     "address",
		Alias:     "long alias",
	}, args3)

	valid4, usage4, args4 := interacter.BoundChainAliasParser("/command", []string{"chain1", "chain2"})
	require.False(t, valid4)
	require.Equal(t, html.EscapeString("Usage: /command <chain> <address> <alias>"), usage4)
	require.Empty(t, args4)
}

func TestSingleChainItem(t *testing.T) {
	t.Parallel()

	interacter := &Interacter{
		Logger: zerolog.Nop(),
	}

	valid1, usage1, args1 := interacter.SingleChainItemParser("/command", []string{}, "proposal")
	require.False(t, valid1)
	require.Equal(t, html.EscapeString("Usage: /command <chain> <proposal>"), usage1)
	require.Empty(t, args1)

	valid2, usage2, args2 := interacter.SingleChainItemParser("/command ID", []string{"chain"}, "proposal")
	require.True(t, valid2)
	require.Empty(t, usage2)
	require.Equal(t, SingleChainItemArgs{
		ChainName: "chain",
		ItemID:    "ID",
	}, args2)

	valid3, usage3, args3 := interacter.SingleChainItemParser("/command chain1 ID", []string{"chain"}, "proposal")
	require.True(t, valid3)
	require.Empty(t, usage3)
	require.Equal(t, SingleChainItemArgs{
		ChainName: "chain1",
		ItemID:    "ID",
	}, args3)

	valid4, usage4, args4 := interacter.SingleChainItemParser("/command chain1 ID", []string{"chain1", "chain2"}, "proposal")
	require.True(t, valid4)
	require.Empty(t, usage4)
	require.Equal(t, SingleChainItemArgs{
		ChainName: "chain1",
		ItemID:    "ID",
	}, args4)
}

func TestBoundChainsNoArgs(t *testing.T) {
	t.Parallel()

	interacter := &Interacter{
		Logger: zerolog.Nop(),
	}

	valid1, usage1, args1 := interacter.BoundChainsNoArgsParser("/command", []string{})
	require.False(t, valid1)
	require.Equal(t, html.EscapeString("Usage: /command [chain]"), usage1)
	require.Empty(t, args1)

	valid2, usage2, args2 := interacter.BoundChainsNoArgsParser("/command chain1,chain2", []string{"chain"})
	require.True(t, valid2)
	require.Empty(t, usage2)
	require.Equal(t, BoundChainsNoArgs{
		ChainNames: []string{"chain1", "chain2"},
	}, args2)

	valid3, usage3, args3 := interacter.BoundChainsNoArgsParser("/command", []string{"chain"})
	require.True(t, valid3)
	require.Empty(t, usage3)
	require.Equal(t, BoundChainsNoArgs{
		ChainNames: []string{"chain"},
	}, args3)
}
