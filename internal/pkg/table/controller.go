package table

import (
	"context"
	"errors"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

var (
	ErrValidation = errors.New("invalid table")
	ErrCreate     = errors.New("unable to create table")
	ErrList       = errors.New("unable to fetch all tables")
)

// StorageProvider provides an interface to the Storage layer.
type StorageProvider interface {
	Create(ctx context.Context, model Table) (uuid.UUID, error)
	List(ctx context.Context) ([]Table, error)
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
func (c Controller) Create(ctx context.Context, model Table) (uuid.UUID, error) {
	err := c.Validator.Struct(model)
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

// List returns all tables from the storage layer.
func (c Controller) List(ctx context.Context) ([]Table, error) {
	tables, err := c.Storage.List(ctx)
	if err != nil {
		c.Logger.WithFields(logrus.Fields{
			"error": err.Error(),
		}).Error(ErrList.Error())

		return []Table{}, ErrList
	}

	return tables, nil
}
