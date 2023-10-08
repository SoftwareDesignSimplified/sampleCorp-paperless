package db

import (
	"context"
	"math/rand"
	"testing"

	"github.com/qwetu_petro/backend/utils"
	"github.com/stretchr/testify/require"
)

func createRandomPurchaseOrderItem(t *testing.T) PurchaseOrderItem {
	description := utils.RandomString(100)
	partNo := utils.RandomString(3)
	uom := utils.RandomString(1)
	purchaseOrder := createRandomPurchaseOrder(t)
	id := purchaseOrder.ID
	currency := utils.RandomString(3)
	qty := int32(utils.RandomInt())
	discount := rand.Float64() * 10000

	arg := CreatePurchaseOrderItemParams{
		Description:     description,
		PartNo:          partNo,
		Uom:             uom,
		Qty:             qty,
		ItemPrice:       100,                    // Generate a random number between 0 and 10000
		Discount:        discount,               // Generate a random number between 0 and 10000
		NetPrice:        rand.Float64() * 10000, // Generate a random number between 0 and 10000
		NetValue:        rand.Float64() * 10000, // Generate a random number between 0 and 10000
		Currency:        currency,
		PurchaseOrderID: id,
	}
	purchaseOrderItem, err := testQueries.CreatePurchaseOrderItem(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, purchaseOrderItem)
	return purchaseOrderItem
}

func TestQueries_CreatePurchaseOrderItem(t *testing.T) {
	createRandomPurchaseOrderItem(t)
}

func TestQueries_GetPurchaseOrderItemsByPurchaseOrderID(t *testing.T) {
	purchaseOrderItem1 := createRandomPurchaseOrderItem(t)
	purchaseOrderItem2, err := testQueries.GetPurchaseOrderItemsByPurchaseOrderID(context.Background(), purchaseOrderItem1.PurchaseOrderID)
	require.NoError(t, err)
	require.NotEmpty(t, purchaseOrderItem2)
}
