package datafetcher

import (
	"fmt"
	"main/pkg/types"
)

func (f *DataFetcher) DoesWalletExist(chain *types.Chain, wallet string) error {
	balances, err := f.NodesManager.GetBalance(chain, wallet)
	if err != nil {
		return err
	}

	if len(balances.Balances) == 0 {
		return fmt.Errorf("wallet %s does not exist", wallet)
	}

	return nil
}
