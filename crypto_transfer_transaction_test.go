package hedera

import (
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
)

func TestSerializeCryptoTransferTransaction(t *testing.T) {
	tx, err := newMockTransaction()
	assert.NoError(t, err)

	assert.Equal(t, `bodyBytes:"\n\016\n\010\010\334\311\007\020\333\237\t\022\002\030\003\022\002\030\003\030\200\302\327/\"\002\010xr\024\n\022\n\007\n\002\030\002\020\307\001\n\007\n\002\030\003\020\310\001"sigMap:<sigPair:<pubKeyPrefix:"\344\361\300\353L}\315\303\347\353\021p\263\010\212=\022\242\227\364\243\353\342\362\205\003\375g5F\355\216"ed25519:"\022&5\226\373\264\034]P\273%\354P\233k\315\231\013\337\274\254)\246+\322<\227+\273\214\212f\313\332i\027T4{\367\363UYn\n\217\253ep\004\366\203\017\272FUP\243\321/\035\235\032\013">>transactionID:<transactionValidStart:<seconds:124124nanos:151515>accountID:<accountNum:3>>nodeAccountID:<accountNum:3>transactionFee:100000000transactionValidDuration:<seconds:120>cryptoTransfer:<transfers:<accountAmounts:<accountID:<accountNum:2>amount:-100>accountAmounts:<accountID:<accountNum:3>amount:100>>>`, strings.ReplaceAll(strings.ReplaceAll(tx.String(), " ", ""), "\n", ""))
}

func TestCryptoTransferTransaction_Execute(t *testing.T) {
	client := newTestClient(t)

	txID, err := NewCryptoTransferTransaction().
		AddSender(client.GetOperatorID(), NewHbar(1)).
		AddRecipient(AccountID{Account: 3}, NewHbar(1)).
		SetMaxTransactionFee(NewHbar(1)).
		Execute(client)
	assert.NoError(t, err)

	_, err = txID.GetReceipt(client)
	assert.NoError(t, err)
}
