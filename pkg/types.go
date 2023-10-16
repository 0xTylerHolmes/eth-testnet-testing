package pkg

// ExecutionClientJSON the json representation of an execution layer client
type ExecutionClientJSON struct {
	Name        string `json:"name"`
	RPCEndpoint string `json:"rpc-endpoint"`
}

type ConsensusClientJSON struct {
	Name        string `json:"name"`
	APIEndpoint string `json:"api-endpoint"`
}

type ConsensusClient struct {
	APIEndpoint string
}

type ExecutionClient struct {
	RPCEndpoint string
}

type TestnetClients struct {
	ConsensusClients map[string]ConsensusClient
	ExecutionClients map[string]ExecutionClient
}

type TestnetClientsJSON struct {
	ConsensusClients []ConsensusClientJSON `json:"consensus-clients"`
	ExecutionClients []ExecutionClientJSON `json:"execution-clients"`
}

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
