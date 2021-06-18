package openapi

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// NewRouter initializes all openapi routes.
func NewRouter(router *gin.Engine, handler Handler) {
	v1 := router.Group("/v1")

	v1.Handle(http.MethodGet, "/docs", handler.Docs)
	v1.Handle(http.MethodGet, "/spec", handler.Specification)
}
