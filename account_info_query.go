package hedera

import (
	"github.com/hashgraph/hedera-sdk-go/v2/proto"
)

type AccountInfoQuery struct {
	Query
	pb *proto.CryptoGetInfoQuery
}

func NewAccountInfoQuery() *AccountInfoQuery {
	header := proto.QueryHeader{}
	query := newQuery(true, &header)
	pb := proto.CryptoGetInfoQuery{Header: &header}
	query.pb.Query = &proto.Query_CryptoGetInfo{
		CryptoGetInfo: &pb,
	}

	return &AccountInfoQuery{
		Query: query,
		pb:    &pb,
	}
}

func accountInfoQuery_mapResponseStatus(_ request, response response) Status {
	return Status(response.query.GetCryptoGetInfo().Header.NodeTransactionPrecheckCode)
}

func accountInfoQuery_getMethod(_ request, channel *channel) method {
	return method{
		query: channel.getCrypto().GetAccountInfo,
	}
}

// SetAccountID sets the AccountID for this AccountInfoQuery.
func (query *AccountInfoQuery) SetAccountID(accountID AccountID) *AccountInfoQuery {
	query.pb.AccountID = accountID.toProtobuf()
	return query
}

func (query *AccountInfoQuery) GetCost(client *Client) (Hbar, error) {
	if client == nil || client.operator == nil {
		return Hbar{}, errNoClientProvided
	}

	paymentTransaction, err := query_makePaymentTransaction(TransactionID{}, AccountID{}, client.operator, Hbar{})
	if err != nil {
		return Hbar{}, err
	}

	query.pbHeader.Payment = paymentTransaction
	query.pbHeader.ResponseType = proto.ResponseType_COST_ANSWER
	query.nodeIDs = client.network.getNodeAccountIDsForExecute()

	resp, err := execute(
		client,
		request{
			query: &query.Query,
		},
		query_shouldRetry,
		costQuery_makeRequest,
		costQuery_advanceRequest,
		costQuery_getNodeAccountID,
		accountInfoQuery_getMethod,
		accountInfoQuery_mapResponseStatus,
		query_mapResponse,
	)

	if err != nil {
		return Hbar{}, err
	}

	cost := int64(resp.query.GetCryptoGetInfo().Header.Cost)
	if cost < 25 {
		return HbarFromTinybar(25), nil
	} else {
		return HbarFromTinybar(cost), nil
	}
}

func (query *AccountInfoQuery) GetAccountID() AccountID {
	if query.pb.AccountID != nil {
		return AccountID{}
	} else {
		return accountIDFromProtobuf(query.pb.AccountID)
	}
}

// SetNodeAccountIDs sets the node AccountID for this AccountInfoQuery.
func (query *AccountInfoQuery) SetNodeAccountIDs(accountID []AccountID) *AccountInfoQuery {
	query.Query.SetNodeAccountIDs(accountID)
	return query
}

//SetQueryPayment sets the Hbar payment to pay the node a fee for handling this query
func (query *AccountInfoQuery) SetQueryPayment(queryPayment Hbar) *AccountInfoQuery {
	query.queryPayment = queryPayment
	return query
}

//SetMaxQueryPayment sets the maximum payment allowable for this query.
func (query *AccountInfoQuery) SetMaxQueryPayment(queryMaxPayment Hbar) *AccountInfoQuery {
	query.maxQueryPayment = queryMaxPayment
	return query
}

func (query *AccountInfoQuery) SetMaxRetry(count int) *AccountInfoQuery {
	query.Query.SetMaxRetry(count)
	return query
}

func (query *AccountInfoQuery) Execute(client *Client) (AccountInfo, error) {
	if client == nil || client.operator == nil {
		return AccountInfo{}, errNoClientProvided
	}

	if len(query.Query.GetNodeAccountIDs()) == 0 {
		query.SetNodeAccountIDs(client.network.getNodeAccountIDsForExecute())
	}

	query.paymentTransactionID = TransactionIDGenerate(client.operator.accountID)

	var cost Hbar
	if query.queryPayment.tinybar != 0 {
		cost = query.queryPayment
	} else {
		if query.maxQueryPayment.tinybar == 0 {
			cost = client.maxQueryPayment
		} else {
			cost = query.maxQueryPayment
		}

		actualCost, err := query.GetCost(client)
		if err != nil {
			return AccountInfo{}, err
		}

		if cost.tinybar < actualCost.tinybar {
			return AccountInfo{}, ErrMaxQueryPaymentExceeded{
				QueryCost:       actualCost,
				MaxQueryPayment: cost,
				query:           "AccountInfoQuery",
			}
		}

		cost = actualCost
	}

	err := query_generatePayments(&query.Query, client, cost)
	if err != nil {
		return AccountInfo{}, err
	}

	resp, err := execute(
		client,
		request{
			query: &query.Query,
		},
		query_shouldRetry,
		query_makeRequest,
		query_advanceRequest,
		query_getNodeAccountID,
		accountInfoQuery_getMethod,
		accountInfoQuery_mapResponseStatus,
		query_mapResponse,
	)

	if err != nil {
		return AccountInfo{}, err
	}

	return accountInfoFromProtobuf(resp.query.GetCryptoGetInfo().AccountInfo)
}
