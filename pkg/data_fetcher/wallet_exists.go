package datafetcher

import (
	"fmt"
	"main/pkg/types"
)

func (f *DataFetcher) DoesWalletExist(chain *types.Chain, wallet string) error {
	rpc := f.GetRPC(chain)

	balances, _, err := rpc.GetBalance(wallet)
	if err != nil {
		return err
	}

	if len(balances.Balances) == 0 {
		return fmt.Errorf("wallet %s does not exist", wallet)
	}

	return nil
}
