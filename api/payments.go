package api

import (
	"fmt"
	"github.com/qwetu_petro/backend/utils"
	"github.com/qwetu_petro/backend/workers"
	"github.com/rs/zerolog/log"
	"net/http"
	"time"

	db "github.com/qwetu_petro/backend/db/sqlc"

	"github.com/gin-gonic/gin"
)

// Petty Cash Routes

type createPettyCashRequest struct {
	Amount       string `json:"amount"`
	Description  string `json:"description"`
	Folio        string `json:"folio"`
	DebitAccount string `json:"debit_account"`
}

// TO DO: modify  to automatically get the user from the context
func (server *Server) createPettyCash(ctx *gin.Context) {

	var req createPettyCashRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(400, errorResponse(err))
		return
	}

	userId, err := getUserIdFromContext(ctx)
	amount, err := utils.StringToNumeric(req.Amount)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	args := db.CreatePettyCashParams{
		Amount:          amount,
		Description:     req.Description,
		EmployeeID:      userId,
		TransactionDate: time.Now(),
		Folio:           req.Folio,
	}

	pettyCash, err := server.store.CreatePettyCashWithAudit(ctx, userId, args)
	if err != nil {
		ctx.JSON(400, errorResponse(err))
		return
	}
	// create task payload
	payload := &workers.CreatePettyCashPayload{
		PettyCashID: pettyCash.TransactionID,
	}

	// create task
	err = server.taskDistributor.DistributeTaskCreatePettyCashPdf(ctx, payload)
	if err != nil {
		ctx.JSON(400, errorResponse(err))
		return
	}
	ctx.JSON(200, pettyCash)
}

type listPettyCashRequest struct {
	PageID   int32 `form:"page_id" binding:"required,min=1"`
	PageSize int32 `form:"page_size" binding:"required,min=5,max=100"`
}

func (server *Server) listPettyCash(ctx *gin.Context) {
	var req listPettyCashRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(400, errorResponse(err))
		return
	}

	userIdInt, err := getUserIdFromContext(ctx)

	args := db.ListEmployeePettyCashParams{
		EmployeeID: userIdInt,
		Limit:      req.PageSize,
		Offset:     (req.PageID - 1) * req.PageSize,
	}
	pettyCash, err := server.store.ListEmployeePettyCash(ctx, args)
	if err != nil {
		ctx.JSON(400, errorResponse(err))
		return
	}
	ctx.JSON(200, pettyCash)

}

type updatePettyCashRequest struct {
	Amount        string `json:"amount"`
	Description   string `json:"description"`
	TransactionId int32  `json:"transaction_id"`
}

func (server *Server) updatePettyCash(ctx *gin.Context) {

	var req updatePettyCashRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(400, errorResponse(err))
		return
	}
	userIdInt, err := getUserIdFromContext(ctx)
	amount, err := utils.StringToNumeric(req.Amount)
	args := db.UpdatePettyCashParams{
		TransactionID: req.TransactionId,
		Amount:        amount,
		Description:   req.Description,
	}

	fmt.Println("----------", args)

	pettyCash, err := server.store.UpdatePettyCashWithAudit(ctx, userIdInt, args)
	if err != nil {
		// Check if the error is "no rows in result set"
		if err.Error() == "no rows in result set" {
			customErrorMessage := "The transaction is either already approved or declined. Please check again."
			ctx.JSON(http.StatusConflict, gin.H{"error": customErrorMessage})
		} else {
			log.Error().Msg(fmt.Sprintf("Error updating petty cash: %v", err))
			ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		}
		return
	}
	ctx.JSON(200, pettyCash)

}

type approvePettyCashRequest struct {
	Status        string `json:"status"`
	TransactionID int32  `json:"transaction_id"`
}

