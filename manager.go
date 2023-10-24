package eth_testnet_tool

import (
	"context"
	"encoding/json"
	"eth-testnet-tool/validator"
	"github.com/attestantio/go-eth2-client/http"
	"github.com/attestantio/go-eth2-client/spec/phase0"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"math/rand"
	"os"
	"time"
)

type TestnetManager struct {
	ConsensusLayerManager ClientManager
}

type ClientManager struct {
	Eth2Clients   map[string]*http.Service
	TestnetConfig *TestnetConfig
	Validators    []*validator.Validator
	SlotsPerEpoch uint64
	SlotDuration  time.Duration
	GenesisTime   time.Time
}

func NewClientManager(testnetClientsConfigFilePath string, testnetConfigFilePath string) (*ClientManager, error) {
	eth2Clients := make(map[string]*http.Service)
	testnetClients, err := testnetClientsFromFile(testnetClientsConfigFilePath)
	if err != nil {
		return nil, errors.Wrap(err, "unable to create clients from config.")
	}

	for name, clClient := range testnetClients.ConsensusClients {
		service, err := http.New(context.Background(), http.WithAddress(clClient.APIEndpoint), http.WithTimeout(5*time.Second), http.WithLogLevel(zerolog.WarnLevel), http.WithEnforceJSON(true))
		if err != nil {
			return nil, errors.Wrap(err, "unable to create service with clients")
		}
		eth2Clients[name] = service.(*http.Service)
	}
	testnetConfig, err := testnetConfigFromFile(testnetConfigFilePath)
	if err != nil {
		return nil, errors.Wrap(err, "unable to get Testnet Config from file")
	}

	validators, err := validator.GetValidatorsFromMnemonic(testnetConfig.ValidatorMnemonic, 0, testnetConfig.GenesisValidatorCount)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create validators")
	}

	clientManager := ClientManager{
		Eth2Clients:   eth2Clients,
		TestnetConfig: testnetConfig,
		Validators:    validators,
	}

	err = clientManager.setTestnetParameters()
	if err != nil {
		return nil, err
	}
	return &clientManager, nil

}

func (c *ClientManager) setTestnetParameters() error {
	randomClient := c.GetRandomConsensusClient()
	genesisTime, err := randomClient.GenesisTime(context.Background())
	if err != nil {
		return errors.Wrap(err, "failed to get genesis time for testnet")
	}
	slotDuration, err := randomClient.SlotDuration(context.Background())
	if err != nil {
		return errors.Wrap(err, "failed to get slot duration")
	}
	slotsPerEpoch, err := randomClient.SlotsPerEpoch(context.Background())
	if err != nil {
		return errors.Wrap(err, "failed to get slots per epoch")
	}
	c.GenesisTime = genesisTime
	c.SlotDuration = slotDuration
	c.SlotsPerEpoch = slotsPerEpoch

	return nil
}

func (c *ClientManager) GetRandomConsensusClient() *http.Service {
	r := rand.Intn(len(c.Eth2Clients))
	for _, client := range c.Eth2Clients {
		if r == 0 {
			return client
		}
		r--
	}
	panic("unreachable")
}

func (c *ClientManager) GetCurrentEpoch() phase0.Epoch {
	return phase0.Epoch(uint64(c.GetCurrentSlot()) / c.SlotsPerEpoch)
}

func (c *ClientManager) GetCurrentSlot() phase0.Slot {
	return phase0.Slot(uint64(time.Since(c.GenesisTime).Seconds()) / uint64(c.SlotDuration.Seconds()))
}

func testnetClientsFromFile(filePath string) (*TestnetClients, error) {
	var testnetClientsJSON TestnetClientsJSON
	var testnetClients TestnetClients

	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, errors.Wrap(err, "failed to open the testnet-clients config")
	}

	err = json.Unmarshal(data, &testnetClientsJSON)
	if err != nil {
		return nil, errors.Wrap(err, "failed to unmarshall the testnet-clients config")
	}

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
