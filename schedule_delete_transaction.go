package hedera

import (
	"github.com/hashgraph/hedera-sdk-go/v2/proto"

	"time"
)

type ScheduleDeleteTransaction struct {
	Transaction
	pb *proto.ScheduleDeleteTransactionBody
}

func NewScheduleDeleteTransaction() *ScheduleDeleteTransaction {
	pb := &proto.ScheduleDeleteTransactionBody{}

	transaction := ScheduleDeleteTransaction{
		pb:          pb,
		Transaction: newTransaction(),
	}
	transaction.SetMaxTransactionFee(NewHbar(5))

	return &transaction
}

func scheduleDeleteTransactionFromProtobuf(transaction Transaction, pb *proto.TransactionBody) ScheduleDeleteTransaction {
	return ScheduleDeleteTransaction{
		Transaction: transaction,
		pb:          pb.GetScheduleDelete(),
	}
}

func (transaction *ScheduleDeleteTransaction) SetScheduleID(scheduleID ScheduleID) *ScheduleDeleteTransaction {
	transaction.requireNotFrozen()
	transaction.pb.ScheduleID = scheduleID.toProtobuf()
	return transaction
}

func (transaction *ScheduleDeleteTransaction) GetScheduleID() ScheduleID {
	return scheduleIDFromProtobuf(transaction.pb.GetScheduleID())
}

//
// The following methods must be copy-pasted/overriden at the bottom of **every** _transaction.go file
// We override the embedded fluent setter methods to return the outer type
//

func scheduleDeleteTransaction_getMethod(request request, channel *channel) method {
	return method{
		transaction: channel.getSchedule().DeleteSchedule,
	}
}

func (transaction *ScheduleDeleteTransaction) IsFrozen() bool {
	return transaction.isFrozen()
}

// Sign uses the provided privateKey to sign the transaction.
func (transaction *ScheduleDeleteTransaction) Sign(
	privateKey PrivateKey,
) *ScheduleDeleteTransaction {
	return transaction.SignWith(privateKey.PublicKey(), privateKey.Sign)
}

func (transaction *ScheduleDeleteTransaction) SignWithOperator(
	client *Client,
) (*ScheduleDeleteTransaction, error) {
	// If the transaction is not signed by the operator, we need
	// to sign the transaction with the operator

	if client == nil {
		return nil, errNoClientProvided
	} else if client.operator == nil {
		return nil, errClientOperatorSigning
	}

	if !transaction.IsFrozen() {
		_, err := transaction.FreezeWith(client)
		if err != nil {
			return transaction, err
		}
	}
	return transaction.SignWith(client.operator.publicKey, client.operator.signer), nil
}

// SignWith executes the TransactionSigner and adds the resulting signature data to the Transaction's signature map
// with the publicKey as the map key.
func (transaction *ScheduleDeleteTransaction) SignWith(
	publicKey PublicKey,
	signer TransactionSigner,
) *ScheduleDeleteTransaction {
	if !transaction.IsFrozen() {
		_, _ = transaction.Freeze()
	} else {
		transaction.transactions = make([]*proto.Transaction, 0)
	}

	if transaction.keyAlreadySigned(publicKey) {
		return transaction
	}

	for index := 0; index < len(transaction.signedTransactions); index++ {
		signature := signer(transaction.signedTransactions[index].GetBodyBytes())

		transaction.signedTransactions[index].SigMap.SigPair = append(
			transaction.signedTransactions[index].SigMap.SigPair,
			publicKey.toSignaturePairProtobuf(signature),
		)
	}

	return transaction
}

