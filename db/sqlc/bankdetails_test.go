package db

import (
	"context"
	"github.com/qwetu_petro/backend/utils"
	"testing"

	"github.com/stretchr/testify/require"
)

func createRandomBankDetails(t *testing.T) BankDetail {
	company := createRandomCompany(t)
	arg := CreateBankDetailsParams{
		BankName:      utils.RandomString(10),
		AccountName:   utils.RandomString(10),
		AccountNumber: utils.RandomString(10),
		Branch:        utils.RandomString(10),
		SwiftCode:     utils.RandomString(4),
		Address:       utils.RandomString(10),
		Country:       utils.RandomString(10),
		Currency:      utils.RandomString(10),
		AccountType:   utils.RandomString(10),
		CompanyID:     int32(company.ID),
	}
	bankDetails, err := testQueries.CreateBankDetails(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, bankDetails)
	return bankDetails
}

func TestQueries_CreateBankDetails(t *testing.T) {
	createRandomBankDetails(t)
}

func TestQueries_GetBankDetailsByAccountNumber(t *testing.T) {
	bankDetails1 := createRandomBankDetails(t)
	bankDetails2, err := testQueries.GetBankDetailsByAccountNumber(context.Background(), bankDetails1.AccountNumber)
	require.NoError(t, err)
	require.NotEmpty(t, bankDetails2)
}

func TestQueries_ListBankDetailsByBankName(t *testing.T) {
	var lastBankDetail BankDetail
	for i := 0; i < 10; i++ {
		lastBankDetail = createRandomBankDetails(t)
	}

	arg := lastBankDetail.BankName
	bankDetails, err := testQueries.ListBankDetailsByBankName(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, bankDetails)
	for _, bankDetail := range bankDetails {
		require.NotEmpty(t, bankDetail)
		require.Equal(t, arg, bankDetail.BankName)
	}
}
