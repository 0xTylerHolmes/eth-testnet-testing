package eth_testnet_tool

import (
	"eth-testnet-tool/consensus"
	v1 "github.com/attestantio/go-eth2-client/api/v1"
	"github.com/attestantio/go-eth2-client/spec"
	"github.com/attestantio/go-eth2-client/spec/capella"
	"github.com/attestantio/go-eth2-client/spec/phase0"
)

// TODO: use internal API until PRs finished to unify

// GetBeaconBlockHeader get the BeaconBlockHeader from the client api endpoint at the provided slot.
func (c *ConsensusClient) GetBeaconBlockHeader(slot string) (*v1.BeaconBlockHeader, error) {
	return consensus.GetBeaconBlockHeader(c.APIEndpoint, slot)
}

func (c *ConsensusClient) GetCheckpoints(slot string) (*v1.Finality, error) {
	return consensus.GetCheckpoints(c.APIEndpoint, slot)
}

func (c *ConsensusClient) PostAttestation(attestation phase0.Attestation) error {
	return consensus.PostAttestation(c.APIEndpoint, attestation)
}

func (c *ConsensusClient) PostAttesterSlashing(slashing phase0.AttesterSlashing) error {
	return consensus.PostAttestorSlashing(c.APIEndpoint, slashing)
}

func (c *ConsensusClient) PostProposerSlashing(slashing phase0.ProposerSlashing) error {
	return consensus.PostProposerSlashing(c.APIEndpoint, slashing)
}

func (c *ConsensusClient) PostSignedVoluntaryExit(exit phase0.SignedVoluntaryExit) error {
	return consensus.PostSignedVoluntaryExit(c.APIEndpoint, exit)
}

func (c *ConsensusClient) PostSignedBLSToExecutionChange(change capella.SignedBLSToExecutionChange) error {
	return consensus.PostSignedBLSToExecutionChange(c.APIEndpoint, change)
}

func (c *ConsensusClient) PostSignedBeaconBlock(block spec.VersionedSignedBeaconBlock) error {
	return consensus.PostSignedBeaconBlock(c.APIEndpoint, block)
}

// Group Functions

func (t *TestnetClients) GetAllBeaconBlockHeaders(slot string) map[string]*v1.BeaconBlockHeader {
	beaconBlockHeaders := make(map[string]*v1.BeaconBlockHeader)
	for clientName, beaconAPI := range t.ConsensusClients {
		bbh, err := beaconAPI.GetBeaconBlockHeader(slot)
		if err != nil {
			beaconBlockHeaders[clientName] = nil
		} else {
			beaconBlockHeaders[clientName] = bbh
		}

	}
	return beaconBlockHeaders
}

func (t *TestnetClients) GetAllFinalityCheckpoints(slot string) map[string]*v1.Finality {
	finalityCheckpoint := make(map[string]*v1.Finality)
	for clientName, beaconAPI := range t.ConsensusClients {
		bbh, err := beaconAPI.GetCheckpoints(slot)
		if err != nil {
			finalityCheckpoint[clientName] = nil
		} else {
			finalityCheckpoint[clientName] = bbh
		}

	}
	return finalityCheckpoint
}
