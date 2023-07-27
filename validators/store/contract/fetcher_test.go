package contract

import (
	"errors"
	"fmt"
	"math/big"
	"testing"

	testHelper "tynmo/helper/tests"
	"tynmo/helper/staking"
	"tynmo/state"
	"tynmo/types"
	"tynmo/validators"
	"github.com/stretchr/testify/assert"
)

func TestFetchValidators(t *testing.T) {
	t.Parallel()

	// only check error handling because of the duplicated tests below
	fakeValidatorType := validators.ValidatorType("fake")
	res, err := FetchValidators(
		fakeValidatorType,
		nil,
		types.ZeroAddress,
	)

	assert.Nil(t, res)
	assert.ErrorContains(t, err, fmt.Sprintf("unsupported validator type: %s", fakeValidatorType))
}

func TestFetchECDSAValidators(t *testing.T) {
	t.Parallel()

	val := staking.DefaultStakedBalance
	bigDefaultStakedBalance, err := types.ParseUint256orHex(&val)
	if err != nil {
		bigDefaultStakedBalance = big.NewInt(0)
	}
	bigDefaultStakedBalance2 := big.NewInt(0).Set(bigDefaultStakedBalance)
	bigDefaultStakedBalance2.Mul(bigDefaultStakedBalance2, big.NewInt(2))

	var (
		ecdsaValidators = validators.NewECDSAValidatorSet(
			validators.NewECDSAValidator(addr1, bigDefaultStakedBalance),
			validators.NewECDSAValidator(addr2, bigDefaultStakedBalance2),
		)
	)

	tests := []struct {
		name        string
		transition  *state.Transition
		from        types.Address
		expectedRes validators.Validators
		expectedErr error
	}{
		{
			name: "should return error if QueryValidators failed",
			transition: newTestTransition(
				t,
			),
			from:        types.ZeroAddress,
			expectedRes: nil,
			expectedErr: errors.New("empty input"),
		},
		{
			name: "should return ECDSA Validators",
			transition: newTestTransitionWithPredeployedStakingContract(
				t,
				ecdsaValidators,
			),
			from:        types.ZeroAddress,
			expectedRes: ecdsaValidators,
			expectedErr: nil,
		},
	}

	for _, test := range tests {
		test := test

		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			res, err := FetchValidators(
				validators.ECDSAValidatorType,
				test.transition,
				test.from,
			)

			assert.Equal(t, test.expectedRes, res)
			testHelper.AssertErrorMessageContains(t, test.expectedErr, err)
		})
	}
}

func TestFetchBLSValidators(t *testing.T) {
	t.Parallel()

	val := staking.DefaultStakedBalance
	bigDefaultStakedBalance, err := types.ParseUint256orHex(&val)
	if err != nil {
		bigDefaultStakedBalance = big.NewInt(0)
	}

	var (
		blsValidators = validators.NewBLSValidatorSet(
			validators.NewBLSValidator(addr1, testBLSPubKey1, bigDefaultStakedBalance),
			validators.NewBLSValidator(addr2, []byte{}), // validator 2 has not set BLS Public Key
		)
	)

	tests := []struct {
		name        string
		transition  *state.Transition
		from        types.Address
		expectedRes validators.Validators
		expectedErr error
	}{
		{
			name: "should return error if QueryValidators failed",
			transition: newTestTransition(
				t,
			),
			from:        types.ZeroAddress,
			expectedRes: nil,
			expectedErr: errors.New("empty input"),
		},
		{
			name: "should return BLS Validators",
			transition: newTestTransitionWithPredeployedStakingContract(
				t,
				blsValidators,
			),
			from: types.ZeroAddress,
			expectedRes: validators.NewBLSValidatorSet(
				validators.NewBLSValidator(addr1, testBLSPubKey1, bigDefaultStakedBalance),
			),
			expectedErr: nil,
		},
	}

	for _, test := range tests {
		test := test

		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			res, err := FetchValidators(
				validators.BLSValidatorType,
				test.transition,
				test.from,
			)

			assert.Equal(t, test.expectedRes, res)
			testHelper.AssertErrorMessageContains(t, test.expectedErr, err)
		})
	}
}
