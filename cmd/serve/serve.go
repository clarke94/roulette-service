// Package serve implements a HTTP server.
package serve

import (
	"context"
	"errors"
	"fmt"
	"github.com/clarke94/roulette-service/cmd/serve/bet"
	betStorage "github.com/clarke94/roulette-service/storage/bet"
	"net/http"
	"os/signal"
	"syscall"
	"time"

	"github.com/clarke94/roulette-service/cmd/serve/openapi"
	"github.com/clarke94/roulette-service/cmd/serve/table"
	storage "github.com/clarke94/roulette-service/storage/table"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	gormlog "gorm.io/gorm/logger"
)

const (
	serverTimeoutSeconds = 10
	serverMaxHeaderBytes = 1 << 20
	dbSlowThreshold      = 200 * time.Millisecond
)

// Handler provides a Run method when the serve command is executed.
type Handler struct{}

// New initializes the serve command.
func New() *cobra.Command {
	handler := &Handler{}

	return &cobra.Command{
		Use:   "serve",
		Short: "Starts the app server",
		Long:  ``,
		Run:   handler.Run,
	}
}

// Run will run a HTTP server and gracefully shutdown on fatal error.
func (h *Handler) Run(_ *cobra.Command, _ []string) {
	router := gin.Default()
	logger := logrus.New()
	db := h.newDatabase(logger)
	validate := validator.New()

	openapi.Module(router, logger)
	table.Module(router, logger, db, validate)
	bet.Module(router, logger, db, validate)

	h.newServer(router, logger)
}

func (h *Handler) newDatabase(logger *logrus.Logger) *gorm.DB {
	config := &gorm.Config{
		Logger: gormlog.New(
			logger,
			gormlog.Config{
				SlowThreshold: dbSlowThreshold,
				LogLevel:      gormlog.Warn,
				Colorful:      true,
			},
		),
	}

	db, err := gorm.Open(postgres.Open(viper.GetString("DATABASE_URL")), config)
	if err != nil {
		logger.WithFields(logrus.Fields{
			"error": err.Error(),
		}).Fatalln("unable to initialize database")

		return nil
	}

	err = db.AutoMigrate(storage.Table{}, betStorage.Bet{})
	if err != nil {
		logger.WithFields(logrus.Fields{
			"error": err.Error(),
		}).Fatalln("unable to migrate database")

		return nil
	}

	return db
}

func (h *Handler) newServer(router *gin.Engine, logger *logrus.Logger) {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	srv := &http.Server{
		Addr:           fmt.Sprintf(":%s", viper.GetString("PORT")),
		ReadTimeout:    serverTimeoutSeconds * time.Second,
		WriteTimeout:   serverTimeoutSeconds * time.Second,
		MaxHeaderBytes: serverMaxHeaderBytes,
		Handler:        router,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			logger.WithFields(logrus.Fields{
				"error": err.Error(),
			}).Error("unable to listen and serve")

			return
		}
	}()

	<-ctx.Done()

	stop()

	ctx, cancel := context.WithTimeout(context.Background(), serverTimeoutSeconds*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		logger.WithFields(logrus.Fields{
			"error": err.Error(),
		}).Error("server forced to shutdown")

		return
	}

	logger.Warn("server exiting")
}
