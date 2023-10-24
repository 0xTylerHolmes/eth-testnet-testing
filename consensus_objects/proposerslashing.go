package consensus_objects

import (
	"github.com/attestantio/go-eth2-client/spec/phase0"
	fuzz "github.com/google/gofuzz"
)

type ProposerSlashingJSON struct {
	SignedHeader1 *phase0.SignedBeaconBlockHeader `json:"signed_header_1"`
	SignedHeader2 *phase0.SignedBeaconBlockHeader `json:"signed_header_2"`
}

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
