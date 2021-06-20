package bet

import (
	"context"
	"errors"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

var (
	ErrValidation = errors.New("invalid bet")
	ErrCreate     = errors.New("unable to create bet")
	ErrList       = errors.New("unable to fetch all bets")
	ErrUpdate     = errors.New("unable to update bet")
	ErrDelete     = errors.New("unable to delete bet")
)

// StorageProvider provides an interface to the Storage layer.
type StorageProvider interface {
	Create(ctx context.Context, model Bet) (uuid.UUID, error)
	List(ctx context.Context, tableID uuid.UUID) ([]Bet, error)
	Update(ctx context.Context, model Bet) (uuid.UUID, error)
	Delete(ctx context.Context, tableID, id uuid.UUID) (uuid.UUID, error)
}

// Controller provides a domain controller.
type Controller struct {
	Logger    *logrus.Logger
	Validator *validator.Validate
	Storage   StorageProvider
}

// New initializes a new Controller.
func New(logger *logrus.Logger, storage StorageProvider, validate *validator.Validate) Controller {
	return Controller{
		Logger:    logger,
		Validator: validate,
		Storage:   storage,
	}
}

// Create validates the model and invokes the repository.
func (c Controller) Create(ctx context.Context, model Bet) (uuid.UUID, error) {
	model.ID = uuid.New()

	err := c.Validator.Struct(model)
	if err != nil {
		c.Logger.WithFields(logrus.Fields{
			"error": err.Error(),
		}).Error(ErrValidation.Error())

		return uuid.Nil, ErrValidation
	}

	err = c.validate(model)
	if err != nil {
		c.Logger.WithFields(logrus.Fields{
			"error": err.Error(),
		}).Error(ErrValidation.Error())

		return uuid.Nil, ErrValidation
	}

	id, err := c.Storage.Create(ctx, model)
	if err != nil {
		c.Logger.WithFields(logrus.Fields{
			"error": err.Error(),
		}).Error(ErrCreate.Error())

		return uuid.Nil, ErrCreate
	}

	return id, nil
}

// List returns all bets from the storage layer.
func (c Controller) List(ctx context.Context, tableID uuid.UUID) ([]Bet, error) {
	bets, err := c.Storage.List(ctx, tableID)
	if err != nil {
		c.Logger.WithFields(logrus.Fields{
			"error": err.Error(),
		}).Error(ErrList.Error())

		return []Bet{}, ErrList
	}

	return bets, nil
}

// Update validates the model and invokes the repository.
func (c Controller) Update(ctx context.Context, model Bet) (uuid.UUID, error) {
	err := c.Validator.Struct(Update{
		ID:  model.ID,
		Bet: model,
	})
	if err != nil {
		c.Logger.WithFields(logrus.Fields{
			"error": err.Error(),
		}).Error(ErrValidation.Error())

		return uuid.Nil, ErrValidation
	}

	err = c.validate(model)
	if err != nil {
		c.Logger.WithFields(logrus.Fields{
			"error": err.Error(),
		}).Error(ErrValidation.Error())

		return uuid.Nil, ErrValidation
	}

	id, err := c.Storage.Update(ctx, model)
	if err != nil {
		c.Logger.WithFields(logrus.Fields{
			"error": err.Error(),
		}).Error(ErrUpdate.Error())

		return uuid.Nil, ErrUpdate
	}

	return id, nil
}

// Delete deletes one from the repository.
func (c Controller) Delete(ctx context.Context, tableID, id uuid.UUID) (uuid.UUID, error) {
	err := c.Validator.Var(id, "required")
	if err != nil {
		c.Logger.WithFields(logrus.Fields{
			"error": err.Error(),
		}).Error(ErrValidation.Error())

		return uuid.Nil, ErrValidation
	}

	deletedID, err := c.Storage.Delete(ctx, tableID, id)
	if err != nil {
		c.Logger.WithFields(logrus.Fields{
			"error": err.Error(),
		}).Error(ErrDelete.Error())

		return uuid.Nil, ErrDelete
	}

	return deletedID, nil
}

func (c Controller) validate(b Bet) error {
	if _, ok := TypeMultiplierMap[b.Type]; !ok {
		return ErrValidation
	}

	switch b.Type {
	case TypeRedBlack:
		return c.Validator.Var(b.Bet, "oneof=red black")
	case TypeOddEven:
		return c.Validator.Var(b.Bet, "oneof=odd even")
	case TypeHighLow:
		return c.Validator.Var(b.Bet, "oneof=high low")
	case TypeColumn:
		return c.Validator.Var(b.Bet, "oneof=1 2 3")
	case TypeDozen:
		return c.Validator.Var(b.Bet, "oneof=1-12 13-24 25-36")
	case TypeStraight:
		return c.Validator.Var(b.Bet, "number,gte=0,lte=36")
	}

	return nil
}
