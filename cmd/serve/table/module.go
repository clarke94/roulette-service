package table

import (
	domain "github.com/clarke94/roulette-service/internal/pkg/table"
	storage "github.com/clarke94/roulette-service/storage/table"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

// Module initializes all table dependencies.
func Module(router *gin.Engine, logger *logrus.Logger, db *gorm.DB) {
	store := storage.New(db)
	controller := domain.New(logger, store)
	handler := NewHandler(controller)
	NewRouter(router, handler)
}
