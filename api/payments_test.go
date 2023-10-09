package api

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	mockdb "github.com/qwetu_petro/backend/db/mock"
	db "github.com/qwetu_petro/backend/db/sqlc"
	"github.com/stretchr/testify/require"
	"io"
	"net"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
	"time"
)

func TestPaymentImproved(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	store := mockdb.NewMockStore(ctrl)

	paymentRequest := db.PaymentRequest{
		Amount:           123.0,
		Description:      "wow",
		Currency:         "USD",
		AmountInWords:    "One hundred",
		PaymentRequestNo: "1234",
		RequestID:        1234,
	}
	// Set up your mock store expectations
	store.EXPECT().ApprovePaymentRequestWithAudit(gomock.Any(), gomock.Any(), gomock.Any()).Return(paymentRequest, nil)

	r := gin.Default()
	server := newTestServer(t, store)
	r.POST("/payment-request/approve", func(c *gin.Context) {
		userId := int64(1234)
		c.Set("user_id", userId)
		server.approvePaymentRequest(c)
	})

	body := `{ "amount": "123", "description": "wow", "currency": "USD", "amount_in_words": "One hundred" }`
	req, err := http.NewRequest(http.MethodPost, "/payment-request/approve", bytes.NewBufferString(body))
	require.NoError(t, err)

	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	require.Equal(t, http.StatusOK, w.Code)

	var expectedPaymentRequest db.PaymentRequest
	err = json.Unmarshal(w.Body.Bytes(), &expectedPaymentRequest)
	if err != nil {
		t.Error(err)
	}

	require.Equal(t, paymentRequest, expectedPaymentRequest)
	require.NoError(t, err)
}

func TestEndToEnd(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	store := mockdb.NewMockStore(ctrl)
	srv := newTestServer(t, store)
	store.EXPECT().GetUserByUserNameOrEmail(gomock.Any(), gomock.Any()).Return(db.User{}, nil)
	store.EXPECT().GetUserRoles(gomock.Any(), gomock.Any()).Return([]db.UserRole{}, nil)

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
	store.EXPECT().ApprovePaymentRequestWithAudit(gomock.Any(), gomock.Any(), gomock.Any()).Return(paymentRequest, nil)

	go func() {
		if err := srv.Start("0.0.0.0:8085"); !errors.Is(err, http.ErrServerClosed) {
			t.Errorf("ListenAndServe(): %v", err)
			return
		}
	}()

	request, err := http.NewRequest(
		"POST",
		"http://localhost:8085/payment-request/approve",
		bytes.NewBufferString(`{ "amount": "123", "description": "Description", "currency": "USD", "amount_in_words": "One hundred" }`),
	)
	if err != nil {
		return
	}

	token, _, err := srv.tokenMaker.CreateToken("username", 300000000)
	if err != nil {
		return
	}
	request.Header.Set("Authorization", "Bearer "+token)

	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		t.Fatalf("Could not make request: %v", err)
	}

	if response.StatusCode != 200 {
		t.Fatalf("Expected status 200, got %v", response.StatusCode)
	}

	bodyBytes, err := io.ReadAll(response.Body)
	if err != nil {
		t.Fatalf("Failed to read response body: %v", err)
	}
	var paymentRequestFromResponse db.PaymentRequest
	err = json.Unmarshal(bodyBytes, &paymentRequestFromResponse)
	if err != nil {
		t.Fatalf("Failed to unmarshal JSON: %v", err)
	}

	requireEqualPaymentRequestsByName(t, paymentRequest, paymentRequestFromResponse)
	require.NoError(t, err)
}

func requireEqualPaymentRequestsByName(t *testing.T, expected, actual db.PaymentRequest) {
	fieldNames := []string{
		"RequestID",
		"PaymentRequestNo",
		"AmountInWords",
		"EmployeeID",
		"Currency",
		"Amount",
		"Description",
		"Status",
	}

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
