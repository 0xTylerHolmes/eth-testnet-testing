package consensus

import (
	"github.com/attestantio/go-eth2-client/spec/phase0"
	fuzz "github.com/google/gofuzz"
)

func RandomProposerSlashing() *phase0.ProposerSlashing {
	var proposerSlashing phase0.ProposerSlashing
	f := fuzz.New().NilChance(0)
	for true {
		f.Fuzz(&proposerSlashing)
		_, err := proposerSlashing.MarshalSSZ()
		if err == nil {
			break
		}
	}
	return &proposerSlashing
}