// Approve Petty Cash
func (server *Server) approvePettyCash(ctx *gin.Context) {
	// TODO check if the user is an admin

	var req approvePettyCashRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(400, errorResponse(err))
		return
	}

	userIdInt, err := getUserIdFromContext(ctx)
	if err != nil {
		ctx.JSON(400, errorResponse(err))
		return
	}
	today := time.Now().UTC()
	args := db.ApprovePettyCashParams{
		Status:        req.Status,
		AuthorisedBy:  &userIdInt,
		ApprovedAt:    &today,
		TransactionID: req.TransactionID,
	}

	pettyCash, err := server.store.ApprovePettyCashWithAudit(ctx, userIdInt, args)
	if err != nil {
		// Check if the error is "no rows in result set"
		if err.Error() == "no rows in result set" {
			customErrorMessage := "The transaction is either already approved or declined. Please check again."
			ctx.JSON(http.StatusConflict, gin.H{"error": customErrorMessage})
		} else {
			// For other errors, return a generic server error
			ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		}
		return
	}

	ctx.JSON(200, pettyCash)

}

type downloadPettyCashPdfRequest struct {
	TransactionID int32 `json:"transaction_id"`
}

// download petty cash pdf
func (server *Server) downloadPettyCashPdf(ctx *gin.Context) {
	conf := server.config
	folder := workers.PettyCashDir
	var req downloadPettyCashPdfRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(400, errorResponse(err))
		return
	}

	pettyCashDetails, err := server.store.GetPettyCash(ctx, req.TransactionID)
	if err != nil {
		ctx.JSON(400, errorResponse(err))
		return
	}

	presignedUrl, err := utils.GeneratePresignedURL(conf, folder, pettyCashDetails.PettyCashNo)
	if err != nil {
		log.Error().Err(err).Msg("error generating presigned url")
		ctx.JSON(500, errorResponse(err))
		return
	}

	ctx.JSON(200, gin.H{"url": presignedUrl})
}

// TODO improve the petty cash table to include petty cash items and their status

//Payment authorisation routes

type createPaymentRequestRequest struct {
	Amount        float64 `json:"amount"`
	Description   string  `json:"description"`
	Currency      string  `json:"currency"`
	AmountInWords string  `json:"amount_in_words"`
}

func (server *Server) createPaymentRequest(ctx *gin.Context) {

	var req createPaymentRequestRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(400, errorResponse(err))
		return
	}

	userIdInt, err := getUserIdFromContext(ctx)
	fmt.Println(userIdInt)
	args := db.CreatePaymentRequestParams{
		Amount:        req.Amount,
		Currency:      req.Currency,
		Description:   req.Description,
		EmployeeID:    userIdInt,
		Status:        "PENDING",
		AmountInWords: req.AmountInWords,
	}
	paymentRequest, err := server.store.CreatePaymentRequestWithAudit(ctx, userIdInt, args)
	if err != nil {
		ctx.JSON(400, errorResponse(err))
		return
	}

	// create task payload
	payload := &workers.CreatePaymentRequestPayload{
		PaymentRequestID: paymentRequest.RequestID,
	}

	// create task
	err = server.taskDistributor.DistributeTaskCreatePaymentRequestPdf(ctx, payload)
	if err != nil {
		ctx.JSON(400, errorResponse(err))
		return
	}
	ctx.JSON(200, paymentRequest)

}

type listPaymentRequestRequest struct {
	PageID   int32 `form:"page_id" binding:"required,min=1"`
	PageSize int32 `form:"page_size" binding:"required,min=5,max=100"`
}

func (server *Server) listPaymentRequest(ctx *gin.Context) {
	var req listPaymentRequestRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	userRoles := ctx.MustGet("roles").([]int64)
	userIdInt, err := getUserIdFromContext(ctx)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	if contains(userRoles, 1) { // admin
		args := db.ListPaymentRequestsParams{
			Limit:  req.PageSize,
			Offset: (req.PageID - 1) * req.PageSize,
		}

		paymentRequest, err := server.store.ListPaymentRequests(ctx, args)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusOK, paymentRequest)
	} else {
		args := db.ListEmployeePaymentRequestsParams{
			EmployeeID: userIdInt,
			Limit:      req.PageSize,
			Offset:     (req.PageID - 1) * req.PageSize,
		}

		paymentRequest, err := server.store.ListEmployeePaymentRequests(ctx, args)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusOK, paymentRequest)
	}
}

