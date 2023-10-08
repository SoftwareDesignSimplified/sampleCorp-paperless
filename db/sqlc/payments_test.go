package db

import (
	"context"
	"github.com/jackc/pgtype"
	"math/rand"
	"testing"
	"time"

	"github.com/qwetu_petro/backend/utils"

	"github.com/stretchr/testify/require"
)

func createRandomPaymentRequest(t *testing.T) PaymentRequest {
	user := CreateRandomUser(t)

	args := CreatePaymentRequestParams{
		EmployeeID:  int32(user.ID),
		Amount:      600.0,
		Status:      "PENDING",
		Description: utils.RandomString(100),
	}

	paymentRequest, err := testQueries.CreatePaymentRequest(context.Background(), args)
	require.NoError(t, err)
	require.NotEmpty(t, paymentRequest)

	return paymentRequest

}

func createRandomPettyCash(t *testing.T) PettyCash {
	user := CreateRandomUser(t)

	randomFloat := rand.Float64() * 1000
	var amount pgtype.Numeric
	err := amount.Set(randomFloat)
	if err != nil {
		panic(err)
	}

	args := CreatePettyCashParams{
		EmployeeID:      int32(user.ID),
		Amount:          amount,
		Description:     utils.RandomString(100),
		TransactionDate: time.Now(),
		Folio:           "BANK",
	}

	pettyCash, err := testQueries.CreatePettyCash(context.Background(), args)
	require.NoError(t, err)
	require.NotEmpty(t, pettyCash)

	return pettyCash

}

func TestQueries_CreatePettyCash(t *testing.T) {
	createRandomPettyCash(t)
}

func TestQueries_CreatePaymentRequest(t *testing.T) {
	createRandomPaymentRequest(t)
}

func TestQueries_UpdatePaymentRequest(t *testing.T) {
	paymentRequest := createRandomPaymentRequest(t)

	args := UpdatePaymentRequestParams{
		RequestID:   paymentRequest.RequestID,
		EmployeeID:  paymentRequest.EmployeeID,
		Amount:      paymentRequest.Amount,
		Description: paymentRequest.Description,
		AdminID:     paymentRequest.AdminID,
		InvoiceID:   paymentRequest.InvoiceID,
	}

	paymentRequest2, err := testQueries.UpdatePaymentRequest(context.Background(), args)
	require.NoError(t, err)
	require.NotEmpty(t, paymentRequest2)
	require.Equal(t, paymentRequest.RequestID, paymentRequest2.RequestID)

}

func TestQueries_ApprovePaymentRequest(t *testing.T) {
	date := time.Now()
	admin := CreateRandomUser(t)
	id := int32(admin.ID)
	paymentRequest := createRandomPaymentRequest(t)
	args := ApprovePaymentRequestParams{
		Status:       "APPROVED",
		ApprovalDate: &date,
		RequestID:    paymentRequest.RequestID,
		AdminID:      &id,
	}

	paymentRequest2, err := testQueries.ApprovePaymentRequest(context.Background(), args)
	require.NoError(t, err)
	require.NotEmpty(t, paymentRequest2)
	require.Equal(t, paymentRequest.RequestID, paymentRequest2.RequestID)
	require.Equal(t, args.Status, paymentRequest2.Status)

}

func TestQueries_DeletePaymentRequest(t *testing.T) {
	paymentRequest := createRandomPaymentRequest(t)
	err := testQueries.DeletePaymentRequest(context.Background(), paymentRequest.RequestID)
	require.NoError(t, err)
}

func TestQueries_ListEmployeePaymentRequests(t *testing.T) {
	var lastPaymentRequest PaymentRequest
	for i := 0; i < 10; i++ {
		lastPaymentRequest = createRandomPaymentRequest(t)
	}

	arg := ListEmployeePaymentRequestsParams{
		EmployeeID: lastPaymentRequest.EmployeeID,
		Limit:      5,
		Offset:     0,
	}

	paymentRequests, err := testQueries.ListEmployeePaymentRequests(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, paymentRequests)

	for _, paymentRequest := range paymentRequests {
		require.NotEmpty(t, paymentRequest)
		require.Equal(t, arg.EmployeeID, paymentRequest.EmployeeID)
	}
}

func TestQueries_ListEmployeePettyCash(t *testing.T) {

	var lastPettyCash PettyCash
	for i := 0; i < 10; i++ {
		lastPettyCash = createRandomPettyCash(t)
	}

	arg := ListEmployeePettyCashParams{
		EmployeeID: lastPettyCash.EmployeeID,
		Limit:      5,
		Offset:     0,
	}

	pettyCash, err := testQueries.ListEmployeePettyCash(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, pettyCash)

	for _, petty := range pettyCash {
		require.NotEmpty(t, petty)
		require.Equal(t, arg.EmployeeID, petty.EmployeeID)
	}
}

func TestQueries_DeletePettyCash(t *testing.T) {
	pettyCash := createRandomPettyCash(t)
	err := testQueries.DeletePettyCash(context.Background(), pettyCash.TransactionID)
	require.NoError(t, err)
}

func TestQueries_UpdatePettyCash(t *testing.T) {
	pettyCash := createRandomPettyCash(t)

	args := UpdatePettyCashParams{
		TransactionID: pettyCash.TransactionID,
		Amount:        pettyCash.Amount,
		Description:   "New Description",
	}

	pettyCash2, err := testQueries.UpdatePettyCash(context.Background(), args)
	require.NoError(t, err)
	require.NotEmpty(t, pettyCash2)
	require.Equal(t, pettyCash.TransactionID, pettyCash2.TransactionID)
}
