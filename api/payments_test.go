package api

import (
	"bytes"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	mockdb "github.com/qwetu_petro/backend/db/mock"
	db "github.com/qwetu_petro/backend/db/sqlc"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"testing"
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
