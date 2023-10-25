package consensus_objects

import (
	"github.com/attestantio/go-eth2-client/spec/capella"
	fuzz "github.com/google/gofuzz"
)

type SignedBLSToExecutionChangeJSON struct {
	Message   *capella.BLSToExecutionChange `json:"message"`
	Signature string                        `json:"signature"`
}

// RandomBLSToExecutionChange creates a random ssz-able BLSToExecutionChange
func RandomBLSToExecutionChange() *capella.BLSToExecutionChange {
	var blsToExecutionChange capella.BLSToExecutionChange
	f := fuzz.New().NilChance(0)
	for {
		f.Fuzz(&blsToExecutionChange)
		_, err := blsToExecutionChange.MarshalSSZ()
		if err == nil {
			break
		}
	}
	return &blsToExecutionChange
}

func RandomSignedBLSToExecutionChange() *capella.SignedBLSToExecutionChange {
	var signedBlsToExecutionChange capella.SignedBLSToExecutionChange
	f := fuzz.New().NilChance(0)
	for {
		f.Fuzz(&signedBlsToExecutionChange)
		_, err := signedBlsToExecutionChange.MarshalSSZ()
		if err == nil {
			break
		}
	}
	return &signedBlsToExecutionChange
}
