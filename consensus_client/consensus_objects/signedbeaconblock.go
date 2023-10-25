package consensus_objects

import (
	"fmt"
	v1 "github.com/attestantio/go-eth2-client/api/v1"
	"github.com/attestantio/go-eth2-client/spec"
	"github.com/attestantio/go-eth2-client/spec/altair"
	"github.com/attestantio/go-eth2-client/spec/bellatrix"
	"github.com/attestantio/go-eth2-client/spec/capella"
	"github.com/attestantio/go-eth2-client/spec/deneb"
	"github.com/attestantio/go-eth2-client/spec/phase0"
	fuzz "github.com/google/gofuzz"
)

// BeaconBlockHeaderResponseJSON represents the response in JSON for a /eth/v1/block/header request
type BeaconBlockHeaderResponseJSON struct {
	Data                *v1.BeaconBlockHeader `json:"data"`
	ExecutionOptimistic bool                  `json:"execution_optimistic"`
	Finalized           bool                  `json:"finalized"`
}

type Phase0SignedBeaconBlockProposalJSON struct {
	Data *phase0.SignedBeaconBlock `json:"data"`
}

type AltairSignedBeaconBlockProposalJSON struct {
	Data *altair.SignedBeaconBlock `json:"data"`
}

type BellatrixSignedBeaconBlockProposalJSON struct {
	Data *bellatrix.SignedBeaconBlock `json:"data"`
}

type CapellaSignedBeaconBlockProposalJSON struct {
	Data *capella.SignedBeaconBlock `json:"data"`
}

type DenebSignedBeaconBlockProposalJSON struct {
	Data *deneb.SignedBeaconBlock `json:"data"`
}

// TODO round trip ssz fails
func RandomPhase0SignedBeaconBlock() *spec.VersionedSignedBeaconBlock {
	var signedBeaconBlock phase0.SignedBeaconBlock
	f := fuzz.New().NilChance(0)
	for {
		f.Fuzz(&signedBeaconBlock)
		_, err := signedBeaconBlock.MarshalSSZ()
		if err == nil {
			break
		}
	}
	return &spec.VersionedSignedBeaconBlock{
		Version: spec.DataVersionPhase0,
		Phase0:  &signedBeaconBlock,
	}
}

// TODO round trip ssz fails
func RandomAltairSingedBeaconBlock() *spec.VersionedSignedBeaconBlock {
	var signedBeaconBlock altair.SignedBeaconBlock
	f := fuzz.New().NilChance(0)
	for {
		f.Fuzz(&signedBeaconBlock)
		_, err := signedBeaconBlock.MarshalSSZ()
		if err == nil {
			break
		}
	}
	return &spec.VersionedSignedBeaconBlock{
		Version: spec.DataVersionAltair,
		Altair:  &signedBeaconBlock,
	}
}

// TODO round trip ssz fails
func RandomBellatrixSignedBeaconBlock() *spec.VersionedSignedBeaconBlock {
	var signedBeaconBlock bellatrix.SignedBeaconBlock
	f := fuzz.New().NilChance(0)
	for {
		f.Fuzz(&signedBeaconBlock)
		_, err := signedBeaconBlock.MarshalSSZ()
		if err == nil {
			break
		}
	}
	return &spec.VersionedSignedBeaconBlock{
		Version:   spec.DataVersionBellatrix,
		Bellatrix: &signedBeaconBlock,
	}
}

// TODO round trip ssz fails
func RandomCapellaSignedBeaconBlock() *spec.VersionedSignedBeaconBlock {
	var signedBeaconBlock capella.SignedBeaconBlock
	f := fuzz.New().NilChance(0)
	for {
		f.Fuzz(&signedBeaconBlock)
		_, err := signedBeaconBlock.MarshalJSON()
		if err == nil {
			break
		}
	}
	return &spec.VersionedSignedBeaconBlock{
		Version: spec.DataVersionCapella,
		Capella: &signedBeaconBlock,
	}
}

// TODO round trip ssz fails
func RandomDenebSignedBeaconBlock() *spec.VersionedSignedBeaconBlock {
	var signedBeaconBlock deneb.SignedBeaconBlock
	f := fuzz.New().NilChance(0)
	for {
		f.Fuzz(&signedBeaconBlock)
		_, err := signedBeaconBlock.MarshalSSZ()
		if err == nil {
			break
		} else {
			fmt.Println("generated block wasn't serializable, retrying.")
		}
	}
	return &spec.VersionedSignedBeaconBlock{
		Version: spec.DataVersionDeneb,
		Deneb:   &signedBeaconBlock,
	}
}
