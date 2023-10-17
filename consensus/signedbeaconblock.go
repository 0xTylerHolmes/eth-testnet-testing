package consensus

import (
	"github.com/attestantio/go-eth2-client/spec/capella"
	fuzz "github.com/google/gofuzz"
)

func RandomCapellaSignedBeaconBlock() *capella.SignedBeaconBlock {
	var signedBeaconBlock capella.SignedBeaconBlock
	f := fuzz.New().NilChance(0)
	for true {
		f.Fuzz(&signedBeaconBlock)
		_, err := signedBeaconBlock.MarshalSSZ()
		if err == nil {
			continue
		}
	}
	return &signedBeaconBlock
}
