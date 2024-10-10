package converter

import (
	upgradeTypes "cosmossdk.io/x/upgrade/types"
	"github.com/cosmos/cosmos-sdk/codec"
	codecTypes "github.com/cosmos/cosmos-sdk/codec/types"
	"github.com/cosmos/cosmos-sdk/std"
	sdkTypes "github.com/cosmos/cosmos-sdk/types"
	govV1Types "github.com/cosmos/cosmos-sdk/x/gov/types/v1"
	govV1beta1Types "github.com/cosmos/cosmos-sdk/x/gov/types/v1beta1"
	paramsProposalTypes "github.com/cosmos/cosmos-sdk/x/params/types/proposal"
	stakingTypes "github.com/cosmos/cosmos-sdk/x/staking/types"
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

func (c *Converter) Unmarshal(bytes []byte, target proto.Message) error {
	return c.parseCodec.UnmarshalJSON(bytes, target)
}

func (c *Converter) UnpackProposal(proposal govV1beta1Types.Proposal) error {
	return proposal.UnpackInterfaces(c.parseCodec)
}

func (c *Converter) GetValidatorConsAddr(validator stakingTypes.Validator) string {
	if err := validator.UnpackInterfaces(c.parseCodec); err != nil {
		panic(err)
	}

	addr, err := validator.GetConsAddr()
	if err != nil {
		panic(err)
	}

	return sdkTypes.ConsAddress(addr).String()
}
