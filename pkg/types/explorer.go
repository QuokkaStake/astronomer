package types

import (
	"errors"
	"fmt"
	"main/pkg/utils"
)

type Explorer struct {
	Chain                string
	Name                 string
	ProposalLinkPattern  string
	WalletLinkPattern    string
	ValidatorLinkPattern string
	MainLink             string
}

func ExplorerFromArgs(args map[string]string) *Explorer {
	explorer := &Explorer{}

	for key, value := range args {
		switch key {
		case "mintscan-prefix":
			return MintscanExplorerFromArgs(args)
		case "ping-prefix":
			return PingExplorerFromArgs(args)
		case "name":
			explorer.Name = value
		case "chain":
			explorer.Chain = value
		case "proposal-link-pattern":
			explorer.ProposalLinkPattern = value
		case "wallet-link-pattern":
			explorer.WalletLinkPattern = value
		case "validator-link-pattern":
			explorer.ValidatorLinkPattern = value
		case "main-link":
			explorer.MainLink = value
		}
	}

	return explorer
}

func MintscanExplorerFromArgs(args map[string]string) *Explorer {
	prefix := args["mintscan-prefix"]

	name := args["name"]
	if name == "" {
		name = "Mintscan"
	}

	return &Explorer{
		Chain:                args["chain"],
		Name:                 name,
		ProposalLinkPattern:  fmt.Sprintf("https://mintscan.io/%s/proposals/%%s", prefix),
		WalletLinkPattern:    fmt.Sprintf("https://mintscan.io/%s/account/%%s", prefix),
		ValidatorLinkPattern: fmt.Sprintf("https://mintscan.io/%s/validators/%%s", prefix),
		MainLink:             fmt.Sprintf("https://mintscan.io/%s", prefix),
	}
}

func PingExplorerFromArgs(args map[string]string) *Explorer {
	prefix := args["ping-prefix"]

	name := args["name"]
	if name == "" {
		name = "Ping"
	}

	host := args["ping-host"]
	if host == "" {
		host = "https://ping.pub"
	}

	return &Explorer{
		Chain:                args["chain"],
		Name:                 name,
		ProposalLinkPattern:  fmt.Sprintf("%s/%s/gov/%%s", host, prefix),
		WalletLinkPattern:    fmt.Sprintf("%s/%s/account/%%s", host, prefix),
		ValidatorLinkPattern: fmt.Sprintf("%s/%s/staking/%%s", host, prefix),
		MainLink:             fmt.Sprintf("%s/%s", host, prefix),
	}
}

func (e *Explorer) Validate() error {
	if e.Chain == "" {
		return errors.New("chain name cannot be empty")
	}

	if e.Name == "" {
		return errors.New("name cannot be empty")
	}

	if e.ProposalLinkPattern == "" {
		return errors.New("proposal link pattern cannot be empty")
	}

	if e.WalletLinkPattern == "" {
		return errors.New("wallet link pattern cannot be empty")
	}

	if e.ValidatorLinkPattern == "" {
		return errors.New("validator link pattern cannot be empty")
	}

	if e.MainLink == "" {
		return errors.New("main link cannot be empty")
	}

	return nil
}

func (e *Explorer) DisplayWarnings(chainName string) []Warning {
	warnings := make([]Warning, 0)

	if e.WalletLinkPattern == "" {
		warnings = append(warnings, Warning{
			Labels:  map[string]string{"chain": chainName},
			Message: "wallet-link-pattern for explorer is not set, cannot generate wallet links",
		})
	}

	if e.ProposalLinkPattern == "" {
		warnings = append(warnings, Warning{
			Labels:  map[string]string{"chain": chainName},
			Message: "proposal-link-pattern for explorer is not set, cannot generate proposal links",
		})
	}

	return warnings
}

type Explorers []*Explorer

func (e Explorers) GetExplorersByChain(chain string) Explorers {
	return utils.Filter(e, func(e *Explorer) bool {
		return e.Chain == chain
	})
}

func (e Explorers) GetValidatorLinks(valoper string) []Link {
	links := make([]Link, len(e))
	for index, explorer := range e {
		links[index] = Link{
			Text: explorer.Name,
			Href: fmt.Sprintf(explorer.ValidatorLinkPattern, valoper),
		}
	}

	return links
}

func (e Explorers) GetProposalLinks(proposalID string) []Link {
	links := make([]Link, len(e))
	for index, explorer := range e {
		links[index] = Link{
			Text: explorer.Name,
			Href: fmt.Sprintf(explorer.ProposalLinkPattern, proposalID),
		}
	}

	return links
}

func (e Explorers) GetChainLinks(chainName string) []Link {
	links := make([]Link, 0)
	for _, explorer := range e {
		if explorer.Chain != chainName {
			continue
		}

		links = append(links, Link{
			Text: explorer.Name,
			Href: explorer.MainLink,
		})
	}

	return links
}

func (e Explorers) GetWalletLinks(wallet *WalletLink) []Link {
	links := make([]Link, 0)
	for _, explorer := range e {
		links = append(links, Link{
			Text: explorer.Name,
			Href: fmt.Sprintf(explorer.WalletLinkPattern, wallet.Address),
		})
	}

	return links
}
