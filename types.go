package eth_testnet_tool

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
	ValidatorMnemonic        string            `json:"validator-mnemonic"`
	GenesisValidatorCount    uint64            `json:"genesis-validator-count,omitempty"`
	DepositContractAddress   string            `json:"deposit-contract-address,omitempty"`
	ExecutionAccountMnemonic string            `json:"execution-account-mnemonic"`
	ExecutionPremines        map[string]uint64 `json:"premines"`
}

// TestnetConfig contains information about the running testnet.
type TestnetConfig struct {
	// consensus related info
	// ValidatorMnemonic the mnemonic used to pre-seed genesis validators and all future added validators
	ValidatorMnemonic string
	// GenesisValidatorCount how many validators were pre-seeded into genesis
	GenesisValidatorCount uint64
	// DepositContractAddress the deposit contract address, if omitted will fetch from a random client
	DepositContractAddress string
	// ExecutionAccountMnemonic used to seed premines
	ExecutionAccountMnemonic string
	// ExecutionPremines the seeded premine addresses and their values
	ExecutionPremines map[string]uint64
}
