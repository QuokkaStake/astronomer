package types

import (
	"fmt"

	"github.com/guregu/null/v5"
)

type WalletLink struct {
	Chain    string
	Reporter string
	UserID   string
	Address  string
	Alias    null.String
}

func (l *WalletLink) Validate() error {
	if l.Chain == "" {
		return fmt.Errorf("empty chain name")
	}

	if l.Reporter == "" {
		return fmt.Errorf("empty reporter")
	}

	if l.UserID == "" {
		return fmt.Errorf("empty user")
	}

	if l.Address == "" {
		return fmt.Errorf("empty address")
	}

	return nil
}
