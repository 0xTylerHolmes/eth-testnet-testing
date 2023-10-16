package pkg

import (
	"fmt"
	v1 "github.com/attestantio/go-eth2-client/api/v1"
	"github.com/stretchr/testify/require"
	"testing"
)

func getTestnetClientsForTest() (*TestnetClients, error) {
	testnetClientsJSON, err := GetTestnetClientsFromConfig("/0xtylerholmes/git/eth-testnet-tool/eth-mon-config.json")
	if err != nil {
		return nil, err
	}
	return testnetClientsJSON.GetTestnetClients(), nil
}

func TestGetBlockHeaders(t *testing.T) {
	results := make(map[string]*v1.BeaconBlockHeader)
	clients, err := getTestnetClientsForTest()
	require.NoError(t, err)
	for name, client := range clients.ConsensusClients {
		bbh, err := client.GetBeaconBlockHeader("head")
		if err != nil {
			fmt.Println(err)
			results[name] = nil
		} else {
			results[name] = bbh
		}
	}
	fmt.Println(results)
}
