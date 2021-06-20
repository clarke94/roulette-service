package bet

import (
	domain "github.com/clarke94/roulette-service/internal/pkg/bet"
	storage "github.com/clarke94/roulette-service/storage/bet"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

// Module initializes all bet dependencies.
func Module(router *gin.Engine, logger *logrus.Logger, db *gorm.DB, validate *validator.Validate) {
	store := storage.New(db)
	controller := domain.New(logger, store, validate)
	handler := NewHandler(controller)
	NewRouter(router, handler)
}
