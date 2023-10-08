package api

import (
	"context"
	"fmt"
	"github.com/qwetu_petro/backend/workers"
	"github.com/rs/zerolog/log"

	db "github.com/qwetu_petro/backend/db/sqlc"
	"github.com/qwetu_petro/backend/token"
	"github.com/qwetu_petro/backend/utils"

	"github.com/gin-gonic/gin"
)

type Server struct {
	store           db.Store
	router          *gin.Engine
	tokenMaker      token.Maker
	config          utils.Config
	taskDistributor workers.TaskDistributor
}

func NewServer(config utils.Config, store db.Store, taskDistributor workers.TaskDistributor) (*Server, error) {

	log.Info().Msg("Hurray  Qwetu Just Gone Live")

	tokenMaker, err := token.NewPasetoMaker(config.TokenSymmetricKey)
	if err != nil {
		return nil, fmt.Errorf("can not create token maker: %w", err)
	}

	server := &Server{store: store, tokenMaker: tokenMaker, config: config, taskDistributor: taskDistributor}

	server.setUpRouter(config)

	return server, nil

}

func (server *Server) Start(address string) error {
	// Create an error channel to communicate errors from goroutines
	errChan := make(chan error)

	// handle http traffic in a goroutine
	go func() {
		if err := server.router.Run(address); err != nil {
			// Send the error to the error channel
			errChan <- err
		}
	}()

	// Wait for an error from either goroutine
	err := <-errChan
	log.Error().Err(err).Msg("Failed to start server")

	return err
}

func (server *Server) Shutdown(ctx context.Context) error {
	err := server.Shutdown(ctx)
	return err
}
