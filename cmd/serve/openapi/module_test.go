package openapi

import (
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"testing"
)

func TestModule(t *testing.T) {
	tests := []struct {
		name   string
		router *gin.Engine
		logger *logrus.Logger
	}{
		{
			name:   "expect Module to init",
			router: gin.New(),
			logger: logrus.New(),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			Module(tt.router, tt.logger)
		})
	}
}
