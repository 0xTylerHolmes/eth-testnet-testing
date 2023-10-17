package consensus

import (
	"github.com/attestantio/go-eth2-client/spec/phase0"
	fuzz "github.com/google/gofuzz"
)

func RandomVoluntaryExit() *phase0.VoluntaryExit {
	var voluntaryExit phase0.VoluntaryExit
	f := fuzz.New().NilChance(0)
	for true {
		f.Fuzz(&voluntaryExit)
		_, err := voluntaryExit.MarshalSSZ()
		if err == nil {
			break
		}
	}
	return &voluntaryExit
}
