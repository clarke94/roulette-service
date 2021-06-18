package openapi

import (
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"

	"github.com/clarke94/roulette-service/internal/pkg/openapi"
)

// Module initializes all openapi dependencies.
func Module(router *gin.Engine, logger *logrus.Logger) {
	controller := openapi.New(logger)
	handler := NewHandler(controller)
	NewRouter(router, handler)
}
