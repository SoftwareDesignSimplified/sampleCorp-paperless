package db

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func createRandomPurchaseOrder(t *testing.T) PurchaseOrder {
	company := createRandomCompany(t)
	signatory := createRandomSignatory(t) // Assuming you have a function to create a random signatory
	arg := CreatePurchaseOrderParams{
		Attn:           "Attn",
		CompanyID:      int32(company.ID),
		Address:        "Address",
		SignatoryID:    int32(signatory.ID), // Include signatory_id in the arg variable
		QuotationID:    nil,
		SentOrReceived: nil,
	}

	purchaseOrder, err := testQueries.CreatePurchaseOrder(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, purchaseOrder)
	return purchaseOrder
}

func TestQueries_CreatePurchaseOrder(t *testing.T) {
	createRandomPurchaseOrder(t)
}

func TestQueries_ApprovePurchaseOrder(t *testing.T) {
	user := CreateRandomUser(t)
	id := user.ID
	d := int32(id)
	status := "APPROVED"
	purchaseOrder := createRandomPurchaseOrder(t)
	arg := ApprovePurchaseOrderParams{
		ID:         purchaseOrder.ID,
		ApprovedBy: &d,
		PoStatus:   &status,
	}
	purchaseOrder, err := testQueries.ApprovePurchaseOrder(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, purchaseOrder)
	require.Equal(t, *purchaseOrder.PoStatus, status)
}

func TestQueries_ListPurchaseOrders(t *testing.T) {
	for i := 0; i < 10; i++ {
		createRandomPurchaseOrder(t)
	}

	arg := ListPurchaseOrdersParams{
		Limit:  5,
		Offset: 5,
	}

	purchaseOrders, err := testQueries.ListPurchaseOrders(context.Background(), arg)
	require.NoError(t, err)
	require.Len(t, purchaseOrders, 5)

	for _, purchaseOrder := range purchaseOrders {
		require.NotEmpty(t, purchaseOrder)
	}
}

func TestQueries_GetPurchaseOrder(t *testing.T) {
	purchaseOrder1 := createRandomPurchaseOrder(t)
	purchaseOrder2, err := testQueries.GetPurchaseOrder(context.Background(), purchaseOrder1.ID)
	require.NoError(t, err)
	require.NotEmpty(t, purchaseOrder2)
	require.Equal(t, purchaseOrder1.ID, purchaseOrder2.ID)
	require.Equal(t, purchaseOrder1.PoNo, purchaseOrder2.PoNo)
	require.Equal(t, purchaseOrder1.Date, purchaseOrder2.Date)
	require.Equal(t, purchaseOrder1.Attn, purchaseOrder2.Attn)
	require.Equal(t, purchaseOrder1.CompanyID, purchaseOrder2.CompanyID)
	require.Equal(t, purchaseOrder1.Address, purchaseOrder2.Address)
	require.Equal(t, purchaseOrder1.SignatoryID, purchaseOrder2.SignatoryID)
	require.Equal(t, purchaseOrder1.QuotationID, purchaseOrder2.QuotationID)
	require.Equal(t, purchaseOrder1.PoStatus, purchaseOrder2.PoStatus)
	require.Equal(t, purchaseOrder1.ApprovedBy, purchaseOrder2.ApprovedBy)
	require.Equal(t, purchaseOrder1.SentOrReceived, purchaseOrder2.SentOrReceived)
	require.Equal(t, purchaseOrder1.PdfUrl, purchaseOrder2.PdfUrl)
}

func TestQueries_UpdatePurchaseOrder(t *testing.T) {
	purchaseOrder1 := createRandomPurchaseOrder(t)
	purchaseOrder2 := createRandomPurchaseOrder(t)

	arg := UpdatePurchaseOrderParams{
		ID:          purchaseOrder1.ID,
		Attn:        "Attn Updated",
		CompanyID:   purchaseOrder2.CompanyID,
		Address:     "Address Updated",
		SignatoryID: purchaseOrder2.SignatoryID,
	}

	purchaseOrder3, err := testQueries.UpdatePurchaseOrder(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, purchaseOrder3)
}
