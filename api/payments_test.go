package api

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"github.com/golang/mock/gomock"
	mockdb "github.com/qwetu_petro/backend/db/mock"
	db "github.com/qwetu_petro/backend/db/sqlc"
	"github.com/stretchr/testify/require"
	"io"
	"net"
	"net/http"
	"reflect"
	"testing"
	"time"
)

func TestGivenAUserWhoIsNotAnAdmin_WhenRequestingAPaymentApproval_ThenReturnUnauthorized(t *testing.T) {
	store, srv := setupTestServerAndMockStore(t)
	store.EXPECT().GetUserByUserNameOrEmail(gomock.Any(), gomock.Any()).Return(db.User{}, nil)
	store.EXPECT().GetUserRoles(gomock.Any(), gomock.Any()).Return([]db.UserRole{}, nil)
	store.EXPECT().GetUserRoles(gomock.Any(), gomock.Any()).Return([]db.UserRole{}, nil)
	paymentRequest := createRandomPaymentRequest()
	store.EXPECT().ApprovePaymentRequestWithAudit(gomock.Any(), gomock.Any(), gomock.Any()).Return(paymentRequest, nil)

	response, err := makePaymentRequestApprovalRequest(t, srv)
	require.NoError(t, err)

	require.Equal(t, http.StatusUnauthorized, response.StatusCode)
	bodyBytes, err := io.ReadAll(response.Body)
	if err != nil {
		t.Fatalf("Failed to read response body: %v", err)
	}
	require.Equal(t, "\"you don't have permission to access this resource\"", string(bodyBytes))
}

func TestPaymentApprovalAPIApprovedForAnyUser(t *testing.T) {
	mockUserRoles := []db.UserRole{
		{
			ID: 1,
		},
	}

	store, srv := setupTestServerAndMockStore(t)
	store.EXPECT().GetUserByUserNameOrEmail(gomock.Any(), gomock.Any()).Return(db.User{}, nil)
	store.EXPECT().GetUserRoles(gomock.Any(), gomock.Any()).Return(mockUserRoles, nil)
	store.EXPECT().GetUserRoles(gomock.Any(), gomock.Any()).Return(mockUserRoles, nil)
	paymentRequest := createRandomPaymentRequest()
	store.EXPECT().ApprovePaymentRequestWithAudit(gomock.Any(), gomock.Any(), gomock.Any()).Return(paymentRequest, nil)

	response, err := makePaymentRequestApprovalRequest(t, srv)
	require.NoError(t, err)

	require.Equal(t, http.StatusOK, response.StatusCode)
	returnedPaymentRequest := extractPaymentRequestFromResponse(t, err, response)

	fieldsToCompare := []string{
		"RequestID",
		"PaymentRequestNo",
		"AmountInWords",
		"EmployeeID",
		"Currency",
		"Amount",
		"Description",
		"Status",
	}
	requirePaymentRequestFieldsEqual(t, fieldsToCompare, paymentRequest, returnedPaymentRequest)
}

func extractPaymentRequestFromResponse(t *testing.T, err error, response *http.Response) db.PaymentRequest {
	bodyBytes, err := io.ReadAll(response.Body)
	if err != nil {
		t.Fatalf("Failed to read response body: %v", err)
	}
	var paymentRequestFromResponse db.PaymentRequest
	err = json.Unmarshal(bodyBytes, &paymentRequestFromResponse)
	if err != nil {
		t.Fatalf("Failed to unmarshal JSON: %v", err)
	}
	return paymentRequestFromResponse
}

func makePaymentRequestApprovalRequest(t *testing.T, srv *Server) (*http.Response, error) {
	request, err := http.NewRequest(
		"POST",
		"http://localhost:8085/payment-request/approve",
		bytes.NewBufferString(`{ "amount": "123", "description": "Description", "currency": "USD", "amount_in_words": "One hundred" }`),
	)
	if err != nil {
		return nil, err
	}

	token, _, err := srv.tokenMaker.CreateToken("username", 300000000)
	if err != nil {
		return nil, err
	}
	request.Header.Set("Authorization", "Bearer "+token)

	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		t.Fatalf("Could not make request: %v", err)
	}
	return response, err
}

func createRandomPaymentRequest() db.PaymentRequest {
	adminID := int32(1234)
	approvalDate := time.Now()
	invoiceID := int32(1234)
	pdfUrl := "http://example.org/pdf"

	paymentRequest := db.PaymentRequest{
		RequestID:        1234,
		PaymentRequestNo: "1234",
		AmountInWords:    "One hundred",
		EmployeeID:       6789,
		Currency:         "USD",
		Amount:           100.00,
		Description:      "Description",
		RequestDate:      time.Now(),
		Status:           "Approved",
		AdminID:          &adminID,
		ApprovalDate:     &approvalDate,
		InvoiceID:        &invoiceID,
		PdfUrl:           &pdfUrl,
	}
	return paymentRequest
}

func setupTestServerAndMockStore(t *testing.T) (*mockdb.MockStore, *Server) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	store := mockdb.NewMockStore(ctrl)
	srv := newTestServer(t, store)

	go func() {
		if err := srv.Start("0.0.0.0:8085"); !errors.Is(err, http.ErrServerClosed) {
			t.Errorf("ListenAndServe(): %v", err)
			return
		}
	}()

	return store, srv
}

func requirePaymentRequestFieldsEqual(t *testing.T, fieldNames []string, expected, actual db.PaymentRequest) {
	valExpected := reflect.ValueOf(expected)
	valActual := reflect.ValueOf(actual)

	for _, fieldName := range fieldNames {
		fieldExpected := valExpected.FieldByName(fieldName)
		fieldActual := valActual.FieldByName(fieldName)
		require.Equal(t, fieldExpected.Interface(), fieldActual.Interface(), "Field mismatch: %s", fieldName)
	}
}

func ShutdownServer(t *testing.T, srv *Server) {
	if err := srv.Shutdown(context.Background()); err != nil {
		t.Fatalf("Server Shutdown Failed:%v", err)
	}
}

func StartServer(t *testing.T) *Server {
	srv := newTestServer(t, nil)
	ready := make(chan struct{})

	go func() {
		ln, err := net.Listen("tcp", "localhost:8081")
		if err != nil {
			t.Errorf("Could not listen: %v", err)
			return
		}
		close(ready) // Signal that we are ready

		srv := &http.Server{}
		if err := srv.Serve(ln); !errors.Is(err, http.ErrServerClosed) {
			t.Errorf("ListenAndServe(): %v", err)
			return
		}
	}()

	<-ready
	return srv
}
