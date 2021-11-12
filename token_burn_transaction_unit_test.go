//+build all unit

package hedera

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUnitTokenBurnTransactionValidate(t *testing.T) {
	client := ClientForTestnet()
	client.SetAutoValidateChecksums(true)
	tokenID, err := TokenIDFromString("0.0.123-rmkyk")
	assert.NoError(t, err)

	tokenBurn := NewTokenBurnTransaction().
		SetTokenID(tokenID)

	err = tokenBurn._ValidateNetworkOnIDs(client)
	assert.NoError(t, err)
}

func TestUnitTokenBurnTransactionValidateWrong(t *testing.T) {
	client := ClientForTestnet()
	client.SetAutoValidateChecksums(true)
	tokenID, err := TokenIDFromString("0.0.123-rmkykd")
	assert.NoError(t, err)

	tokenBurn := NewTokenBurnTransaction().
		SetTokenID(tokenID)

	err = tokenBurn._ValidateNetworkOnIDs(client)
	assert.Error(t, err)
	if err != nil {
		assert.Equal(t, "network mismatch or wrong checksum given, given checksum: rmkykd, correct checksum rmkyk, network: testnet", err.Error())
	}
}