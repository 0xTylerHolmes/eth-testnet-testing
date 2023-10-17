package eth_testnet_tool

import (
	"math/rand"
	"reflect"
	"time"
)

// ExecutionClientJSON the json representation of an execution layer client
type ExecutionClientJSON struct {
	Name        string `json:"name"`
	RPCEndpoint string `json:"rpc-endpoint"`
}

// ConsensusClientJSON the json representation of a consensus layer client
type ConsensusClientJSON struct {
	Name        string `json:"name"`
	APIEndpoint string `json:"api-endpoint"`
}

// ConsensusClient is the api-endpoint we interact with for a CL
type ConsensusClient struct {
	APIEndpoint string
}

// ExecutionClient is the rpc-endpoint we interact with for an EL
type ExecutionClient struct {
	RPCEndpoint string
}

// TestnetClients a collection of all the consensus and execution clients
// the map is keyed by name with the value being the client endpoint.
type TestnetClients struct {
	ConsensusClients map[string]ConsensusClient
	ExecutionClients map[string]ExecutionClient
}

// TestnetClientsJSON the json representation of the TestnetClients
type TestnetClientsJSON struct {
	ConsensusClients []ConsensusClientJSON `json:"consensus-clients"`
	ExecutionClients []ExecutionClientJSON `json:"execution-clients"`
}

// TestnetConfigJSON the json structure holding testnet-parameters for various utilities
type TestnetConfigJSON struct {
	// consensus related info
	// ValidatorMnemonic the mnemonic used to pre-seed genesis validators and all future added validators
	ValidatorMnemonic string `json:"validator-mnemonic"`
	// DepositContractAddress the deposit contract address, if omitted will fetch from a random client
	DepositContractAddress string `json:"deposit-contract-address,omitempty"`
	// ExecutionAccountMnemonic used to seed premines
	ExecutionAccountMnemonic string `json:"execution-account-mnemonic"`
	// ExecutionPremines the seeded premine addresses and their values
	ExecutionPremines map[string]uint64 `json:"premines"`
}

// GetTestnetClients creates a *TestnetClients that can be used to interact with clients in the testnet.
func (t *TestnetClientsJSON) GetTestnetClients() *TestnetClients {
	var testnetClients TestnetClients
	testnetClients.ConsensusClients = make(map[string]ConsensusClient)
	testnetClients.ExecutionClients = make(map[string]ExecutionClient)
	for _, consensusClientJSON := range t.ConsensusClients {
		testnetClients.ConsensusClients[consensusClientJSON.Name] = ConsensusClient{APIEndpoint: consensusClientJSON.APIEndpoint}
	}
	for _, executionClientJSON := range t.ExecutionClients {
		testnetClients.ExecutionClients[executionClientJSON.Name] = ExecutionClient{RPCEndpoint: executionClientJSON.RPCEndpoint}
	}
	return &testnetClients
}

func (t *TestnetClients) GetRandomConsensusClient() *ConsensusClient {
	// non-efficient but straight forward
	rand.Seed(time.Now().UnixNano())
	keys := reflect.ValueOf(t.ConsensusClients).MapKeys()
	randIdx := rand.Intn(len(keys))
	randomClient := t.ConsensusClients[keys[randIdx].String()]
	return &randomClient
}
