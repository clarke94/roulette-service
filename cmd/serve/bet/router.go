package bet

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// NewRouter initializes all bet routes.
func NewRouter(router *gin.Engine, handler Handler) {
	v1 := router.Group("/v1")

	v1.Handle(http.MethodPost, "/table/:table/bet", handler.Create)
	v1.Handle(http.MethodGet, "/table/:table/bet", handler.List)
	v1.Handle(http.MethodPut, "/table/:table/bet", handler.Update)
	v1.Handle(http.MethodDelete, "/table/:table/bet/:bet", handler.Delete)
}
