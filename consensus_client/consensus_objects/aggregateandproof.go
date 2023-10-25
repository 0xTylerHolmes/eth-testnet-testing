package consensus_objects

import (
	"github.com/attestantio/go-eth2-client/spec/phase0"
	fuzz "github.com/google/gofuzz"
)

func RandomAggregateAndProof() *phase0.SignedAggregateAndProof {
	var signedAggregateAndProof phase0.SignedAggregateAndProof
	f := fuzz.New().NilChance(0)
	for {
		f.Fuzz(&signedAggregateAndProof)
		_, err := signedAggregateAndProof.MarshalSSZ()
		if err == nil {
			break
		}
	}
	return &signedAggregateAndProof
}
