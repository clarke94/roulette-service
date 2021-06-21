package bet

import (
	"context"
	"net/http"

	"github.com/clarke94/roulette-service/internal/pkg/bet"
	"github.com/gin-gonic/gin"
)

// ControllerProvider provides an interface for the domain controller.
type ControllerProvider interface {
	Create(ctx context.Context, model bet.Bet) (string, error)
	List(ctx context.Context, tableID string) ([]bet.Bet, error)
	Update(ctx context.Context, model bet.Bet) (string, error)
	Delete(ctx context.Context, tableID, id string) (string, error)
	Play(ctx context.Context, tableID string) (bet.Result, error)
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
	var params TableParam
	if err := ctx.BindUri(&params); err != nil {
		return
	}

	var model Bet
	if err := ctx.BindJSON(&model); err != nil {
		return
	}

	id, err := h.Controller.Create(ctx, presentationToDomain(model, params.Table))
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, Error{Error: err.Error()})

		return
	}

	ctx.JSON(http.StatusCreated, Upsert{ID: id})
}

// List invokes the List controller and returns response.
func (h Handler) List(ctx *gin.Context) {
	var params TableParam
	if err := ctx.BindUri(&params); err != nil {
		return
	}

	bets, err := h.Controller.List(ctx, params.Table)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, Error{Error: err.Error()})

		return
	}

	ctx.JSON(http.StatusOK, domainListToPresentation(bets))
}

// Update invokes the Update controller and returns response.
func (h Handler) Update(ctx *gin.Context) {
	var params TableParam
	if err := ctx.BindUri(&params); err != nil {
		return
	}

	var model Update
	if err := ctx.BindJSON(&model); err != nil {
		return
	}

	domainModel := presentationToDomain(Bet{
		ID:       model.ID,
		Bet:      model.Bet.Bet,
		Type:     model.Bet.Type,
		Amount:   model.Bet.Amount,
		Currency: model.Bet.Currency,
	}, params.Table)

	id, err := h.Controller.Update(ctx, domainModel)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, Error{Error: err.Error()})

		return
	}

	ctx.JSON(http.StatusOK, Upsert{ID: id})
}

// Delete invokes the Delete controller and returns an id.
func (h Handler) Delete(ctx *gin.Context) {
	var tableParam TableParam
	if err := ctx.BindUri(&tableParam); err != nil {
		return
	}

	var betParam IDParam
	if err := ctx.BindUri(&betParam); err != nil {
		return
	}

	deletedID, err := h.Controller.Delete(ctx, tableParam.Table, betParam.Bet)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, Error{Error: err.Error()})

		return
	}

	ctx.JSON(http.StatusOK, Upsert{ID: deletedID})
}

// Play invokes the Play controller and returns response.
func (h Handler) Play(ctx *gin.Context) {
	var params TableParam
	if err := ctx.BindUri(&params); err != nil {
		return
	}

	results, err := h.Controller.Play(ctx, params.Table)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, Error{Error: err.Error()})

		return
	}

	ctx.JSON(http.StatusOK, domainResultToDomain(results))
}
