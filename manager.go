package eth_testnet_tool

import (
	"context"
	"encoding/json"
	"github.com/attestantio/go-eth2-client/http"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"os"
	"time"
)

type TestnetManager struct {
	ConsensusLayerManager ClientManager
}

type ClientManager struct {
	Eth2Clients   map[string]*http.Service
	TestnetConfig *TestnetConfig
}

func NewClientManager(testnetClientsConfigFilePath string, testnetConfigFilePath string) (*ClientManager, error) {
	eth2Clients := make(map[string]*http.Service)
	testnetClients, err := testnetClientsFromFile(testnetClientsConfigFilePath)
	if err != nil {
		return nil, errors.Wrap(err, "unable to create clients from config.")
	}
	for name, client := range testnetClients.ConsensusClients {
		service, err := http.New(context.Background(), http.WithAddress(client.APIEndpoint), http.WithTimeout(5*time.Second), http.WithLogLevel(zerolog.WarnLevel))
		if err != nil {
			return nil, errors.Wrap(err, "unable to create service with clients")
		}
		eth2Clients[name] = service.(*http.Service)
	}

	testnetConfig, err := testnetConfigFromFile(testnetConfigFilePath)
	if err != nil {
		return nil, errors.Wrap(err, "unable to get Testnet Config from file")
	}
	return &ClientManager{
		Eth2Clients:   eth2Clients,
		TestnetConfig: testnetConfig,
	}, nil

}

func testnetClientsFromFile(filePath string) (*TestnetClients, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, errors.Wrap(err, "failed to open the testnet-clients config")
	}
	var testnetClientsJSON TestnetClientsJSON
	err = json.Unmarshal(data, &testnetClientsJSON)
	if err != nil {
		return nil, errors.Wrap(err, "failed to unmarshall the testnet-clients config")
	}
	var testnetClients TestnetClients
	testnetClients.ConsensusClients = make(map[string]ConsensusClient)
	testnetClients.ExecutionClients = make(map[string]ExecutionClient)
	for _, consensusClientJSON := range testnetClientsJSON.ConsensusClients {
		testnetClients.ConsensusClients[consensusClientJSON.Name] = ConsensusClient{APIEndpoint: consensusClientJSON.APIEndpoint}
	}
	for _, executionClientJSON := range testnetClientsJSON.ExecutionClients {
		testnetClients.ExecutionClients[executionClientJSON.Name] = ExecutionClient{RPCEndpoint: executionClientJSON.RPCEndpoint}
	}
	return &testnetClients, nil
}

func testnetConfigFromFile(filePath string) (*TestnetConfig, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, errors.Wrap(err, "failed to open the testnet-config file")
	}
	var testnetConfig TestnetConfigJSON
	err = json.Unmarshal(data, &testnetConfig)
	if err != nil {
		return nil, errors.Wrap(err, "failed to unmarshal the testnet-config")
	}
	return &TestnetConfig{
		ValidatorMnemonic:        testnetConfig.ValidatorMnemonic,
		GenesisValidatorCount:    testnetConfig.GenesisValidatorCount,
		DepositContractAddress:   testnetConfig.DepositContractAddress,
		ExecutionAccountMnemonic: testnetConfig.ExecutionAccountMnemonic,
		ExecutionPremines:        testnetConfig.ExecutionPremines,
	}, nil
}
