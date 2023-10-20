package consensus

import (
	"encoding/json"
	"fmt"
	"github.com/attestantio/go-eth2-client/spec"
	"github.com/attestantio/go-eth2-client/spec/capella"
	"github.com/attestantio/go-eth2-client/spec/phase0"
	"github.com/stretchr/testify/require"
	"testing"
)

const ENDPOINT = "http://10.0.20.4:5052"

// / broken due to formatting for genesis (expected) github.com/attestantio/go-eth2-client@v0.18.3/api/v1/genesis.go:61
// / will not fix unless necessary.
func TestGetGenesis(t *testing.T) {
	var genesis GenesisResponseJSON
	resp, err := getRequest(ENDPOINT, []byte("/eth/v1/beacon/genesis"))
	require.NoError(t, err)
	fmt.Println(string(resp))
	err = json.Unmarshal(resp, &genesis)
	require.NoError(t, err)
	fmt.Println(genesis.Data.String())
	fmt.Println(string(resp))
}

func TestGetBeaconBlockHeader(t *testing.T) {
	bbh, err := GetBeaconBlockHeader(ENDPOINT, "head")
	require.NoError(t, err)
	fmt.Println(bbh.String())
}

func TestGetCheckpoint(t *testing.T) {
	checkpoint, err := GetCheckpoints(ENDPOINT, "head")
	require.NoError(t, err)
	fmt.Println(checkpoint.String())
}

func TestPostAttestation(t *testing.T) {
	var attestation phase0.Attestation
	attestation = *RandomAttestation()
	err := PostAttestation(ENDPOINT, attestation)
	// we expect an error for a random attestation
	fmt.Println(err)
}

func TestPostAttestorSlashing(t *testing.T) {
	var attesteterSlashing phase0.AttesterSlashing
	attesteterSlashing = *RandomAttesterSlashing()
	err := PostAttestorSlashing(ENDPOINT, attesteterSlashing)
	// we expect errors from random slashings
	fmt.Println(err)
}

func TestPostProposerSlashing(t *testing.T) {
	var proposerSlashing phase0.ProposerSlashing
	proposerSlashing = *RandomProposerSlashing()
	err := PostProposerSlashing(ENDPOINT, proposerSlashing)
	// we expect errors from random slashings
	fmt.Println(err)
}

func TestPostVoluntaryExit(t *testing.T) {
	var signedVoluntaryExit phase0.SignedVoluntaryExit
	signedVoluntaryExit = *RandomSignedVoluntaryExit()
	err := PostSignedVoluntaryExit(ENDPOINT, signedVoluntaryExit)
	// we expect errors from random exits
	fmt.Println(err)
}

func TestPostSignedBLSToExecutionChange(t *testing.T) {
	var signedBLSToExecutionChange capella.SignedBLSToExecutionChange
	signedBLSToExecutionChange = *RandomSignedBLSToExecutionChange()
	err := PostSignedBLSToExecutionChange(ENDPOINT, signedBLSToExecutionChange)
	fmt.Println(err)
}

// TODO failing
func TestPostSignedBeaconBlock(t *testing.T) {
	var signedBeaconBlock spec.VersionedSignedBeaconBlock
	signedBeaconBlock = *RandomDenebSignedBeaconBlock()
	fmt.Printf("Generated random deneb signed beacon blocK: %s\n", signedBeaconBlock.String())
	err := PostSignedBeaconBlock(ENDPOINT, signedBeaconBlock)
	fmt.Println(err)
}
