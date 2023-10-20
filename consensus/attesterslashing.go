package consensus

import (
	"github.com/attestantio/go-eth2-client/spec/phase0"
	fuzz "github.com/google/gofuzz"
)

type AttestorSlashingJSON struct {
	Attestation1 *phase0.IndexedAttestation `json:"attestation_1"`
	Attestation2 *phase0.IndexedAttestation `json:"attestation_2"`
}

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
