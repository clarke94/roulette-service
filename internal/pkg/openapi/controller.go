package openapi

import (
	"embed"
	"errors"

	"github.com/sirupsen/logrus"
)

var (
	// ErrDocumentation is the external error for unable to read html file.
	ErrDocumentation = errors.New("unable to render documentation")
	// ErrSpecification is the external error for unable to read swagger spec file.
	ErrSpecification = errors.New("unable to render specification")
)

//go:embed docs/*
var docs embed.FS

// Controller provides a domain controller.
type Controller struct {
	Logger *logrus.Logger

	swagger string
	docs    string
}

// New initializes a new Controller.
func New(logger *logrus.Logger) Controller {
	return Controller{
		Logger:  logger,
		swagger: "docs/swagger.json",
		docs:    "docs/redoc.html",
	}
}

// Specification returns the swagger specification.
func (c Controller) Specification() ([]byte, error) {
	data, err := docs.ReadFile(c.swagger)
	if err != nil {
		c.Logger.WithFields(logrus.Fields{
			"error": err.Error(),
		}).Error(ErrSpecification.Error())

		return nil, ErrSpecification
	}

	return data, nil
}

// Docs returns the html bytes for the redoc documentation.
func (c Controller) Docs() ([]byte, error) {
	data, err := docs.ReadFile(c.docs)
	if err != nil {
		c.Logger.WithFields(logrus.Fields{
			"error": err.Error(),
		}).Error(ErrDocumentation.Error())

		return nil, ErrDocumentation
	}

	return data, nil
}
