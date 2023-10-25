package eth_testnet_tool

import (
	"context"
	"eth-testnet-tool/consensus_client"
	"eth-testnet-tool/consensus_client/consensus_objects"
	"eth-testnet-tool/validator"
	"github.com/attestantio/go-eth2-client/http"
	"github.com/attestantio/go-eth2-client/spec/bellatrix"
	"github.com/attestantio/go-eth2-client/spec/phase0"
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

func getTestConsensusClient() (*consensus_client.ConsensusClient, error) {
	service, err := http.New(context.Background(), http.WithAddress(BeaconAPI), http.WithTimeout(Timeout), http.WithEnforceJSON(ForceJson))
	if err != nil {
		return nil, errors.Wrap(err, "unable to get test client for test")
	}

	return &consensus_client.ConsensusClient{
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

func TestCreateAndSubmitBLSToExecutionChange(t *testing.T) {
	// RUN on minimal to see the effect.
	testConsensusClient, err := getTestConsensusClient()
	require.NoError(t, err)
	validators, err := getTestValidators(20)
	require.NoError(t, err)
	randomAddress := bellatrix.ExecutionAddress{0x69, 0x69, 0x069, 0x69, 0x69, 0x69, 0x69, 0x69, 0x69, 0x69, 0x69, 0x69, 0x69, 0x69, 0x69, 0x69, 0x69, 0x69, 0x69, 0x69}
	testValidator := validators[14]

	signedBlsToExecutionChange, err := BuildSignedBLSToExecutionChangeFromClientView(testConsensusClient, testValidator, randomAddress)
	require.NoError(t, err)

	err = testConsensusClient.SubmitBLSToExecutionChange(signedBlsToExecutionChange)
	if err != nil {
		t.Log(err.Error())
	}
	require.NoError(t, err)
}

func TestCreateAndSubmitInvalidBLSToExecutionChange(t *testing.T) {
	testConsensusClient, err := getTestConsensusClient()
	require.NoError(t, err)
	validators, err := getTestValidators(20)
	require.NoError(t, err)
	testValidator := validators[14]

	blsToExecutionChange := consensus_objects.RandomBLSToExecutionChange()
	signedBLSToExecutionChange, err := SignBLSToExecutionChangeWithValidator(testConsensusClient, testValidator, blsToExecutionChange)
	require.NoError(t, err)

	err = testConsensusClient.SubmitBLSToExecutionChange(signedBLSToExecutionChange)
	if err != nil {
		t.Log(err.Error())
	}
	require.Error(t, err)
}

func TestCreateAndSubmitSignedVoluntaryExit(t *testing.T) {
	testConsensusClient, err := getTestConsensusClient()
	require.NoError(t, err)
	validators, err := getTestValidators(20)
	require.NoError(t, err)
	testValidator := validators[14]

	shardCommitteePeriod, err := testConsensusClient.GetShardCommitteePeriod()
	require.NoError(t, err)

	currentEpoch, err := testConsensusClient.GetCurrentEpoch()
	require.NoError(t, err)

	if uint64(currentEpoch) < shardCommitteePeriod {
		t.Log("We are going to fail in an expected way. The testnet you are using for testing has not reached epoch > SHARD_COMMITTEE_PERIOD")
		t.Logf("current epoch: %d, shard_committee_period: %d", currentEpoch, shardCommitteePeriod)
		require.NoError(t, errors.New("testnet has not progressed long enough to test voluntary exits"))
	}

	voluntaryExit, err := BuildVoluntaryExitFromClientView(testConsensusClient, testValidator, phase0.Epoch(shardCommitteePeriod))
	require.NoError(t, err)

	err = testConsensusClient.SubmitValidatorExit(voluntaryExit)
	require.NoError(t, err)
}
