package converter

import (
	upgradeTypes "cosmossdk.io/x/upgrade/types"
	"github.com/cosmos/cosmos-sdk/codec"
	codecTypes "github.com/cosmos/cosmos-sdk/codec/types"
	"github.com/cosmos/cosmos-sdk/std"
	govV1Types "github.com/cosmos/cosmos-sdk/x/gov/types/v1"
	govV1beta1Types "github.com/cosmos/cosmos-sdk/x/gov/types/v1beta1"
	paramsProposalTypes "github.com/cosmos/cosmos-sdk/x/params/types/proposal"
	"github.com/cosmos/gogoproto/proto"
)

type Converter struct {
	registry   codecTypes.InterfaceRegistry
	parseCodec *codec.ProtoCodec
}

func NewConverter() *Converter {
	interfaceRegistry := codecTypes.NewInterfaceRegistry()
	std.RegisterInterfaces(interfaceRegistry)
	govV1Types.RegisterInterfaces(interfaceRegistry)
	govV1beta1Types.RegisterInterfaces(interfaceRegistry)
	paramsProposalTypes.RegisterInterfaces(interfaceRegistry)
	upgradeTypes.RegisterInterfaces(interfaceRegistry)

	parseCodec := codec.NewProtoCodec(interfaceRegistry)

	return &Converter{
		registry:   interfaceRegistry,
		parseCodec: parseCodec,
	}
}

func (c *Converter) UnmarshalJSON(bytes []byte, target proto.Message) error {
	return c.parseCodec.UnmarshalJSON(bytes, target)
}

func (c *Converter) UnpackProposal(proposal govV1beta1Types.Proposal) error {
	return proposal.UnpackInterfaces(c.parseCodec)
}
