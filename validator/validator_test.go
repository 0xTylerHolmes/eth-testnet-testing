package validator

import (
	"encoding/hex"
	"fmt"
	"github.com/stretchr/testify/require"
	"testing"
)

var (
	ValidatorMnemonic     = "ocean style run case glory clip into nature guess jacket document firm fiscal hello kite disagree symptom tide net coral envelope wink render festival"
	GenesisValidatorCount = 20
)

func TestValidator_GenerateValidators(t *testing.T) {
	validators, err := GetValidatorsFromMnemonic(ValidatorMnemonic, 0, uint64(GenesisValidatorCount))
	require.NoError(t, err)
	for _, v := range validators {
		require.Equal(t, v.ValidatorPublicKey.String(), fmt.Sprintf("0x%s", hex.EncodeToString(v.ValidatorKey.PublicKey().Marshal()[:])))
	}

}
