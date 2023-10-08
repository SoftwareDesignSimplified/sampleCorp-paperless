package workers

import (
	"context"
	"github.com/hibiken/asynq"
	db "github.com/qwetu_petro/backend/db/sqlc"
	"github.com/qwetu_petro/backend/mail"
	"github.com/qwetu_petro/backend/sms"
	"github.com/qwetu_petro/backend/utils"
	"github.com/rs/zerolog/log"
)

type TaskProcessor interface {
	Start() error
	ProcessTaskCreateInvoicePdf(ctx context.Context, task *asynq.Task) error
	ProcessTaskCreatePaymentRequestPdf(ctx context.Context, task *asynq.Task) error
	ProcessTaskCreatePettyCashPdf(ctx context.Context, task *asynq.Task) error
	ProcessTaskCreatePurchaseOrderPdf(ctx context.Context, task *asynq.Task) error
	ProcessTaskCreateQuotationPdf(ctx context.Context, task *asynq.Task) error
}

type RedisTaskProcessor struct {
	server    *asynq.Server
	db        db.Store
	smsSender sms.SmsSender
	mailer    mail.EmailSender
	conf      utils.Config
}

func NewRedisTaskProcessor(redisOpts asynq.RedisClientOpt, db db.Store) TaskProcessor {
	conf, err := utils.LoadConfig("../..")
	if err != nil {
		log.Fatal().Err(err).Msg("Error loading config")
	}
	server := asynq.NewServer(redisOpts,
		asynq.Config{
			Concurrency: -1, // unlimited concurrency (default)
			ErrorHandler: asynq.ErrorHandlerFunc(
				func(ctx context.Context, task *asynq.Task, err error) {
					log.Error().Err(err).Msgf("failed to process task %s with payload %v", task.Type(), task.Payload())
				}),
			Logger: NewLogger(),
		})
	return &RedisTaskProcessor{server: server, db: db, conf: conf}
}

func (processor *RedisTaskProcessor) Start() error {
	mux := asynq.NewServeMux()
	mux.HandleFunc(TaskTypeCreateInvoicePdf, processor.ProcessTaskCreateInvoicePdf)
	mux.HandleFunc(TaskTypeCreatePaymentRequestPdf, processor.ProcessTaskCreatePaymentRequestPdf)
	mux.HandleFunc(TaskTypeCreatePettyCashPdf, processor.ProcessTaskCreatePettyCashPdf)
	mux.HandleFunc(TaskTypeCreatePurchaseOrderPdf, processor.ProcessTaskCreatePurchaseOrderPdf)
	mux.HandleFunc(TaskTypeCreateQuotationPdf, processor.ProcessTaskCreateQuotationPdf)
	return processor.server.Start(mux)
}
