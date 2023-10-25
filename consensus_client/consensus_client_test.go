package consensus_client

import (
	"context"
	"encoding/hex"
	"eth-testnet-tool/validator"
	"fmt"
	"github.com/attestantio/go-eth2-client/http"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

const (
	Timeout               = 3 * time.Second
	ForceJson             = true // force json for testing with minimal testnets
	BeaconAPI             = "http://10.0.20.5:5052"
	ValidatorMnemonic     = "ocean style run case glory clip into nature guess jacket document firm fiscal hello kite disagree symptom tide net coral envelope wink render festival"
	GenesisValidatorCount = 20
)

func getTestConsensusClient() (*ConsensusClient, error) {
	service, err := http.New(context.Background(), http.WithAddress(BeaconAPI), http.WithTimeout(Timeout), http.WithEnforceJSON(ForceJson))
	if err != nil {
		return nil, errors.Wrap(err, "unable to get test client for test")
	}

	return &ConsensusClient{
		Name:          "testClient",
		BeaconAPI:     BeaconAPI,
		BeaconService: service.(*http.Service),
	}, nil
}

func getTestValidators(maxAcc uint64) ([]*validator.Validator, error) {
	validators, err := validator.GetValidatorsFromMnemonic(ValidatorMnemonic, 0, maxAcc)
	if err != nil {
		return nil, errors.Wrap(err, "failed to populate test validators")
	}
	return validators, nil
}

func TestConsensusClient_GetValidators(t *testing.T) {
	testConsensusClient, err := getTestConsensusClient()
	require.NoError(t, err)
	genesisValidators, err := getTestValidators(GenesisValidatorCount)
	require.NoError(t, err)
	clientValidators, err := testConsensusClient.GetAllValidators("head")
	require.NoError(t, err)
	for ndx, v := range clientValidators {
		require.Equal(t, genesisValidators[ndx].ValidatorIndex, uint64(ndx))
		require.Equal(t, fmt.Sprintf("0x%s", hex.EncodeToString(genesisValidators[ndx].ValidatorKey.PublicKey().Marshal())), v.Validator.PublicKey.String())
	}
}

func TestConsensusClient_GetAllActiveValidators(t *testing.T) {
	testConsensusClient, err := getTestConsensusClient()
	require.NoError(t, err)
	clientValidators, err := testConsensusClient.GetAllActiveValidators("head")
	require.NoError(t, err)
	for _, v := range clientValidators {
		t.Log(v.String())
		require.Equal(t, v.Status.IsActive(), true)
	}
}

func TestConsensusClient_GetValidatorByPublicKey(t *testing.T) {
	testConsensusClient, err := getTestConsensusClient()
	require.NoError(t, err)
	genesisValidators, err := getTestValidators(GenesisValidatorCount)
	require.NoError(t, err)
	validator, err := testConsensusClient.GetValidatorByPublicKey("head", genesisValidators[4].ValidatorPublicKey)
	require.Equal(t, validator.Validator.PublicKey, genesisValidators[4].ValidatorPublicKey)
	require.NoError(t, err)
}

func TestConsensusClient_GetShardCommitteePeriod(t *testing.T) {
	testConsensusClient, err := getTestConsensusClient()
	require.NoError(t, err)
	shardCommitteePeriod, err := testConsensusClient.GetShardCommitteePeriod()
	require.NoError(t, err)
	t.Logf("got shard committee period: %d from %s\n", shardCommitteePeriod, testConsensusClient.String())
}
