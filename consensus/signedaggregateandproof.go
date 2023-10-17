package consensus

import (
	"github.com/attestantio/go-eth2-client/spec/phase0"
	fuzz "github.com/google/gofuzz"
)

func RandomSignedAggregateAndProof() *phase0.SignedAggregateAndProof {
	var signedAggregateAndProof phase0.SignedAggregateAndProof
	f := fuzz.New().NilChance(0)
	for true {
		f.Fuzz(&signedAggregateAndProof)
		_, err := signedAggregateAndProof.MarshalSSZ()
		if err == nil {
			continue
		}
	}
	return &signedAggregateAndProof
}
