package api

import (
	"fmt"
	"github.com/gin-contrib/logger"
	"github.com/gin-contrib/requestid"
	"github.com/newrelic/go-agent/v3/integrations/nrgin"
	"github.com/qwetu_petro/backend/newrelic"
	"github.com/rs/zerolog/log"
	"io"
	"time"

	"github.com/qwetu_petro/backend/utils"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func (server *Server) setUpRouter(config utils.Config) {
	router := gin.Default()

	var allowOrigins []string

	if !config.SendLogsToStdOut {
		log.Logger = log.Output(io.Discard)
	}
	log.Info().Msg("Setting up router")
	log.Info().Msg("ENVIRONMENT: " + config.Environment)
	// check if ENVIRONMENT is set to production or development
	if config.Environment == "production" {
		allowOrigins = []string{config.CorsAllowedOrigin}
	} else {
		allowOrigins = []string{"*"}
	}

	// CORS
	fmt.Println("CORS ALLOWED ORIGIN: ", allowOrigins)

	router.Use(cors.New(cors.Config{
		AllowOrigins:     allowOrigins,
		AllowMethods:     []string{"PUT", "PATCH", "POST", "GET", "DELETE"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Content-Length", "Accept-Encoding", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	// set up newrelic
	newRelicApp, err := newrelic.RelicApp(config)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to initialize newrelic")
	}
	router.Use(nrgin.Middleware(newRelicApp))

	//router.Use(RequestMetrics())
	router.Use(requestid.New())
	router.Use(logger.SetLogger())

	// Routes  that don't require authentication
	router.GET("/auth/current_user", server.currentUserMiddleware(), server.currentUser)
	router.POST("/auth/users", server.createUser) // keep here for the time being
	// @todo: add auth to this route. Only authenticated users should have access to create users
	router.POST("/auth/users/login", server.loginUser)

	// access token routes
	router.POST("/token/renew_access", server.renewAccessToken)

	// Routes that require authentication
	authRoutes := router.Group("/").Use(authMiddleware(server.tokenMaker)).Use(server.currentUserMiddleware())

	// User routes
	//authRoutes.POST("/users", server.createUser)
	authRoutes.GET("/users", server.listUsers)

	// Petty Cash Routes

	authRoutes.POST("/petty-cash", server.createPettyCash)
	authRoutes.GET("/petty-cash", server.listPettyCash)
	authRoutes.PUT("/petty-cash", server.updatePettyCash)
	authRoutes.POST("/petty-cash/approve", server.approvePettyCash)
	authRoutes.GET("/petty-cash/download", server.downloadPettyCashPdf)
	//Payment authorisation routes
	authRoutes.POST("/payment-request", server.createPaymentRequest)
	authRoutes.GET("/payment-request", server.listPaymentRequest)
	authRoutes.PUT("/payment-request", server.updatePaymentRequest)
	authRoutes.POST("/payment-request/approve", server.approvePaymentRequest)
	authRoutes.GET("/payment-request/download", server.downloadPaymentRequestPdf)

	//Invoice routes
	authRoutes.POST("/invoice", server.createInvoice)
	authRoutes.GET("/invoice", server.listInvoice)
	authRoutes.PUT("/invoice", server.updateInvoice)
	authRoutes.POST("/invoice/approve", server.approveInvoice)
	authRoutes.GET("/invoice/download", server.downloadInvoicePdf)

	// Purchase Order routes
	authRoutes.POST("/purchase_order", server.createPurchaseOrder)
	authRoutes.GET("/purchase_order", server.listPurchaseOrder)
	authRoutes.PUT("/purchase_order", server.updatePurchaseOrder)
	// TO DO: Ask if this is required
	authRoutes.POST("/purchase_order/approve", server.approvePurchaseOrder)
	authRoutes.GET("/purchase_order/download", server.downloadPurchaseOrderPdf)

	// Roles routes
	authRoutes.POST("/roles", server.createRole)
	authRoutes.GET("/roles", server.listRoles)
	authRoutes.POST("/roles/:id", server.updateRole)
	authRoutes.POST("/roles/delete", server.deleteRole)

	// UserRoles routes
	authRoutes.POST("/user-roles", server.createUserRole)
	//authRoutes.GET("/user-roles", server.listUserRoles)
	//authRoutes.POST("/user-roles/:id", server.updateUserRole)
	authRoutes.POST("/users-roles/delete/:id", server.deleteUserRole)

	// Quotation routes
	authRoutes.POST("/quotation", server.createQuotation)
	authRoutes.GET("/quotation", server.listQuotations)
	authRoutes.PUT("/quotation", server.updateQuotation)
	//authRoutes.POST("/quotation/approve", server.approveQuotation)
	authRoutes.GET("/quotation/download", server.downloadQuotationPdf)

	// Bank Details routes
	authRoutes.POST("/bank-details", server.createBankDetails)
	authRoutes.GET("/bank-details", server.getBankDetails)
	authRoutes.GET("/banks", server.listBanks)

	// Company Details routes
	authRoutes.POST("/company", server.createCompany)
	authRoutes.GET("/company/delete", server.deleteCompany)
	authRoutes.GET("/company", server.listCompanies)

	// Signatory routes
	authRoutes.POST("/signatory", server.createSignatory)
	authRoutes.GET("/signatory", server.listSignatories)

	// Car routes
	authRoutes.POST("/car", server.createCar)
	authRoutes.GET("/car", server.listCars)

	// Fuel routes
	authRoutes.POST("/fuel", server.createFuelConsumption)
	authRoutes.GET("/fuel", server.listFuelConsumption)
	authRoutes.GET("/fuel/car/:id", server.getCarFuelConsumption)
	router.GET("/fuel/range", server.getCarFuelByDateRange)

	server.router = router

}
