package openapi

import (
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"
)

// ControllerProvider provides an interface for the domain controller.
type ControllerProvider interface {
	Specification() ([]byte, error)
	Docs() ([]byte, error)
}

// Handler provides a presentation handler.
type Handler struct {
	Controller ControllerProvider
}

// NewHandler initializes a new Handler.
func NewHandler(controller ControllerProvider) Handler {
	return Handler{
		Controller: controller,
	}
}

// Docs returns the Open API documentation.
func (h Handler) Docs(ctx *gin.Context) {
	html, err := h.Controller.Docs()
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, err)

		return
	}

	ctx.Data(http.StatusOK, "text/html; charset=utf-8", html)
}

// Specification returns the Open API specification.
func (h Handler) Specification(ctx *gin.Context) {
	data, err := h.Controller.Specification()
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, err)

		return
	}

	var spec json.RawMessage

	err = json.Unmarshal(data, &spec)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, err)

		return
	}

	ctx.JSON(http.StatusOK, spec)
}
