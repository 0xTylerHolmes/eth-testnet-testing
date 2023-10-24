package eth_testnet_tool

import (
	"context"
	"fmt"
	"github.com/attestantio/go-eth2-client/spec/capella"
	"github.com/attestantio/go-eth2-client/spec/phase0"
	"github.com/pkg/errors"
	"github.com/protolambda/zrnt/eth2/beacon/common"
)

func (c *ClientManager) SignBLSToExecutionChange(change *capella.BLSToExecutionChange, signerValidatorIndex uint64) (*capella.SignedBLSToExecutionChange, error) {
	client := c.GetRandomConsensusClient()

	var signedBLSToExecutionChange = capella.SignedBLSToExecutionChange{
		Message:   change,
		Signature: phase0.BLSSignature{},
	}

	spec, err := client.Spec(context.Background())
	if err != nil {
		return nil, errors.Wrap(err, "failed to get spec for domain calculation")
	}

	domainType := spec.Data["DOMAIN_BLS_TO_EXECUTION_CHANGE"].(phase0.DomainType)
	exitRoot, err := change.HashTreeRoot()
	if err != nil {
		return nil, errors.Wrap(err, "failed to get hash tree root of exit")
	}
	// use genesis domain for BLSToExecutionChange
	domain, err := client.GenesisDomain(context.Background(), domainType)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get domain")
	}

	if uint64(len(c.Validators)) < signerValidatorIndex {
		return nil, fmt.Errorf("supplied validator index out of range")
	}
	signingRoot := common.ComputeSigningRoot(exitRoot, common.BLSDomain(domain))
	signature := c.Validators[signerValidatorIndex].WithdrawalKey.Sign(signingRoot[:])
	copy(signedBLSToExecutionChange.Signature[:], signature.Marshal()[:])
	return &signedBLSToExecutionChange, nil
}
