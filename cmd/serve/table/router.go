package table

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// NewRouter initializes all table routes.
func NewRouter(router *gin.Engine, handler Handler) {
	v1 := router.Group("/v1")

	v1.Handle(http.MethodPost, "/table", handler.Create)
	v1.Handle(http.MethodGet, "/table", handler.List)
	v1.Handle(http.MethodPut, "/table", handler.Update)
}
