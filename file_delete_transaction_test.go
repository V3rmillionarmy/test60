package hedera

import (
	"github.com/stretchr/testify/assert"
	"os"

	"testing"
)

// func TestSerializeFileDeleteTransaction(t *testing.T) {
// 	mockClient, err := newMockClient()
// 	assert.NoError(t, err)

// 	privateKey, err := PrivateKeyFromString(mockPrivateKey)
// 	assert.NoError(t, err)

// 	tx, err := NewFileDeleteTransaction().
// 		SetFileID(FileID{File: 5}).
// 		SetMaxTransactionFee(HbarFromTinybar(1e6)).
// 		SetTransactionID(testTransactionID).
// 		FreezeWith(mockClient)

// 	assert.NoError(t, err)

// 	tx.Sign(privateKey)

// 	assert.Equal(t, `bodyBytes:"\n\016\n\010\010\334\311\007\020\333\237\t\022\002\030\003\022\002\030\003\030\300\204=\"\002\010x\222\001\004\022\002\030\005"sigMap:<sigPair:<pubKeyPrefix:"\344\361\300\353L}\315\303\347\353\021p\263\010\212=\022\242\227\364\243\353\342\362\205\003\375g5F\355\216"ed25519:"cX\335@\024\365\365\3065\211NT.\355\245\224\364\230@\301\221\343\\T\343H\374\003\261W\252a\272\3401-)\251?N\204\305C\034\301\375\306\327K7a`+"`"+`r\262]\247\231I\332*:\2432\010">>transactionID:<transactionValidStart:<seconds:124124nanos:151515>accountID:<accountNum:3>>nodeAccountID:<accountNum:3>transactionFee:1000000transactionValidDuration:<seconds:120>fileDelete:<fileID:<fileNum:5>>`, strings.ReplaceAll(strings.ReplaceAll(tx.String(), " ", ""), "\n", ""))
// }

func TestFileDeleteTransaction_Execute(t *testing.T) {
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

	client.SetMaxTransactionFee(NewHbar(2))

	resp, err := NewFileCreateTransaction().
		SetKeys(client.GetOperatorKey()).
		SetContents([]byte("Hello, World")).
		SetTransactionMemo("go sdk e2e tests").
		Execute(client)

	assert.NoError(t, err)

	receipt, err := resp.GetReceipt(client)
	assert.NoError(t, err)

	fileID := *receipt.FileID
	assert.NotNil(t, fileID)

	nodeIDs := make([]AccountID, 1)
	nodeIDs[0] = resp.NodeID

	resp, err = NewFileDeleteTransaction().
		SetFileID(fileID).
		SetNodeAccountIDs(nodeIDs).
		Execute(client)
	assert.NoError(t, err)

	_, err = resp.GetReceipt(client)
	assert.NoError(t, err)
}
