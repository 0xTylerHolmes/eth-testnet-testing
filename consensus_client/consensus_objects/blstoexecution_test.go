package consensus_objects

import (
	"github.com/attestantio/go-eth2-client/spec/capella"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestRoundTripExecutionChange(t *testing.T) {
	rtChange := capella.BLSToExecutionChange{}
	for i := 0; i < 10; i++ {
		randExecutionChange := RandomBLSToExecutionChange()
		changeSSZ, err := randExecutionChange.MarshalSSZ()
		require.NoError(t, err)
		err = rtChange.UnmarshalSSZ(changeSSZ)
		require.NoError(t, err)
		require.Equal(t, *randExecutionChange, rtChange)
	}
}

func TestRoundTripSignedExecutionChange(t *testing.T) {
	rtSignedChange := capella.SignedBLSToExecutionChange{}
	for i := 0; i < 10; i++ {
		randSignedExecutionChange := RandomSignedBLSToExecutionChange()
		signedChangeSSZ, err := randSignedExecutionChange.MarshalSSZ()
		require.NoError(t, err)
		err = rtSignedChange.UnmarshalSSZ(signedChangeSSZ)
		require.NoError(t, err)
		require.Equal(t, *randSignedExecutionChange, rtSignedChange)
	}

}
