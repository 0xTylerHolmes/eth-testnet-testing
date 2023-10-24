package consensus_objects

import (
	"fmt"
	"github.com/attestantio/go-eth2-client/spec/phase0"
	"github.com/stretchr/testify/require"
	"testing"
)

//func TestRandomCheckpoint(t *testing.T) {
//	checkpointMap := make(map[phase0.Checkpoint]int)
//	for i := 0; i < 10; i++ {
//		checkPoint := RandomCheckpoint()
//		if _, ok := checkpointMap[*checkPoint]; ok {
//			checkpointMap[*checkPoint] += 1
//		} else {
//			checkpointMap[*checkPoint] = 1
//		}
//	}
//	for _, numUnique := range checkpointMap {
//		require.Equal(t, 1, numUnique)
//	}
//}

func TestRandomAttestation(t *testing.T) {
	var pastAttestations []phase0.Attestation

	for i := 0; i < 10; i++ {
		attestation := RandomAttestation()
		_, err := attestation.MarshalSSZ()
		require.NoError(t, err)
		pastAttestations = append(pastAttestations, *attestation)
	}
	fmt.Println(pastAttestations)
	//TODO uniqueness test that isn't visual
}

//func TestRandomCheckpointAtEpoch(t *testing.T) {
//	targetEpoch := uint64(777)
//	checkpointMap := make(map[phase0.Checkpoint]int)
//	for i := 0; i < 10; i++ {
//		checkPoint := RandomCheckpointAtEpoch(targetEpoch)
//		if _, ok := checkpointMap[*checkPoint]; ok {
//			checkpointMap[*checkPoint] += 1
//		} else {
//			checkpointMap[*checkPoint] = 1
//		}
//	}
//	for checkpoint, _ := range checkpointMap {
//		require.Equal(t, targetEpoch, uint64(checkpoint.Epoch))
//	}
//}
