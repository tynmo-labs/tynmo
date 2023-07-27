package validators

import (
	"encoding/json"
	"fmt"
	"math/big"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestECDSAValidatorsMarshalJSON(t *testing.T) {
	t.Parallel()

	validators := &Set{
		ValidatorType: ECDSAValidatorType,
		Validators: []Validator{
			&ECDSAValidator{addr1, big.NewInt(0), big.NewInt(0)},
			&ECDSAValidator{addr2, big.NewInt(0), big.NewInt(0)},
		},
	}

	res, err := json.Marshal(validators)

	assert.NoError(t, err)

	assert.JSONEq(
		t,
		fmt.Sprintf(
			`[
				{
					"Address": "%s",
					"Stake": 0,
					"Balance": 0
				},
				{
					"Address": "%s",
					"Stake": 0,
					"Balance": 0
				}
			]`,
			addr1.String(),
			addr2.String(),
		),
		string(res),
	)
}

func TestECDSAValidatorsUnmarshalJSON(t *testing.T) {
	t.Parallel()

	inputStr := fmt.Sprintf(
		`[
			{
				"Address": "%s",
				"Stake": 0,
				"Balance": 0
			},
			{
				"Address": "%s",
				"Stake": 0,
				"Balance": 0
			}
		]`,
		addr1.String(),
		addr2.String(),
	)

	validators := NewECDSAValidatorSet()

	assert.NoError(
		t,
		json.Unmarshal([]byte(inputStr), validators),
	)

	assert.Equal(
		t,
		&Set{
			ValidatorType: ECDSAValidatorType,
			Validators: []Validator{
				&ECDSAValidator{addr1, big.NewInt(0), big.NewInt(0)},
				&ECDSAValidator{addr2, big.NewInt(0), big.NewInt(0)},
			},
		},
		validators,
	)
}

func TestBLSValidatorsMarshalJSON(t *testing.T) {
	t.Parallel()

	validators := &Set{
		ValidatorType: BLSValidatorType,
		Validators: []Validator{
			&BLSValidator{addr1, testBLSPubKey1, big.NewInt(0), big.NewInt(0)},
			&BLSValidator{addr2, testBLSPubKey2, big.NewInt(0), big.NewInt(0)},
		},
	}

	res, err := json.Marshal(validators)

	assert.NoError(t, err)

	assert.JSONEq(
		t,
		fmt.Sprintf(
			`[
				{
					"Address": "%s",
					"BLSPublicKey": "%s",
					"Stake": 0,
					"Balance": 0
				},
				{
					"Address": "%s",
					"BLSPublicKey": "%s",
					"Stake": 0,
					"Balance": 0
				}
			]`,
			addr1,
			testBLSPubKey1,
			addr2,
			testBLSPubKey2,
		),
		string(res),
	)
}

func TestBLSValidatorsUnmarshalJSON(t *testing.T) {
	t.Parallel()

	inputStr := fmt.Sprintf(
		`[
			{
				"Address": "%s",
				"BLSPublicKey": "%s",
				"Stake": 0,
				"Balance": 0
			},
			{
				"Address": "%s",
				"BLSPublicKey": "%s",
				"Stake": 0,
				"Balance": 0
			}
		]`,
		addr1,
		testBLSPubKey1,
		addr2,
		testBLSPubKey2,
	)

	validators := NewBLSValidatorSet()

	assert.NoError(
		t,
		json.Unmarshal([]byte(inputStr), validators),
	)

	assert.Equal(
		t,
		&Set{
			ValidatorType: BLSValidatorType,
			Validators: []Validator{
				&BLSValidator{addr1, testBLSPubKey1, big.NewInt(0), big.NewInt(0)},
				&BLSValidator{addr2, testBLSPubKey2, big.NewInt(0), big.NewInt(0)},
			},
		},
		validators,
	)
}
