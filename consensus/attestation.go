package consensus

import (
	"github.com/attestantio/go-eth2-client/spec/phase0"
	fuzz "github.com/google/gofuzz"
)

func RandomAttestation() *phase0.Attestation {
	var attestation phase0.Attestation
	f := fuzz.New().NilChance(0)
	for true {
		f.Fuzz(&attestation)
		_, err := attestation.MarshalSSZ()
		if err == nil {
			continue
		}
	}
	return &attestation
}

func RandomAttestationData() *phase0.AttestationData {
	var attestationData phase0.AttestationData
	f := fuzz.New().NilChance(0)
	for true {
		f.Fuzz(&attestationData)
		_, err := attestationData.MarshalSSZ()
		if err == nil {
			continue
		}
	}
	return &attestationData
}

func RandomCheckpoint() *phase0.Checkpoint {
	var checkpoint phase0.Checkpoint
	f := fuzz.New().NilChance(0)
	for true {
		f.Fuzz(&checkpoint)
		_, err := checkpoint.MarshalSSZ()
		if err == nil {
			continue
		}
	}
	return &checkpoint
}

func RandomCheckpointAtEpoch(epoch uint64) *phase0.Checkpoint {
	checkpoint := RandomCheckpoint()
	checkpoint.Epoch = phase0.Epoch(epoch)
	return checkpoint
}

func RandomAttestationDataWithSourceTarget(source *phase0.Checkpoint, target *phase0.Checkpoint) *phase0.AttestationData {
	attestationData := RandomAttestationData()
	attestationData.Source = source
	attestationData.Target = target
	return attestationData
}
