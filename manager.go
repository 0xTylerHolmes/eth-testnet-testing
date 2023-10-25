package eth_testnet_tool

import (
	"context"
	"encoding/json"
	"eth-testnet-tool/consensus_client"
	"eth-testnet-tool/execution_client"
	"eth-testnet-tool/validator"
	"fmt"
	"github.com/attestantio/go-eth2-client/http"
	"github.com/attestantio/go-eth2-client/spec/phase0"
	"github.com/attestantio/go-execution-client/jsonrpc"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"math/rand"
	"os"
	"time"
)

type ClientManager struct {
	ConsensusClients map[string]*consensus_client.ConsensusClient
	ExecutionClients map[string]*execution_client.ExecutionClient
	TestnetConfig    *TestnetConfig
	Validators       []*validator.Validator
	SlotsPerEpoch    uint64
	SlotDuration     time.Duration
	GenesisTime      time.Time
}

func NewClientManager(testnetClientsConfigFilePath string, testnetConfigFilePath string) (*ClientManager, error) {
	consensusClients, err := getConsensusClientsFromFile(testnetClientsConfigFilePath, 5*time.Second, zerolog.WarnLevel)
	if err != nil {
		return nil, errors.Wrap(err, "unable to create consensus clients from config.")
	}

	executionClients, err := getExecutionClientsFromFile(testnetClientsConfigFilePath, 5*time.Second, zerolog.WarnLevel)
	if err != nil {
		return nil, errors.Wrap(err, "unable to create execution clients from config.")
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
		ConsensusClients: consensusClients,
		ExecutionClients: executionClients,
		TestnetConfig:    testnetConfig,
		Validators:       validators,
	}

	err = clientManager.setTestnetParameters()
	if err != nil {
		return nil, err
	}
	return &clientManager, nil

}

// setTestnetParameters reads the config from a random client and populates the local testnet params.
func (c *ClientManager) setTestnetParameters() error {
	randomClient := c.GetRandomConsensusClient()
	genesisTime, err := randomClient.BeaconService.GenesisTime(context.Background())
	if err != nil {
		return errors.Wrap(err, "failed to get genesis time for testnet")
	}
	slotDuration, err := randomClient.BeaconService.SlotDuration(context.Background())
	if err != nil {
		return errors.Wrap(err, "failed to get slot duration")
	}
	slotsPerEpoch, err := randomClient.BeaconService.SlotsPerEpoch(context.Background())
	if err != nil {
		return errors.Wrap(err, "failed to get slots per epoch")
	}
	c.GenesisTime = genesisTime
	c.SlotDuration = slotDuration
	c.SlotsPerEpoch = slotsPerEpoch

	return nil
}

func (c *ClientManager) GetRandomConsensusClient() *consensus_client.ConsensusClient {
	r := rand.Intn(len(c.ConsensusClients))
	for _, client := range c.ConsensusClients {
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

func getExecutionClientsFromFile(filePath string, timeout time.Duration, logLevel zerolog.Level) (map[string]*execution_client.ExecutionClient, error) {
	var executionTestnetClients = make(map[string]*execution_client.ExecutionClient)
	var testnetClientsJSON TestnetClientsJSON
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, errors.Wrap(err, "failed to open the testnet-clients config")
	}

	err = json.Unmarshal(data, &testnetClientsJSON)
	if err != nil {
		return nil, errors.Wrap(err, "failed to unmarshall the testnet-clients config")
	}

	for _, executionClient := range testnetClientsJSON.ExecutionClients {
		service, err := jsonrpc.New(context.Background(), jsonrpc.WithAddress(executionClient.RPCEndpoint), jsonrpc.WithLogLevel(logLevel), jsonrpc.WithTimeout(timeout))
		if err != nil {
			return nil, errors.Wrap(err, fmt.Sprintf("failed to get client: %s at %s", executionClient.Name, executionClient.RPCEndpoint))
		}
		executionTestnetClients[executionClient.Name] = &execution_client.ExecutionClient{
			Name:       executionClient.Name,
			JsonRPC:    executionClient.RPCEndpoint,
			RPCService: service.(*jsonrpc.Service),
		}
	}
	return executionTestnetClients, nil
}

func getConsensusClientsFromFile(filePath string, timeout time.Duration, logLevel zerolog.Level) (map[string]*consensus_client.ConsensusClient, error) {
	var consensusTestnetClients = make(map[string]*consensus_client.ConsensusClient)
	var testnetClientsJSON TestnetClientsJSON
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, errors.Wrap(err, "failed to open the testnet-clients config")
	}

	err = json.Unmarshal(data, &testnetClientsJSON)
	if err != nil {
		return nil, errors.Wrap(err, "failed to unmarshall the testnet-clients config")
	}

	for _, consensusClient := range testnetClientsJSON.ConsensusClients {
		service, err := http.New(context.Background(), http.WithAddress(consensusClient.APIEndpoint), http.WithLogLevel(logLevel), http.WithTimeout(timeout), http.WithEnforceJSON(true))
		if err != nil {
			return nil, errors.Wrap(err, fmt.Sprintf("failed to get client: %s at %s", consensusClient.Name, consensusClient.APIEndpoint))
		}
		consensusTestnetClients[consensusClient.Name] = &consensus_client.ConsensusClient{
			Name:          consensusClient.Name,
			BeaconAPI:     consensusClient.APIEndpoint,
			BeaconService: service.(*http.Service),
		}
	}
	return consensusTestnetClients, nil
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
