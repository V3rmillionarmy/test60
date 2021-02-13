package hedera

import (
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestSerializeAccountCreateTransaction(t *testing.T) {
	date := time.Unix(1554158542, 0)

	key, err := Ed25519PrivateKeyFromString("302e020100300506032b6570042204203b054fade7a2b0869c6bd4a63b7017cbae7855d12acc357bea718e2c3e805962")
	assert.NoError(t, err)

	tx, err := NewAccountCreateTransaction().
		SetKey(key.PublicKey()).
		SetInitialBalance(HbarFromTinybar(450)).
		SetProxyAccountID(AccountID{Account: 1020}).
		SetReceiverSignatureRequired(true).
		SetNodeAccountID(AccountID{Account: 3}).
		SetTransactionID(TransactionID{
			AccountID:  AccountID{Account: 2},
			ValidStart: date,
		}).
		SetMaxTransactionFee(HbarFromTinybar(1e6)).
		Build(nil)

	assert.NoError(t, err)

	tx.Sign(key)

	assert.Equal(t, `bodyBytes:"\n\014\n\006\010\316\247\212\345\005\022\002\030\002\022\002\030\003\030\300\204=\"\002\010xZI\n\"\022\344\361\300\353L}\315\303\347\353\021p\263\010\212=\022\242\227\364\243\353\342\362\205\003\375g5F\355\216\020\302\003\032\003\030\374\0070\377\377\377\377\377\377\377\377\1778\377\377\377\377\377\377\377\377\177@\001J\005\010\320\310\341\003"sigMap:<sigPair:<pubKeyPrefix:"\344\361\300\353L}\315\303\347\353\021p\263\010\212=\022\242\227\364\243\353\342\362\205\003\375g5F\355\216"ed25519:"Y\361\353\220n\377\223\354\257\356\363\263\245;=\372L>\r\332?q\336\014\3713\253\222\031]\212\313\213\326v\343}\273\376\363\302\004\306u\221=x]&j\315:-\364\006l\nf\362\322Xd\220\013">>transactionID:<transactionValidStart:<seconds:1554158542>accountID:<accountNum:2>>nodeAccountID:<accountNum:3>transactionFee:1000000transactionValidDuration:<seconds:120>cryptoCreateAccount:<key:<ed25519:"\344\361\300\353L}\315\303\347\353\021p\263\010\212=\022\242\227\364\243\353\342\362\205\003\375g5F\355\216">initialBalance:450proxyAccountID:<accountNum:1020>sendRecordThreshold:9223372036854775807receiveRecordThreshold:9223372036854775807receiverSigRequired:trueautoRenewPeriod:<seconds:7890000>>`, strings.ReplaceAll(strings.ReplaceAll(tx.String(), " ", ""), "\n", ""))
}

func TestAccountCreateTransaction_Execute(t *testing.T) {
	client := newTestClient(t)

	newKey, err := GenerateEd25519PrivateKey()
	assert.NoError(t, err)

	newBalance := NewHbar(1)

	assert.Equal(t, HbarUnits.Hbar.numberOfTinybar(), newBalance.tinybar)

	txID, err := NewAccountCreateTransaction().
		SetKey(newKey.PublicKey()).
		SetMaxTransactionFee(NewHbar(2)).
		SetInitialBalance(newBalance).
		Execute(client)
	assert.NoError(t, err)

	receipt, err := txID.GetReceipt(client)
	assert.NoError(t, err)

	accountID := receipt.GetAccountID()
	assert.NoError(t, err)

	tx, err := NewAccountDeleteTransaction().
		SetDeleteAccountID(accountID).
		SetTransferAccountID(client.GetOperatorID()).
		SetMaxTransactionFee(NewHbar(1)).
		SetTransactionID(NewTransactionID(accountID)).
		Build(client)
	assert.NoError(t, err)

	txID, err = tx.
		Sign(newKey).
		Execute(client)
	assert.NoError(t, err)

	_, err = txID.GetReceipt(client)
	assert.NoError(t, err)
}

func TestAccountCreateTransactionAddSignature(t *testing.T) {
	client := newTestClient(t)

	newKey, err := GenerateEd25519PrivateKey()
	assert.NoError(t, err)

	resp, err := NewAccountCreateTransaction().
		SetKey(newKey.PublicKey()).
		Execute(client)
	assert.NoError(t, err)

	println(resp.String())

	receipt, err := resp.GetReceipt(client)
	assert.NoError(t, err)

	account := receipt.GetAccountID()

	tx, err := NewAccountDeleteTransaction().
		SetNodeAccountID(AccountID{Account: 3}).
		SetDeleteAccountID(account).
		SetTransferAccountID(client.GetOperatorID()).
		Build(client)
	assert.NoError(t, err)

	updateBytes, err := tx.MarshalBinary()
	assert.NoError(t, err)

	sig1, err := newKey.SignTransaction(&tx)
	assert.NoError(t, err)

	tx2, err := NewAccountDeleteTransaction().Build(client)
	assert.NoError(t, err)

	println(account.String())

	err = tx2.UnmarshalBinary(updateBytes)
	assert.NoError(t, err)

	resp, err = tx2.AddSignature(newKey.PublicKey(), sig1).Execute(client)
	assert.NoError(t, err)

	_, err = resp.GetReceipt(client)
	assert.NoError(t, err)
}
