package consensus

import (
	"github.com/attestantio/go-eth2-client/spec/capella"
	fuzz "github.com/google/gofuzz"
)

// RandomBLSToExecutionChange creates a random ssz-able BLSToExecutionChange
func RandomBLSToExecutionChange() *capella.BLSToExecutionChange {
	var blsToExecutionChange capella.BLSToExecutionChange
	f := fuzz.New().NilChance(0)
	for true {
		f.Fuzz(&blsToExecutionChange)
		_, err := blsToExecutionChange.MarshalSSZ()
		if err == nil {
			break
		}
	}
	return &blsToExecutionChange
}
