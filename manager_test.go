package eth_testnet_tool

import (
	"context"
	"encoding/hex"
	"fmt"
	"github.com/attestantio/go-eth2-client/api"
	"github.com/stretchr/testify/require"
	"testing"
)

const ExampleTestnetClientsConfigFilePath = "./example/configs/example-testnet-clients-config.json"
const ExampleTestnetConfigFilePath = "./example/configs/example-testnet-config.json"

func getTestnetClientManager() (*ClientManager, error) {
	return NewClientManager(ExampleTestnetClientsConfigFilePath, ExampleTestnetConfigFilePath)
}

func TestClientManager_Validators(t *testing.T) {
	manager, err := getTestnetClientManager()
	require.NoError(t, err)
	t.Log(manager.Validators)
}

func TestClient_VerifyValidators(t *testing.T) {
	//WARN: some clients don't support this endpoint when running minimals
	manager, err := getTestnetClientManager()
	require.NoError(t, err)
	for _, consensusClient := range manager.ConsensusClients {
		// some consensus clients don't support this endpoint.
		onChainValidatorsIndexes, err := consensusClient.BeaconService.Validators(context.Background(), &api.ValidatorsOpts{
			State: "head",
		})
		if err != nil {
			fmt.Printf("Warning %s doesn't support this validators call (%s)\n", consensusClient.String(), err.Error())
		} else {
			fmt.Printf("Checking validators from client: %s\n", consensusClient.String())
			for k, v := range onChainValidatorsIndexes.Data {
				chainKey := v.Validator.PublicKey.String()
				localKey := hex.EncodeToString(manager.Validators[k].ValidatorKey.PublicKey().Marshal())
				require.Equal(t, chainKey, fmt.Sprintf("0x%s", localKey))
			}
		}
	}
}
