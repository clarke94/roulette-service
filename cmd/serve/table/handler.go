package table

import (
	"context"
	"net/http"

	"github.com/clarke94/roulette-service/internal/pkg/table"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// ControllerProvider provides an interface for the domain controller.
type ControllerProvider interface {
	Create(ctx context.Context, model table.Table) (uuid.UUID, error)
	List(ctx context.Context) ([]table.Table, error)
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

// Create invokes the Create controller and returns response.
func (h Handler) Create(ctx *gin.Context) {
	var model Table
	if err := ctx.ShouldBindJSON(&model); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, err)
	}

	id, err := h.Controller.Create(ctx, presentationToDomain(model))
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, Error{Error: err.Error()})

		return
	}

	ctx.JSON(http.StatusCreated, Create{ID: id})
}

// List invokes the List controller and returns response.
func (h Handler) List(ctx *gin.Context) {
	tables, err := h.Controller.List(ctx)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, Error{Error: err.Error()})

		return
	}

	ctx.JSON(http.StatusOK, domainListToPresentation(tables))
}
