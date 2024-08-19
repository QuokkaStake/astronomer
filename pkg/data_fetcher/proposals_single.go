package datafetcher

import (
	"fmt"
	"main/pkg/tendermint"
	"main/pkg/types"
)

func (f *DataFetcher) GetSingleProposal(chainName string, proposalID string) types.SingleProposal {
	response := types.SingleProposal{}

	chains, err := f.Database.GetChainsByNames([]string{chainName})
	if err != nil {
		response.Error = err
		return response
	}

	explorers, err := f.Database.GetExplorersByChains([]string{chainName})
	if err != nil {
		response.Error = err
		return response
	}

	if len(chains) != 1 {
		response.Error = fmt.Errorf("chain '%s' is not found", chainName)
		return response
	}

	response.Chain = chains[0]
	response.Explorers = explorers.GetExplorersByChain(chainName)

	rpc := tendermint.NewRPC(chains[0], 10, f.Logger)
	proposal, _, err := rpc.GetSingleProposal(proposalID)

	if err != nil {
		response.Error = err
	} else {
		response.Proposal = proposal
	}

	return response
}
