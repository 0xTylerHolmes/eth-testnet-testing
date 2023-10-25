package execution_client

import (
	"fmt"
	ethclient "github.com/attestantio/go-execution-client/jsonrpc"
)

type ExecutionClient struct {
	Name       string
	JsonRPC    string
	RPCService *ethclient.Service
}

func (e *ExecutionClient) String() string {
	return fmt.Sprintf("%s @ %s", e.Name, e.JsonRPC)
}