// Execute executes the Transaction with the provided client
func (transaction *ScheduleDeleteTransaction) Execute(
	client *Client,
) (TransactionResponse, error) {
	if client == nil || client.operator == nil {
		return TransactionResponse{}, errNoClientProvided
	}

	if transaction.freezeError != nil {
		return TransactionResponse{}, transaction.freezeError
	}

	if !transaction.IsFrozen() {
		_, err := transaction.FreezeWith(client)
		if err != nil {
			return TransactionResponse{}, err
		}
	}

	transactionID := transaction.GetTransactionID()

	if !client.GetOperatorAccountID().isZero() && client.GetOperatorAccountID().equals(transactionID.AccountID) {
		transaction.SignWith(
			client.GetOperatorPublicKey(),
			client.operator.signer,
		)
	}

	resp, err := execute(
		client,
		request{
			transaction: &transaction.Transaction,
		},
		transaction_shouldRetry,
		transaction_makeRequest,
		transaction_advanceRequest,
		transaction_getNodeAccountID,
		scheduleDeleteTransaction_getMethod,
		transaction_mapResponseStatus,
		transaction_mapResponse,
	)

	if err != nil {
		return TransactionResponse{
			TransactionID: transaction.GetTransactionID(),
			NodeID:        resp.transaction.NodeID,
		}, err
	}

	hash, err := transaction.GetTransactionHash()

	return TransactionResponse{
		TransactionID: transaction.GetTransactionID(),
		NodeID:        resp.transaction.NodeID,
		Hash:          hash,
	}, nil
}

func (transaction *ScheduleDeleteTransaction) onFreeze(
	pbBody *proto.TransactionBody,
) bool {
	pbBody.Data = &proto.TransactionBody_ScheduleDelete{
		ScheduleDelete: transaction.pb,
	}

	return true
}

func (transaction *ScheduleDeleteTransaction) Freeze() (*ScheduleDeleteTransaction, error) {
	return transaction.FreezeWith(nil)
}

func (transaction *ScheduleDeleteTransaction) FreezeWith(client *Client) (*ScheduleDeleteTransaction, error) {
	transaction.initFee(client)
	if err := transaction.initTransactionID(client); err != nil {
		return transaction, err
	}

	if !transaction.onFreeze(transaction.pbBody) {
		return transaction, nil
	}

	return transaction, transaction_freezeWith(&transaction.Transaction, client)
}

func (transaction *ScheduleDeleteTransaction) GetMaxTransactionFee() Hbar {
	return transaction.Transaction.GetMaxTransactionFee()
}

// SetMaxTransactionFee sets the max transaction fee for this ScheduleDeleteTransaction.
func (transaction *ScheduleDeleteTransaction) SetMaxTransactionFee(fee Hbar) *ScheduleDeleteTransaction {
	transaction.requireNotFrozen()
	transaction.Transaction.SetMaxTransactionFee(fee)
	return transaction
}

func (transaction *ScheduleDeleteTransaction) GetTransactionMemo() string {
	return transaction.Transaction.GetTransactionMemo()
}

// SetTransactionMemo sets the memo for this ScheduleDeleteTransaction.
func (transaction *ScheduleDeleteTransaction) SetTransactionMemo(memo string) *ScheduleDeleteTransaction {
	transaction.requireNotFrozen()
	transaction.Transaction.SetTransactionMemo(memo)
	return transaction
}

func (transaction *ScheduleDeleteTransaction) GetTransactionValidDuration() time.Duration {
	return transaction.Transaction.GetTransactionValidDuration()
}

// SetTransactionValidDuration sets the valid duration for this ScheduleDeleteTransaction.
func (transaction *ScheduleDeleteTransaction) SetTransactionValidDuration(duration time.Duration) *ScheduleDeleteTransaction {
	transaction.requireNotFrozen()
	transaction.Transaction.SetTransactionValidDuration(duration)
	return transaction
}

func (transaction *ScheduleDeleteTransaction) GetTransactionID() TransactionID {
	return transaction.Transaction.GetTransactionID()
}

// SetTransactionID sets the TransactionID for this ScheduleDeleteTransaction.
func (transaction *ScheduleDeleteTransaction) SetTransactionID(transactionID TransactionID) *ScheduleDeleteTransaction {
	transaction.requireNotFrozen()

	transaction.Transaction.SetTransactionID(transactionID)
	return transaction
}

// SetNodeAccountID sets the node AccountID for this ScheduleDeleteTransaction.
func (transaction *ScheduleDeleteTransaction) SetNodeAccountIDs(nodeID []AccountID) *ScheduleDeleteTransaction {
	transaction.requireNotFrozen()
	transaction.Transaction.SetNodeAccountIDs(nodeID)
	return transaction
}

func (transaction *ScheduleDeleteTransaction) SetMaxRetry(count int) *ScheduleDeleteTransaction {
	transaction.Transaction.SetMaxRetry(count)
	return transaction
}
