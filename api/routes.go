package api

import (
	"github.com/gin-contrib/logger"
	"github.com/gin-contrib/requestid"
	"github.com/newrelic/go-agent/v3/integrations/nrgin"
	"github.com/qwetu_petro/backend/newrelic"
	"github.com/rs/zerolog/log"
	"io"
	"strings"
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

	log.Info().Msg("CORS ALLOWED ORIGIN: " + strings.Join(allowOrigins, ","))

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
	loggedInRoutes := router.Group("/").Use(authMiddleware(server.tokenMaker)).Use(server.currentUserMiddleware())

	// User routes
	//loggedInRoutes.POST("/users", server.createUser)
	loggedInRoutes.GET("/users", server.listUsers)

	// Petty Cash Routes

	loggedInRoutes.POST("/petty-cash", server.createPettyCash)
	loggedInRoutes.GET("/petty-cash", server.listPettyCash)
	loggedInRoutes.PUT("/petty-cash", server.updatePettyCash)
	loggedInRoutes.POST("/petty-cash/approve", server.approvePettyCash)
	loggedInRoutes.GET("/petty-cash/download", server.downloadPettyCashPdf)
	//Payment authorisation routes
	loggedInRoutes.POST("/payment-request", server.createPaymentRequest)
	loggedInRoutes.GET("/payment-request", server.listPaymentRequest)
	loggedInRoutes.PUT("/payment-request", server.updatePaymentRequest)

	// Admin Routes
	adminRoutes := loggedInRoutes.Use(server.mustBeAdminMiddleware())
	adminRoutes.POST("/payment-request/approve", server.approvePaymentRequest)

	loggedInRoutes.GET("/payment-request/download", server.downloadPaymentRequestPdf)

	//Invoice routes
	loggedInRoutes.POST("/invoice", server.createInvoice)
	loggedInRoutes.GET("/invoice", server.listInvoice)
	loggedInRoutes.PUT("/invoice", server.updateInvoice)
	loggedInRoutes.POST("/invoice/approve", server.approveInvoice)
	loggedInRoutes.GET("/invoice/download", server.downloadInvoicePdf)

	// Purchase Order routes
	loggedInRoutes.POST("/purchase_order", server.createPurchaseOrder)
	loggedInRoutes.GET("/purchase_order", server.listPurchaseOrder)
	loggedInRoutes.PUT("/purchase_order", server.updatePurchaseOrder)
	// TO DO: Ask if this is required
	loggedInRoutes.POST("/purchase_order/approve", server.approvePurchaseOrder)
	loggedInRoutes.GET("/purchase_order/download", server.downloadPurchaseOrderPdf)

	// Roles routes
	loggedInRoutes.POST("/roles", server.createRole)
	loggedInRoutes.GET("/roles", server.listRoles)
	loggedInRoutes.POST("/roles/:id", server.updateRole)
	loggedInRoutes.POST("/roles/delete", server.deleteRole)

	// UserRoles routes
	loggedInRoutes.POST("/user-roles", server.createUserRole)
	//loggedInRoutes.GET("/user-roles", server.listUserRoles)
	//loggedInRoutes.POST("/user-roles/:id", server.updateUserRole)
	loggedInRoutes.POST("/users-roles/delete/:id", server.deleteUserRole)

	// Quotation routes
	loggedInRoutes.POST("/quotation", server.createQuotation)
	loggedInRoutes.GET("/quotation", server.listQuotations)
	loggedInRoutes.PUT("/quotation", server.updateQuotation)
	//loggedInRoutes.POST("/quotation/approve", server.approveQuotation)
	loggedInRoutes.GET("/quotation/download", server.downloadQuotationPdf)

	// Bank Details routes
	loggedInRoutes.POST("/bank-details", server.createBankDetails)
	loggedInRoutes.GET("/bank-details", server.getBankDetails)
	loggedInRoutes.GET("/banks", server.listBanks)

	// Company Details routes
	loggedInRoutes.POST("/company", server.createCompany)
	loggedInRoutes.GET("/company/delete", server.deleteCompany)
	loggedInRoutes.GET("/company", server.listCompanies)

	// Signatory routes
	loggedInRoutes.POST("/signatory", server.createSignatory)
	loggedInRoutes.GET("/signatory", server.listSignatories)

	// Car routes
	loggedInRoutes.POST("/car", server.createCar)
	loggedInRoutes.GET("/car", server.listCars)

	// Fuel routes
	loggedInRoutes.POST("/fuel", server.createFuelConsumption)
	loggedInRoutes.GET("/fuel", server.listFuelConsumption)
	loggedInRoutes.GET("/fuel/car/:id", server.getCarFuelConsumption)
	router.GET("/fuel/range", server.getCarFuelByDateRange)

	server.router = router

}
