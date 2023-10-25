package validator

import (
	"encoding/hex"
	"fmt"
	"github.com/attestantio/go-eth2-client/spec/phase0"
	hbls "github.com/herumi/bls-eth-go-binary/bls"
	"github.com/pkg/errors"
	"github.com/tyler-smith/go-bip39"
	e2types "github.com/wealdtech/go-eth2-types/v2"
	util "github.com/wealdtech/go-eth2-util"
	"strings"
)

func init() {
	_ = hbls.Init(hbls.BLS12_381)
	_ = hbls.SetETHmode(hbls.EthModeLatest)
}

type Validator struct {
	ValidatorIndex uint64
	ValidatorKey   e2types.PrivateKey
	WithdrawalKey  e2types.PrivateKey
	// ValidatorPublicKey contains a type friendly version of the public key for go-eth2-client
	ValidatorPublicKey phase0.BLSPubKey
}

func (v *Validator) String() string {
	return fmt.Sprintf("pubKey: 0x%s", hex.EncodeToString(v.ValidatorKey.PublicKey().Marshal()))
}

func MnemonicToSeed(mnemonic string) ([]byte, error) {
	mnem := strings.TrimSpace(mnemonic)
	if bip39.IsMnemonicValid(mnem) {
		return bip39.NewSeed(mnem, ""), nil
	}
	return nil, errors.New("invalid mnemonic")

}

func GetValidatorsFromMnemonic(mnemonic string, minAcc uint64, maxAcc uint64) ([]*Validator, error) {
	// adopted from eth2-val-tools: https://github.com/protolambda/eth2-val-tools
	var validators []*Validator

	seed, err := MnemonicToSeed(mnemonic)
	if err != nil {
		return nil, err
	}
	for idx := minAcc; idx < maxAcc; idx++ {
		valAccPath := fmt.Sprintf("m/12381/3600/%d/0/0", idx)
		withdrawalAccPath := fmt.Sprintf("m/12381/3600/%d/0", idx)
		validatorKey, err := util.PrivateKeyFromSeedAndPath(seed, valAccPath)
		if err != nil {
			return nil, fmt.Errorf("account %s cannot be derived, continuing to next account", valAccPath)
		}
		withdrawalKey, err := util.PrivateKeyFromSeedAndPath(seed, withdrawalAccPath)
		if err != nil {
			return nil, fmt.Errorf("withdrawal %s cannot be derived, continuing to next account", valAccPath)
		}
		var pubKey phase0.BLSPubKey
		copy(pubKey[:], validatorKey.PublicKey().Marshal()[:])

		validators = append(validators, &Validator{
			ValidatorIndex:     idx,
			WithdrawalKey:      withdrawalKey,
			ValidatorKey:       validatorKey,
			ValidatorPublicKey: pubKey,
		})
	}

	return validators, err
}
