package consensus

import (
	"github.com/attestantio/go-eth2-client/spec/phase0"
	fuzz "github.com/google/gofuzz"
)

func RandomAttesterSlashing() *phase0.AttesterSlashing {
	var attesterSlashing phase0.AttesterSlashing
	f := fuzz.New().NilChance(0)
	for true {
		f.Fuzz(&attesterSlashing)
		_, err := attesterSlashing.MarshalSSZ()
		if err == nil {
			break
		}
	}
	return &attesterSlashing
}
