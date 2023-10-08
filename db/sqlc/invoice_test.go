package db

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func createRandomInvoice(t *testing.T) Invoice {
	status := "sent"
	sig := createRandomSignatory(t)
	company := createRandomCompany(t)
	po := createRandomPurchaseOrder(t)
	bank := createRandomBankDetails(t)
	bankID := bank.ID
	arg := CreateInvoiceParams{
		Attn:                "Attn",
		PurchaseOrderNumber: po.PoNo,
		CompanyID:           int32(company.ID),
		AmountDue:           1000,
		BankDetails:         bankID,
		Site:                "Site",
		SignatoryID:         int32(sig.ID),
		SentOrReceived:      status,
	}
	invoice, err := testQueries.CreateInvoice(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, invoice)

	return invoice
}

func TestQueries_CreateInvoice(t *testing.T) {
	createRandomInvoice(t)
}

func TestQueries_GetInvoicesByPurchaseOrderNumber(t *testing.T) {
	invoice1 := createRandomInvoice(t)
	invoice2, err := testQueries.GetInvoicesByPurchaseOrderNumber(context.Background(), invoice1.PurchaseOrderNumber)
	require.NoError(t, err)
	require.NotEmpty(t, invoice2)
}
