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
	Update(ctx context.Context, model table.Table) (uuid.UUID, error)
	Delete(ctx context.Context, id uuid.UUID) (uuid.UUID, error)
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

		return
	}

	id, err := h.Controller.Create(ctx, presentationToDomain(model))
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, Error{Error: err.Error()})

		return
	}

	ctx.JSON(http.StatusCreated, Upsert{ID: id})
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

// Update invokes the Update controller and returns response.
func (h Handler) Update(ctx *gin.Context) {
	var model Table
	if err := ctx.ShouldBindJSON(&model); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, err)

		return
	}

	id, err := h.Controller.Update(ctx, presentationToDomain(model))
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, Error{Error: err.Error()})

		return
	}

	ctx.JSON(http.StatusOK, Upsert{ID: id})
}

// Delete invokes the Delete controller and returns an id.
func (h Handler) Delete(ctx *gin.Context) {
	id, err := uuid.Parse(ctx.Param("table"))
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, Error{Error: err.Error()})

		return
	}

	deletedID, err := h.Controller.Delete(ctx, id)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, Error{Error: err.Error()})

		return
	}

	ctx.JSON(http.StatusOK, Upsert{ID: deletedID})
}
