package datafetcher

import (
	"main/pkg/types"
)

func (f *DataFetcher) GetSingleProposal(chain *types.Chain, proposalID string) types.SingleProposal {
	response := types.SingleProposal{}

	explorers, err := f.Database.GetExplorersByChains([]string{chain.Name})
	if err != nil {
		response.Error = err
		return response
	}

	response.Chain = chain
	response.Explorers = explorers.GetExplorersByChain(chain.Name)

	proposal, err := f.NodesManager.GetSingleProposal(chain, proposalID)

	if err != nil {
		response.Error = err
	} else {
		response.Proposal = proposal
	}

	return response
}