type updatePaymentRequestRequest struct {
	RequestId   int32   `json:"request_id"`
	Amount      float64 `json:"amount"`
	Description string  `json:"description"`
}

func (server *Server) updatePaymentRequest(ctx *gin.Context) {
	var req updatePaymentRequestRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(400, errorResponse(err))
		return
	}

	userIdInt, err := getUserIdFromContext(ctx)

	if err != nil {
		ctx.JSON(400, errorResponse(err))
		return
	}

	args := db.UpdatePaymentRequestParams{
		Amount:      req.Amount,
		Description: req.Description,
		EmployeeID:  userIdInt,
		RequestID:   req.RequestId,
	}

	paymentRequest, err := server.store.UpdatePaymentRequestWithAudit(ctx, userIdInt, args)
	if err != nil {
		// Check if the error is "no rows in result set"
		if err.Error() == "no rows in result set" {
			customErrorMessage := "The transaction is either already approved or declined. Please check again."
			ctx.JSON(http.StatusConflict, gin.H{"error": customErrorMessage})
		} else {
			// For other errors, return a generic server error
			ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		}

		// create task payload
		payload := &workers.CreatePaymentRequestPayload{
			PaymentRequestID: paymentRequest.RequestID,
		}

		// create task
		err = server.taskDistributor.DistributeTaskCreatePaymentRequestPdf(ctx, payload)
		if err != nil {
			ctx.JSON(400, errorResponse(err))
			return
		}
		return
	}

	ctx.JSON(200, paymentRequest)

}

type approvePaymentRequestRequest struct {
	Status    string `json:"status"`
	RequestId int32  `json:"request_id"`
}

func (server *Server) approvePaymentRequest(ctx *gin.Context) {
	userId, _ := ctx.Get("user_id")
	userIdInt, _ := userId.(int32)

	var req approvePaymentRequestRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(400, errorResponse(err))
		return
	}

	today := time.Now().UTC()
	args := db.ApprovePaymentRequestParams{
		Status:       req.Status,
		AdminID:      &userIdInt,
		ApprovalDate: &today,
		RequestID:    req.RequestId,
	}

	paymentRequest, err := server.store.ApprovePaymentRequestWithAudit(ctx, userIdInt, args)
	if err != nil {
		// Check if the error is "no rows in result set"
		if err.Error() == "no rows in result set" {
			customErrorMessage := "The transaction is either already approved or declined. Please check again."
			ctx.JSON(http.StatusConflict, gin.H{"error": customErrorMessage})
		} else {
			// For other errors, return a generic server error
			ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		}
		return
	}
	// create task payload
	payload := &workers.CreatePaymentRequestPayload{
		PaymentRequestID: paymentRequest.RequestID,
	}

	// create task
	err = server.taskDistributor.DistributeTaskCreatePaymentRequestPdf(ctx, payload)
	if err != nil {
		ctx.JSON(400, errorResponse(err))
		return
	}
	ctx.JSON(200, paymentRequest)

}

type downloadPaymentRequestPdfRequest struct {
	RequestId int32 `json:"request_id"`
}

func (server *Server) downloadPaymentRequestPdf(ctx *gin.Context) {
	var req downloadPaymentRequestPdfRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(400, errorResponse(err))
		return
	}
	conf := server.config
	folder := workers.PaymentRequestDir
	// todo improve this in the workers package name it better
	paymentRequestDetails, err := server.store.GetPaymentRequest(ctx, req.RequestId)

	key := fmt.Sprintf("%s.pdf", paymentRequestDetails.PaymentRequestNo)

	presignedUrl, err := utils.GeneratePresignedURL(conf, folder, key)
	if err != nil {
		log.Error().Err(err).Msg("error generating presigned url")
		ctx.JSON(500, errorResponse(err))
		return
	}

	if err != nil {
		log.Error().Err(err).Msg("error generating presigned url")
		ctx.JSON(500, errorResponse(err))
		return
	}

	ctx.JSON(200, gin.H{"url": presignedUrl})

}
