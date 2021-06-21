package bet

import (
	"context"
	"crypto/rand"
	"errors"
	"math/big"
	"strconv"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

var (
	ErrCreate = errors.New("unable to create bet")
	ErrList   = errors.New("unable to fetch all bets")
	ErrUpdate = errors.New("unable to update bet")
	ErrDelete = errors.New("unable to delete bet")
)

// StorageProvider provides an interface to the Storage layer.
type StorageProvider interface {
	Create(ctx context.Context, model Bet) (string, error)
	List(ctx context.Context, tableID string, filters ...Bet) ([]Bet, error)
	Update(ctx context.Context, model Bet) (string, error)
	Delete(ctx context.Context, tableID, id string) (string, error)
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
func (c Controller) Create(ctx context.Context, model Bet) (string, error) {
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

// List returns all bets from the storage layer.
func (c Controller) List(ctx context.Context, tableID string) ([]Bet, error) {
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
func (c Controller) Update(ctx context.Context, model Bet) (string, error) {
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
func (c Controller) Delete(ctx context.Context, tableID, id string) (string, error) {
	deletedID, err := c.Storage.Delete(ctx, tableID, id)
	if err != nil {
		c.Logger.WithFields(logrus.Fields{
			"error": err.Error(),
		}).Error(ErrDelete.Error())

		return "", ErrDelete
	}

	return deletedID, nil
}

// Play runs the roulette algorithm, clears the table and returns the winners.
func (c Controller) Play(ctx context.Context, tableID string) (Result, error) {
	number := c.getNumber()
	color := c.getColor(number)

	filters := c.winnerFilters(number, color)

	bets, err := c.Storage.List(ctx, tableID, filters...)
	if err != nil {
		c.Logger.WithFields(logrus.Fields{
			"error": err.Error(),
		}).Error(ErrList.Error())

		return Result{}, ErrList
	}

	result := Result{
		Number:  number,
		Color:   color,
		Winners: betListToWinner(bets),
	}

	return result, nil
}

func (c Controller) winnerFilters(number int, color string) []Bet {
	filters := make([]Bet, 0)

	filters = append(filters, Bet{
		Bet:  strconv.Itoa(number),
		Type: TypeStraight,
	})

	filters = append(filters, Bet{
		Bet:  color,
		Type: TypeRedBlack,
	})

	return filters
}

func (c Controller) getNumber() int {
	bg := big.NewInt(36)

	number, _ := rand.Int(rand.Reader, bg)

	return int(number.Int64())
}

func (c Controller) getColor(number int) string {
	if number == 0 {
		return colorGreen
	}

	if number >= 1 && number <= 10 || number >= 19 && number <= 28 {
		if number%2 == 0 {
			return colorBlack
		}

		return colorRed
	}

	if number%2 == 0 {
		return colorRed
	}

	return colorBlack
}
