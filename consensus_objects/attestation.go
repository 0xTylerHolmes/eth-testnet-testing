package consensus_objects

import (
	v1 "github.com/attestantio/go-eth2-client/api/v1"
	"github.com/attestantio/go-eth2-client/spec/phase0"
	fuzz "github.com/google/gofuzz"
)

type CheckpointResponseJSON struct {
	Data                *v1.Finality `json:"data"`
	ExecutionOptimistic bool         `json:"execution_optimistic"`
	Finalized           bool         `json:"finalized"`
}

func RandomAttestation() *phase0.Attestation {
	var attestation phase0.Attestation
	f := fuzz.New().NilChance(0)
	for true {
		f.Fuzz(&attestation)
		_, err := attestation.MarshalSSZ()
		if err == nil {
			break
		}
	}
	return &attestation
}
