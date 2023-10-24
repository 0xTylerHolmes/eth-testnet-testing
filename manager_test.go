package eth_testnet_tool

import (
	"context"
	"encoding/hex"
	"eth-testnet-tool/consensus_objects"
	"fmt"
	"github.com/attestantio/go-eth2-client/api"
	"github.com/attestantio/go-eth2-client/spec/bellatrix"
	"github.com/attestantio/go-eth2-client/spec/capella"
	"github.com/attestantio/go-eth2-client/spec/phase0"
	"github.com/stretchr/testify/require"
	"testing"
)

const ExampleTestnetClientsConfigFilePath = "./example/configs/example-testnet-clients-config.json"
const ExampleTestnetConfigFilePath = "./example/configs/example-testnet-config.json"

func getTestnetClientManager() (*ClientManager, error) {
	return NewClientManager(ExampleTestnetClientsConfigFilePath, ExampleTestnetConfigFilePath)
}

func TestClientManager_New(t *testing.T) {
	manager, err := getTestnetClientManager()
	require.NoError(t, err)
	for name, client := range manager.Eth2Clients {
		t.Logf("%s: %s\n", name, client.Address())
		t.Log(client.BeaconBlockHeader(context.Background(), &api.BeaconBlockHeaderOpts{Block: "head"}))
	}
}

func TestClientManager_SendRandomAttestations(t *testing.T) {
	manager, err := getTestnetClientManager()
	require.NoError(t, err)
	for name, client := range manager.Eth2Clients {
		randomAttestation := consensus_objects.RandomAttestation()
		fmt.Printf("%s <- %s\n", name, randomAttestation)
		err = client.SubmitAttestations(context.Background(), []*phase0.Attestation{randomAttestation})
		fmt.Printf("%s: %s\n", name, err.Error())
	}
}

func TestClientManager_Validators(t *testing.T) {
	manager, err := getTestnetClientManager()
	require.NoError(t, err)
	t.Log(manager.Validators)
}

func TestClient_VerifyValidators(t *testing.T) {
	manager, err := getTestnetClientManager()
	require.NoError(t, err)
	client := manager.GetRandomConsensusClient()
	onChainValidatorsIndexes, err := client.Validators(context.Background(), &api.ValidatorsOpts{
		State: "head",
	})
	require.NoError(t, err)
	for k, v := range onChainValidatorsIndexes.Data {
		chainKey := v.Validator.PublicKey.String()
		localKey := hex.EncodeToString(manager.Validators[k].ValidatorKey.PublicKey().Marshal())
		require.Equal(t, fmt.Sprintf("%s", chainKey), fmt.Sprintf("0x%s", localKey))
	}
}

func TestClientManager_Configuration(t *testing.T) {
	manager, err := getTestnetClientManager()
	require.NoError(t, err)
	client := manager.GetRandomConsensusClient()
	spec, err := client.Spec(context.Background())
	require.NoError(t, err)
	for k, v := range spec.Data {
		t.Logf("%s: %s\n", k, v)
	}
	fmt.Println(spec.Data)
}

func TestClientManager_SignedBLSToExecutionChange(t *testing.T) {
	manager, err := getTestnetClientManager()
	require.NoError(t, err)
	validatorIndex := uint64(8)
	var withdrawalPubKey phase0.BLSPubKey
	copy(withdrawalPubKey[:], manager.Validators[validatorIndex].WithdrawalKey.PublicKey().Marshal())

	nullAddress := bellatrix.ExecutionAddress{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}

	change := &capella.BLSToExecutionChange{
		ValidatorIndex:     phase0.ValidatorIndex(validatorIndex),
		FromBLSPubkey:      withdrawalPubKey,
		ToExecutionAddress: nullAddress,
	}

	signedChange, err := manager.SignBLSToExecutionChange(change, validatorIndex)
	require.NoError(t, err)
	fmt.Printf("Generated signedBLSToExecutionChange: %s\n", signedChange)
	err = manager.GetRandomConsensusClient().SubmitBLSToExecutionChanges(context.Background(), []*capella.SignedBLSToExecutionChange{signedChange})
	if err != nil {
		fmt.Printf("got a potentially expected error: %s\n", err.Error())
	} else {
		fmt.Printf("submitted without an error.\n")
	}
}

func TestClientManager_CurrentTimeValues(t *testing.T) {
	manager, err := getTestnetClientManager()
	require.NoError(t, err)
	client := manager.GetRandomConsensusClient()
	slotsPerEpoch, err := client.SlotsPerEpoch(context.Background())
	require.NoError(t, err)
	currSlot := manager.GetCurrentSlot()
	currEpoch := manager.GetCurrentEpoch()
	blockHeader, err := client.BeaconBlockHeader(context.Background(), &api.BeaconBlockHeaderOpts{Block: "head"})
	blockSlot := blockHeader.Data.Header.Message.Slot
	require.Equal(t, currSlot, blockSlot)
	require.Equal(t, uint64(currSlot)/slotsPerEpoch, uint64(currEpoch))
}
