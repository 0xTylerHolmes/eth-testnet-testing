package eth_testnet_tool

import (
	"context"
	"eth-testnet-tool/consensus_client"
	"eth-testnet-tool/validator"
	"fmt"
	"github.com/attestantio/go-eth2-client/spec/bellatrix"
	"github.com/attestantio/go-eth2-client/spec/capella"
	"github.com/attestantio/go-eth2-client/spec/phase0"
	"github.com/pkg/errors"
	"github.com/protolambda/zrnt/eth2/beacon/common"
)

var (
	BlsToExecutionChangeDomainLookup = "DOMAIN_BLS_TO_EXECUTION_CHANGE"
	VoluntaryExitDomainLookup        = "DOMAIN_VOLUNTARY_EXIT"
)

// Various useful operations for testnet testing
// The following should always work under healthy network conditions

// BuildVoluntaryExitFromClientView  uses the specified consensus client to build a SignedVoluntaryExit message
func BuildVoluntaryExitFromClientView(consensusClient *consensus_client.ConsensusClient, validator *validator.Validator, epoch phase0.Epoch) (*phase0.SignedVoluntaryExit, error) {
	clientValidatorView, err := consensusClient.GetValidatorByPublicKey("head", validator.ValidatorPublicKey)
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("couldn't create bls to execution change from client %s", consensusClient.Name))
	}

	shardCommitteePeriod, err := consensusClient.GetShardCommitteePeriod()
	if err != nil {
		return nil, errors.Wrap(err, "couldn't get shard committee period to check for valid exit epoch")
	}
	if uint64(epoch) < shardCommitteePeriod {
		return nil, errors.New("can not submit an exit before the SHARD_COMMITTEE_PERIOD")
	}

	voluntaryExit := phase0.VoluntaryExit{
		Epoch:          epoch,
		ValidatorIndex: clientValidatorView.Index,
	}

	return SignVoluntaryExitWithValidator(consensusClient, validator, &voluntaryExit)
}

// BuildSignedBLSToExecutionChangeFromClientView uses the specified consensus client to create and sign a BLSToExecutionChange
// This will always create a valid BLSToExecutionChange with respect to this client.
func BuildSignedBLSToExecutionChangeFromClientView(consensusClient *consensus_client.ConsensusClient, validator *validator.Validator, address bellatrix.ExecutionAddress) (*capella.SignedBLSToExecutionChange, error) {

	clientValidatorView, err := consensusClient.GetValidatorByPublicKey("head", validator.ValidatorPublicKey)
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("couldn't create bls to execution change from client %s", consensusClient.Name))
	}

	if clientValidatorView.Validator.PublicKey != validator.ValidatorPublicKey {
		return nil, errors.New("Could not create a signedBLSToExecutionChange, clients validator doesn't match the supplied one.")
	}
	var blsToExecutionChange = &capella.BLSToExecutionChange{
		ValidatorIndex:     clientValidatorView.Index,
		FromBLSPubkey:      phase0.BLSPubKey(validator.WithdrawalKey.PublicKey().Marshal()),
		ToExecutionAddress: address,
	}

	return SignBLSToExecutionChangeWithValidator(consensusClient, validator, blsToExecutionChange)
}

// Signing methods allow you to sign with the wrong key for testing purposes.

// SignBLSToExecutionChangeWithValidator does an unverified sign with the validator on the supplied execution change
// WARN: using the wrong validator can lead to an invalid signature.  If this is not your intentions use BuildSignedBLSToExecutionChangeFromClientView to create a valid signed payload.
func SignBLSToExecutionChangeWithValidator(consensusClient *consensus_client.ConsensusClient, validator *validator.Validator, blsToExecutionChange *capella.BLSToExecutionChange) (*capella.SignedBLSToExecutionChange, error) {
	var signedBLSToExecutionChange = capella.SignedBLSToExecutionChange{
		Message:   blsToExecutionChange,
		Signature: phase0.BLSSignature{},
	}
	domainType, err := consensusClient.GetDomainTypeFromSpec(BlsToExecutionChangeDomainLookup)
	if err != nil {
		return nil, err
	}
	domain, err := consensusClient.GetDomainFromGenesis(domainType)
	if err != nil {
		return nil, err
	}
	messageRoot, err := blsToExecutionChange.HashTreeRoot()
	if err != nil {
		return nil, err
	}
	signingRoot := common.ComputeSigningRoot(messageRoot, common.BLSDomain(domain))
	signature := validator.WithdrawalKey.Sign(signingRoot[:])
	copy(signedBLSToExecutionChange.Signature[:], signature.Marshal()[:])

	return &signedBLSToExecutionChange, nil
}

// SignVoluntaryExitWithValidator sign a volunatry exit with a validator
// WARN: using the wrong validator can lead to an invalid signature. If this is not your intentions use BuildVoluntaryExitFromClientView to create a valid signed payload.
func SignVoluntaryExitWithValidator(consensusClient *consensus_client.ConsensusClient, validator *validator.Validator, voluntaryExit *phase0.VoluntaryExit) (*phase0.SignedVoluntaryExit, error) {

	var signedVoluntaryExit = phase0.SignedVoluntaryExit{
		Message:   voluntaryExit,
		Signature: phase0.BLSSignature{},
	}

	operationRoot, err := voluntaryExit.HashTreeRoot()
	if err != nil {
		return nil, errors.Wrap(err, "failed to get message hash tree root")
	}

	domainType, err := consensusClient.GetDomainTypeFromSpec(VoluntaryExitDomainLookup)
	if err != nil {
		return nil, err
	}
	domain, err := consensusClient.BeaconService.Domain(context.Background(), domainType, voluntaryExit.Epoch)
	if err != nil {
		return nil, err
	}
	signingRoot := common.ComputeSigningRoot(operationRoot, common.BLSDomain(domain))
	signature := validator.ValidatorKey.Sign(signingRoot[:])
	copy(signedVoluntaryExit.Signature[:], signature.Marshal()[:])
	return &signedVoluntaryExit, nil
}
