package hedera

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
)

func TestSerializeFileInfoQuery(t *testing.T) {
	query := NewFileInfoQuery().
		SetFileID(FileID{File: 3}).
		Query

	assert.Equal(t, `fileGetInfo:{header:{}fileID:{fileNum:3}}`, strings.ReplaceAll(query.pb.String(), " ", ""))
}

func Test_FileInfo_Transaction(t *testing.T) {
	client := newTestClient(t)

	resp, err := NewFileCreateTransaction().
		SetKeys(client.GetOperatorPublicKey()).
		SetContents([]byte("Hello, World")).
		SetTransactionMemo("go sdk e2e tests").
		Execute(client)

	assert.NoError(t, err)

	receipt, err := resp.GetReceipt(client)
	assert.NoError(t, err)

	fileID := receipt.FileID
	assert.NotNil(t, fileID)

	info, err := NewFileInfoQuery().
		SetFileID(*fileID).
		SetNodeAccountIDs([]AccountID{resp.NodeID}).
		SetMaxQueryPayment(NewHbar(1)).
		SetQueryPayment(HbarFromTinybar(25)).
		Execute(client)
	assert.NoError(t, err)

	assert.Equal(t, *fileID, info.FileID)
	assert.Equal(t, info.Size, int64(12))
	assert.False(t, info.IsDeleted)
	assert.NotNil(t, info.Keys)

	resp, err = NewFileDeleteTransaction().
		SetFileID(*fileID).
		SetNodeAccountIDs([]AccountID{resp.NodeID}).
		Execute(client)
	assert.NoError(t, err)

	_, err = resp.GetReceipt(client)
	assert.NoError(t, err)
}

func Test_FileInfoCost_Transaction(t *testing.T) {
	client := newTestClient(t)

	resp, err := NewFileCreateTransaction().
		SetKeys(client.GetOperatorPublicKey()).
		SetContents([]byte("Hello, World")).
		SetTransactionMemo("go sdk e2e tests").
		Execute(client)

	assert.NoError(t, err)

	receipt, err := resp.GetReceipt(client)
	assert.NoError(t, err)

	fileID := receipt.FileID
	assert.NotNil(t, fileID)

	fileInfo := NewFileInfoQuery().
		SetFileID(*fileID).
		SetMaxQueryPayment(NewHbar(1)).
		SetNodeAccountIDs([]AccountID{resp.NodeID})

	cost, err := fileInfo.GetCost(client)
	assert.NoError(t, err)

	info, err := fileInfo.SetQueryPayment(cost).Execute(client)
	assert.NoError(t, err)

	assert.Equal(t, *fileID, info.FileID)
	assert.Equal(t, info.Size, int64(12))
	assert.False(t, info.IsDeleted)

	resp, err = NewFileDeleteTransaction().
		SetFileID(*fileID).
		SetNodeAccountIDs([]AccountID{resp.NodeID}).
		Execute(client)
	assert.NoError(t, err)

	_, err = resp.GetReceipt(client)
	assert.NoError(t, err)
}

func Test_FileInfoCost_BigMax_Transaction(t *testing.T) {
	client := newTestClient(t)

	resp, err := NewFileCreateTransaction().
		SetKeys(client.GetOperatorPublicKey()).
		SetContents([]byte("Hello, World")).
		SetTransactionMemo("go sdk e2e tests").
		Execute(client)

	assert.NoError(t, err)

	receipt, err := resp.GetReceipt(client)
	assert.NoError(t, err)

	fileID := receipt.FileID
	assert.NotNil(t, fileID)

	fileInfo := NewFileInfoQuery().
		SetFileID(*fileID).
		SetMaxQueryPayment(NewHbar(10000)).
		SetNodeAccountIDs([]AccountID{resp.NodeID})

	_, err = fileInfo.GetCost(client)
	assert.NoError(t, err)

	info, err := fileInfo.Execute(client)
	assert.NoError(t, err)

	assert.Equal(t, *fileID, info.FileID)
	assert.Equal(t, info.Size, int64(12))
	assert.False(t, info.IsDeleted)

	resp, err = NewFileDeleteTransaction().
		SetFileID(*fileID).
		SetNodeAccountIDs([]AccountID{resp.NodeID}).
		Execute(client)
	assert.NoError(t, err)

	_, err = resp.GetReceipt(client)
	assert.NoError(t, err)
}

func Test_FileInfoCost_SmallMax_Transaction(t *testing.T) {
	client := newTestClient(t)

	resp, err := NewFileCreateTransaction().
		SetKeys(client.GetOperatorPublicKey()).
		SetContents([]byte("Hello, World")).
		SetTransactionMemo("go sdk e2e tests").
		Execute(client)

	assert.NoError(t, err)

	receipt, err := resp.GetReceipt(client)
	assert.NoError(t, err)

	fileID := receipt.FileID
	assert.NotNil(t, fileID)

	fileInfo := NewFileInfoQuery().
		SetFileID(*fileID).
		SetMaxQueryPayment(HbarFromTinybar(1)).
		SetNodeAccountIDs([]AccountID{resp.NodeID})

	cost, err := fileInfo.GetCost(client)
	assert.NoError(t, err)

	_, err = fileInfo.Execute(client)
	if err != nil {
		assert.Equal(t, fmt.Sprintf("cost of FileInfoQuery ("+cost.String()+") without explicit payment is greater than the max query payment of 1 tħ"), err.Error())
	}

	resp, err = NewFileDeleteTransaction().
		SetFileID(*fileID).
		SetNodeAccountIDs([]AccountID{resp.NodeID}).
		Execute(client)
	assert.NoError(t, err)

	_, err = resp.GetReceipt(client)
	assert.NoError(t, err)
}

func Test_FileInfoCost_InsufficientFee_Transaction(t *testing.T) {
	client := newTestClient(t)

	resp, err := NewFileCreateTransaction().
		SetKeys(client.GetOperatorPublicKey()).
		SetContents([]byte("Hello, World")).
		SetTransactionMemo("go sdk e2e tests").
		Execute(client)

	assert.NoError(t, err)

	receipt, err := resp.GetReceipt(client)
	assert.NoError(t, err)

	fileID := receipt.FileID
	assert.NotNil(t, fileID)

	fileInfo := NewFileInfoQuery().
		SetFileID(*fileID).
		SetMaxQueryPayment(NewHbar(1)).
		SetNodeAccountIDs([]AccountID{resp.NodeID})

	_, err = fileInfo.GetCost(client)
	assert.NoError(t, err)

	_, err = fileInfo.SetQueryPayment(HbarFromTinybar(1)).Execute(client)
	if err != nil {
		assert.Equal(t, fmt.Sprintf("exceptional precheck status INSUFFICIENT_TX_FEE"), err.Error())
	}

	resp, err = NewFileDeleteTransaction().
		SetFileID(*fileID).
		SetNodeAccountIDs([]AccountID{resp.NodeID}).
		Execute(client)
	assert.NoError(t, err)

	_, err = resp.GetReceipt(client)
	assert.NoError(t, err)
}

func Test_FileInfoQuery_NoFileID(t *testing.T) {
	client := newTestClient(t)

	_, err := NewFileInfoQuery().
		SetQueryPayment(NewHbar(1)).
		Execute(client)
	assert.Error(t, err)
	if err != nil {
		assert.Equal(t, fmt.Sprintf("exceptional precheck status INVALID_FILE_ID"), err.Error())
	}
}
