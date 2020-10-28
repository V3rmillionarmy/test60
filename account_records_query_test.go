package hedera

import (
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestAccountRecordQuery(t *testing.T) {
	var client *Client

	network := os.Getenv("HEDERA_NETWORK")

	if network == "previewnet" {
		client = ClientForPreviewnet()
	}

	client, err := ClientFromJsonFile(os.Getenv("CONFIG_FILE"))

	if err != nil {
		client = ClientForTestnet()
	}

	configOperatorID := os.Getenv("OPERATOR_ID")
	configOperatorKey := os.Getenv("OPERATOR_KEY")

	if configOperatorID != "" && configOperatorKey != "" {
		operatorAccountID, err := AccountIDFromString(configOperatorID)
		assert.NoError(t, err)

		operatorKey, err := PrivateKeyFromString(configOperatorKey)
		assert.NoError(t, err)

		client.SetOperator(operatorAccountID, operatorKey)
	}

	newKey, err := GeneratePrivateKey()
	assert.NoError(t, err)

	resp, err := NewAccountCreateTransaction().
		SetKey(newKey.PublicKey()).
		SetMaxTransactionFee(NewHbar(2)).
		SetInitialBalance(NewHbar(1)).
		Execute(client)
	assert.NoError(t, err)

	receipt, err := resp.GetReceipt(client)
	assert.NoError(t, err)

	account := *receipt.AccountID

	nodeIDs := make([]AccountID, 1)
	nodeIDs[0] = resp.NodeID

	_, err = NewCryptoTransferTransaction().
		SetNodeAccountIDs(nodeIDs).
		AddRecipient(account, NewHbar(1)).
		AddSender(client.GetOperatorID(), NewHbar(1)).
		Execute(client)
	assert.NoError(t, err)

	recordsQuery, err := NewAccountRecordsQuery().
		SetNodeAccountIDs(nodeIDs).
		SetAccountID(client.GetOperatorID()).
		SetMaxQueryPayment(NewHbar(1)).
		Execute(client)
	assert.NoError(t, err)

	assert.Equal(t, 0, len(recordsQuery))

	//assert.Equal(t, `cryptoGetAccountRecords:{header:{payment:{bodyBytes:"\n\x0e\n\x08\x08\xdc\xc9\x07\x10۟\t\x12\x02\x18\x03\x12\x02\x18\x03\x18\x80\xc2\xd7/\"\x02\x08xr\x14\n\x12\n\x07\n\x02\x18\x02\x10\xc7\x01\n\x07\n\x02\x18\x03\x10\xc8\x01"sigMap:{sigPair:{pubKeyPrefix:"\xe4\xf1\xc0\xebL}\xcd\xc3\xe7\xeb\x11p\xb3\x08\x8a=\x12\xa2\x97\xf4\xa3\xeb\xe2\xf2\x85\x03\xfdg5F\xed\x8e"ed25519:"\x12&5\x96\xfb\xb4\x1c]P\xbb%\xecP\x9bk͙\x0b߼\xac)\xa6+\xd2<\x97+\xbb\x8c\x8af\xcb\xdai\x17T4{\xf7\xf3UYn\n\x8f\xabep\x04\xf6\x83\x0f\xbaFUP\xa3\xd1/\x1d\x9d\x1a\x0b"}}}}accountID:{accountNum:3}}`, strings.ReplaceAll(query.QueryBuilder.pb.String(), " ", ""))
}
