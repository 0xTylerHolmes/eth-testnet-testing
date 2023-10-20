package eth_testnet_tool

import (
	"eth-testnet-tool/consensus"
	"fmt"
	"github.com/stretchr/testify/require"
	"testing"
)

const ExampleTestnetConfigPath = "./example/configs/example-testnet-clients-config.json"

func getTestnetClientsFromExampleConfig() (*TestnetClients, error) {
	testnetClientsJSON, err := GetTestnetClientsFromConfig(ExampleTestnetConfigPath)
	if err != nil {
		return nil, err
	}
	testnetClients := testnetClientsJSON.GetTestnetClients()
	return testnetClients, nil
}

func TestConsensusClient_GetBeaconBlockHeader(t *testing.T) {
	testnetClients, err := getTestnetClientsFromExampleConfig()
	require.NoError(t, err)
	beaconBlockHeaders := testnetClients.GetAllBeaconBlockHeaders("head")
	for clientName, header := range beaconBlockHeaders {
		fmt.Printf("%s : %s\n", clientName, header.String())
	}
}

func TestConsensusClient_GetFinalityCheckpoint(t *testing.T) {
	testnetClients, err := getTestnetClientsFromExampleConfig()
	require.NoError(t, err)
	finalityCheckpoints := testnetClients.GetAllFinalityCheckpoints("head")
	for clientName, checkpoint := range finalityCheckpoints {
		fmt.Printf("%s : %s\n", clientName, checkpoint.String())
	}
}

func TestPostRandomAttestationToRandomClient(t *testing.T) {
	testnetClients, err := getTestnetClientsFromExampleConfig()
	require.NoError(t, err)
	consensusClient := testnetClients.GetRandomConsensusClient()
	randomAttestation := consensus.RandomAttestation()
	fmt.Printf("generated random attestation: \n%s\n", randomAttestation.String())
	err = consensusClient.PostAttestation(*randomAttestation)
	if err != nil {
		fmt.Printf("client at: %s responded with error: %s", consensusClient.APIEndpoint, err.Error())
	}
}

func TestPostAttestorSlashingToRandomClient(t *testing.T) {
	testnetClients, err := getTestnetClientsFromExampleConfig()
	require.NoError(t, err)
	consensusClient := testnetClients.GetRandomConsensusClient()
	randomAttestorSlashing := consensus.RandomAttesterSlashing()
	fmt.Printf("generated random attester slashing: %s", randomAttestorSlashing.String())
	err = consensusClient.PostAttesterSlashing(*randomAttestorSlashing)
	if err != nil {
		fmt.Printf("client at: %s responded with error: %s", consensusClient.APIEndpoint, err.Error())
	}
}

func TestPostProposerSlashingToRandomClient(t *testing.T) {
	testnetClients, err := getTestnetClientsFromExampleConfig()
	require.NoError(t, err)
	consensusClient := testnetClients.GetRandomConsensusClient()
	randomProposerSlashing := consensus.RandomProposerSlashing()
	fmt.Printf("generated random proposer slashing: %s", randomProposerSlashing.String())
	err = consensusClient.PostProposerSlashing(*randomProposerSlashing)
	if err != nil {
		fmt.Printf("client at: %s responded with error: %s", consensusClient.APIEndpoint, err.Error())
	}
}

func TestPostVoluntaryExitToRandomClient(t *testing.T) {
	testnetClients, err := getTestnetClientsFromExampleConfig()
	require.NoError(t, err)
	consensusClient := testnetClients.GetRandomConsensusClient()
	randomSignedVoluntaryExit := consensus.RandomSignedVoluntaryExit()
	fmt.Printf("generated random voluntary exit: %s", randomSignedVoluntaryExit.String())
	err = consensusClient.PostSignedVoluntaryExit(*randomSignedVoluntaryExit)
	if err != nil {
		fmt.Printf("client at: %s responded with error: %s", consensusClient.APIEndpoint, err.Error())
	}
}

func TestPostSignedBLSToExecutionChangeToRandomClient(t *testing.T) {
	testnetClients, err := getTestnetClientsFromExampleConfig()
	require.NoError(t, err)
	consensusClient := testnetClients.GetRandomConsensusClient()
	randomSignedBLSToExecutionChange := consensus.RandomSignedBLSToExecutionChange()
	fmt.Printf("generated random bls to execution change: %s", randomSignedBLSToExecutionChange.String())
	err = consensusClient.PostSignedBLSToExecutionChange(*randomSignedBLSToExecutionChange)
	if err != nil {
		fmt.Printf("client at: %s responded with error: %s", consensusClient.APIEndpoint, err.Error())
	}
}

// TODO: fix this blocks currently are not working with submission
//func TestPostSignedBeaconBlockToRandomClient(t *testing.T) {
//	testnetClients, err := getTestnetClientsFromExampleConfig()
//	require.NoError(t, err)
//	consensusClient := testnetClients.GetRandomConsensusClient()
//	randomBellarixSignedBeaconBlock := consensus.RandomBellatrixSignedBeaconBlock()
//	fmt.Printf("generated random deneb signed beacon block: %s", randomBellarixSignedBeaconBlock.String())
//	err = consensusClient.PostSignedBeaconBlock(*randomBellarixSignedBeaconBlock)
//	if err != nil {
//		fmt.Printf("client at: %s responded with error: %s", consensusClient.APIEndpoint, err.Error())
//	}
//}
