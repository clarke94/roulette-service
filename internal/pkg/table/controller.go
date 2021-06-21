package table

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

var (
	ErrCreate = errors.New("unable to create table")
	ErrList   = errors.New("unable to fetch all tables")
	ErrUpdate = errors.New("unable to update table")
	ErrDelete = errors.New("unable to delete table")
)

// StorageProvider provides an interface to the Storage layer.
type StorageProvider interface {
	Create(ctx context.Context, model Table) (string, error)
	List(ctx context.Context) ([]Table, error)
	Update(ctx context.Context, model Table) (string, error)
	Delete(ctx context.Context, id string) (string, error)
}

// Controller provides a domain controller.
type Controller struct {
	Logger  *logrus.Logger
	Storage StorageProvider
}

// New initializes a new Controller.
func New(logger *logrus.Logger, storage StorageProvider) Controller {
	return Controller{
		Logger:  logger,
		Storage: storage,
	}
}

// Create validates the model and invokes the repository.
func (c Controller) Create(ctx context.Context, model Table) (string, error) {
	model.ID = uuid.New().String()

	id, err := c.Storage.Create(ctx, model)
	if err != nil {
		c.Logger.WithFields(logrus.Fields{
			"error": err.Error(),
		}).Error(ErrCreate.Error())

		return "", ErrCreate
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

// Update validates the model and invokes the repository.
func (c Controller) Update(ctx context.Context, model Table) (string, error) {
	id, err := c.Storage.Update(ctx, model)
	if err != nil {
		c.Logger.WithFields(logrus.Fields{
			"error": err.Error(),
		}).Error(ErrUpdate.Error())

		return "", ErrUpdate
	}

	return id, nil
}

// Delete deletes one from the repository.
func (c Controller) Delete(ctx context.Context, id string) (string, error) {
	deletedID, err := c.Storage.Delete(ctx, id)
	if err != nil {
		c.Logger.WithFields(logrus.Fields{
			"error": err.Error(),
		}).Error(ErrDelete.Error())

		return "", ErrDelete
	}

	return deletedID, nil
}
