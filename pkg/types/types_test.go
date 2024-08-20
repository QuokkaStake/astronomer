package types

import (
	"main/pkg/constants"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestValidatorInfoIsActive(t *testing.T) {
	t.Parallel()

	require.False(t, ValidatorInfo{Status: ""}.Active())
	require.True(t, ValidatorInfo{Status: constants.ValidatorStatusBonded}.Active())
}

func TestValidatorInfoFormatCommission(t *testing.T) {
	t.Parallel()

	require.Equal(t, "5.00", ValidatorInfo{Commission: 0.05}.FormatCommission())
}

func TestValidatorInfoGetVotingPowerPercent(t *testing.T) {
	t.Parallel()

	require.Equal(t, "5.00", ValidatorInfo{VotingPowerPercent: 0.05}.GetVotingPowerPercent())
}
