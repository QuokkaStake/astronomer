package types

import (
	"fmt"
)

type ValidatorLink struct {
	Chain    string
	Reporter string
	UserID   string
	Address  string
}

func (l *ValidatorLink) Validate() error {
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
