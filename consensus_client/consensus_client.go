package consensus_client

import (
	"context"
	"fmt"
	"github.com/attestantio/go-eth2-client/api"
	v1 "github.com/attestantio/go-eth2-client/api/v1"
	eth2client "github.com/attestantio/go-eth2-client/http"
	"github.com/attestantio/go-eth2-client/spec/capella"
	"github.com/attestantio/go-eth2-client/spec/phase0"
	"github.com/pkg/errors"
)

type ConsensusClient struct {
	Name          string
	BeaconAPI     string
	BeaconService *eth2client.Service
}

func (c *ConsensusClient) String() string {
	return fmt.Sprintf("%s @ %s", c.Name, c.BeaconAPI)
}

// Information about spec

// GetDomainTypeFromSpec returns the domain type from the clients spec
func (c *ConsensusClient) GetDomainTypeFromSpec(domainTypeString string) (phase0.DomainType, error) {
	spec, err := c.BeaconService.Spec(context.Background())
	if err != nil {
		return phase0.DomainType{}, errors.Wrap(err, "failed to get spec for domain calculation")
	}

	d, ok := spec.Data[domainTypeString]
	if !ok {
		return phase0.DomainType{}, errors.New("failed to find domain type from spec")
	}
	domainType := d.(phase0.DomainType)
	return domainType, nil
}

// GetDomainFromGenesis fetches the genesis domain for the domain type
func (c *ConsensusClient) GetDomainFromGenesis(domainType phase0.DomainType) (phase0.Domain, error) {
	domain, err := c.BeaconService.GenesisDomain(context.Background(), domainType)
	if err != nil {
		return phase0.Domain{}, errors.Wrap(err, "failed to get domain from consensus client")
	}
	return domain, nil
}

func (c *ConsensusClient) GetSlotsPerEpoch() (uint64, error) {
	slotsPerEpoch, err := c.BeaconService.SlotsPerEpoch(context.Background())
	if err != nil {
		return 0, err
	}
	return slotsPerEpoch, nil
}

func (c *ConsensusClient) GetShardCommitteePeriod() (uint64, error) {
	spec, err := c.BeaconService.Spec(context.Background())
	if err != nil {
		return 0, errors.Wrap(err, "failed to get spec for shard committee period")
	}
	shardCommitteePeriod, ok := spec.Data["SHARD_COMMITTEE_PERIOD"].(uint64)
	if !ok {
		return 0, errors.New("failed to get shard committee period from spec")
	}
	return shardCommitteePeriod, nil
}

func (c *ConsensusClient) GetCurrentEpoch() (phase0.Epoch, error) {
	currentHeader, err := c.GetBlockHeader("head")
	if err != nil {
		return phase0.Epoch(0), errors.Wrap(err, "failed to get current header to determine epoch")
	}
	slotsPerEpoch, err := c.GetSlotsPerEpoch()
	if err != nil {
		return phase0.Epoch(0), errors.Wrap(err, "failed to fetch slots per epoch from client spec")
	}
	return phase0.Epoch(uint64(currentHeader.Header.Message.Slot) / slotsPerEpoch), nil
}

// GetBlockHeader returns the clients view of the block header at the provided block_id
func (c *ConsensusClient) GetBlockHeader(block string) (*v1.BeaconBlockHeader, error) {
	resp, err := c.BeaconService.BeaconBlockHeader(context.Background(), &api.BeaconBlockHeaderOpts{Block: block})
	if err != nil {
		return nil, errors.Wrapf(err, "failed to get block header for client: %s", c.Name)
	}
	return resp.Data, nil
}

// getValidators returns the view of the validators from the perspective of this client at a given state
// this is used by some of the wrappers to remove boilerplate code in testnet experiments
func (c *ConsensusClient) getValidators(opts *api.ValidatorsOpts) (map[phase0.ValidatorIndex]*v1.Validator, error) {
	resp, err := c.BeaconService.Validators(context.Background(), opts)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to get validators for client: %s", c.Name)
	}
	return resp.Data, nil
}

// GetAllValidators returns all the validators in view of the client at the provided state (head/genesis/finalized/justified/slot/0xstateRoot)
func (c *ConsensusClient) GetAllValidators(stateID string) (map[phase0.ValidatorIndex]*v1.Validator, error) {
	resp, err := c.BeaconService.Validators(context.Background(), &api.ValidatorsOpts{State: stateID})
	if err != nil {
		return nil, errors.Wrapf(err, "failed to get validators for client: %s", c.Name)
	}
	return resp.Data, nil
}

// GetAllActiveValidators fetches all the validators by a pending state
func (c *ConsensusClient) GetAllActiveValidators(stateID string) (map[phase0.ValidatorIndex]*v1.Validator, error) {
	activeValidators := make(map[phase0.ValidatorIndex]*v1.Validator)
	allValidators, err := c.getValidators(&api.ValidatorsOpts{State: stateID})
	if err != nil {
		return nil, err
	}
	for _, validator := range allValidators {
		if validator.Status.IsActive() {
			activeValidators[validator.Index] = validator
		}
	}
	return activeValidators, nil
}

// GetValidatorByPublicKey fetches the validator entry with the same public key from the client
func (c *ConsensusClient) GetValidatorByPublicKey(stateID string, pubKey phase0.BLSPubKey) (*v1.Validator, error) {
	resp, err := c.getValidators(&api.ValidatorsOpts{
		State:   stateID,
		PubKeys: []phase0.BLSPubKey{pubKey},
	})
	if err != nil {
		return nil, err
	}
	if len(resp) == 0 {
		return nil, fmt.Errorf("this client does not have a view containg a validator with public key: %s", pubKey.String())
	}
	if len(resp) > 1 {
		return nil, fmt.Errorf("unexepcted error with multiple validators with public key: %s", pubKey.String())
	}
	for _, v := range resp {
		return v, nil
	}
	panic("unreachable")
}

// SubmitBLSToExecutionChange submits the SignedBLSToExecutionChange to the clients bls to execution change pool
func (c *ConsensusClient) SubmitBLSToExecutionChange(change *capella.SignedBLSToExecutionChange) error {
	return c.BeaconService.SubmitBLSToExecutionChanges(context.Background(), []*capella.SignedBLSToExecutionChange{change})
}

// SubmitValidatorExit submits the voluntary exit to the clients validator exit pool.
func (c *ConsensusClient) SubmitValidatorExit(exit *phase0.SignedVoluntaryExit) error {
	return c.BeaconService.SubmitVoluntaryExit(context.Background(), exit)
}
