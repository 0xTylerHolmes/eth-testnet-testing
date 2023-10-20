package consensus

import (
	v1 "github.com/attestantio/go-eth2-client/api/v1"
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

type Phase0BeaconBlockProposalJSON struct {
	Data *phase0.BeaconBlock `json:"data"`
}

type AltairBeaconBlockProposalJSON struct {
	Data *altair.BeaconBlock `json:"data"`
}

type BellatrixBeaconBlockProposalJSON struct {
	Data *bellatrix.BeaconBlock `json:"data"`
}

type CapellaBeaconBlockProposalJSON struct {
	Data *capella.BeaconBlock `json:"data"`
}

type DenebBeaconBlockProposalJSON struct {
	Data *deneb.BeaconBlock `json:"data"`
}

func RandomPhase0SignedBeaconBlock() *phase0.SignedBeaconBlock {
	var signedBeaconBlock phase0.SignedBeaconBlock
	f := fuzz.New().NilChance(0)
	for true {
		f.Fuzz(&signedBeaconBlock)
		_, err := signedBeaconBlock.MarshalSSZ()
		if err == nil {
			break
		}
	}
	return &signedBeaconBlock
}
func RandomAltairSingedBeaconBlock() *altair.SignedBeaconBlock {
	var signedBeaconBlock altair.SignedBeaconBlock
	f := fuzz.New().NilChance(0)
	for true {
		f.Fuzz(&signedBeaconBlock)
		_, err := signedBeaconBlock.MarshalSSZ()
		if err == nil {
			break
		}
	}
	return &signedBeaconBlock
}

func RandomBellatrixSignedBeaconBlock() *bellatrix.SignedBeaconBlock {
	var signedBeaconBlock bellatrix.SignedBeaconBlock
	f := fuzz.New().NilChance(0)
	for true {
		f.Fuzz(&signedBeaconBlock)
		_, err := signedBeaconBlock.MarshalSSZ()
		if err == nil {
			break
		}
	}
	return &signedBeaconBlock
}

func RandomCapellaSignedBeaconBlock() *capella.SignedBeaconBlock {
	var signedBeaconBlock capella.SignedBeaconBlock
	f := fuzz.New().NilChance(0)
	for true {
		f.Fuzz(&signedBeaconBlock)
		_, err := signedBeaconBlock.MarshalSSZ()
		if err == nil {
			break
		}
	}
	return &signedBeaconBlock
}

func RandomDenebSignedBeaconBlock() *deneb.SignedBeaconBlock {
	var signedBeaconBlock deneb.SignedBeaconBlock
	f := fuzz.New().NilChance(0)
	for true {
		f.Fuzz(&signedBeaconBlock)
		_, err := signedBeaconBlock.MarshalSSZ()
		if err == nil {
			break
		}
	}
	return &signedBeaconBlock
}
