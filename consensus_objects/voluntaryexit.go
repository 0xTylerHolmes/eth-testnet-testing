package consensus_objects

import (
	"github.com/attestantio/go-eth2-client/spec/phase0"
	fuzz "github.com/google/gofuzz"
)

type SignedVoluntaryExitJSON struct {
	Message   *phase0.VoluntaryExit `json:"message"`
	Signature string                `json:"signature"`
}

type VoluntaryExitJSON struct {
	Epoch          string `json:"epoch"`
	ValidatorIndex string `json:"validator_index"`
}

// RandomVoluntaryExit create a random voluntary exit
func RandomVoluntaryExit() *phase0.VoluntaryExit {
	var voluntaryExit phase0.VoluntaryExit
	f := fuzz.New().NilChance(0)
	for {
		f.Fuzz(&voluntaryExit)
		_, err := voluntaryExit.MarshalSSZ()
		if err == nil {
			break
		}
	}
	return &voluntaryExit
}

// RandomSignedVoluntaryExit create a random signed voluntary exit
func RandomSignedVoluntaryExit() *phase0.SignedVoluntaryExit {
	var signedVoluntaryExit phase0.SignedVoluntaryExit
	f := fuzz.New().NilChance(0)
	for {
		f.Fuzz(&signedVoluntaryExit)
		_, err := signedVoluntaryExit.MarshalSSZ()
		if err == nil {
			break
		}
	}
	return &signedVoluntaryExit
}
