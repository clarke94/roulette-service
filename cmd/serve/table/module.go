package table

import (
	domain "github.com/clarke94/roulette-service/internal/pkg/table"
	storage "github.com/clarke94/roulette-service/storage/table"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

// Module initializes all table dependencies.
func Module(router *gin.Engine, logger *logrus.Logger, db *gorm.DB, validate *validator.Validate) {
	store := storage.New(db)
	controller := domain.New(logger, store, validate)
	handler := NewHandler(controller)
	NewRouter(router, handler)
}
