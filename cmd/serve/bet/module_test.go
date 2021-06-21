package bet

import (
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"testing"
)

func TestModule(t *testing.T) {
	tests := []struct {
		name     string
		router   *gin.Engine
		logger   *logrus.Logger
		db       *gorm.DB
		validate *validator.Validate
	}{
		{
			name:     "expect Module to init",
			router:   gin.New(),
			logger:   logrus.New(),
			db:       &gorm.DB{},
			validate: validator.New(),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			Module(tt.router, tt.logger, tt.db)
		})
	}
}
