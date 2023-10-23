package eth_testnet_tool

import (
	"context"
	"eth-testnet-tool/consensus"
	"fmt"
	"github.com/attestantio/go-eth2-client/api"
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
	manager, err := NewClientManager(ExampleTestnetClientsConfigFilePath, ExampleTestnetConfigFilePath)
	require.NoError(t, err)
	for name, client := range manager.Eth2Clients {
		fmt.Println(name)
		fmt.Println(client.Address())
		fmt.Println(client.BeaconBlockHeader(context.Background(), &api.BeaconBlockHeaderOpts{Block: "head"}))
	}
}

func TestClientManager_SendRandomAttestations(t *testing.T) {
	manager, err := getTestnetClientManager()
	require.NoError(t, err)
	for name, client := range manager.Eth2Clients {
		randomAttestation := consensus.RandomAttestation()
		fmt.Printf("%s <- %s\n", name, randomAttestation)
		err = client.SubmitAttestations(context.Background(), []*phase0.Attestation{randomAttestation})
		fmt.Printf("%s: %s\n", name, err.Error())
	}
}
