package eth_testnet_tool

import (
	"encoding/json"
	"github.com/pkg/errors"
	"os"
)

func GetTestnetClientsFromConfig(filePath string) (*TestnetClientsJSON, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, errors.Wrap(err, "failed to open the testnet-clients config")
	}
	var testnetClients TestnetClientsJSON
	err = json.Unmarshal(data, &testnetClients)
	if err != nil {
		return nil, errors.Wrap(err, "failed to unmarshall the testnet-clients config")
	}
	return &testnetClients, nil
}
