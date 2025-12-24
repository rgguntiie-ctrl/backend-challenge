package meta

import (
	"errors"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
)

type MetaErrorHandler struct {
	// FiberFramework represents the instance of the Fiber web framework
	// used within the middleware for handling HTTP requests and responses.
	FiberFramework *fiberFramework
}

type fiberFramework struct {
	isLog bool
}

func NewMetaErrorHandler(optionFuncs ...func(*MetaErrorHandlerOptions)) *MetaErrorHandler {
	defaultOptions := getDefaultMetaErrorHandlerOptions()
	options := &defaultOptions
	for _, optionFunc := range optionFuncs {
		optionFunc(options)
	}
	return &MetaErrorHandler{
		FiberFramework: &fiberFramework{
			isLog: options.isLog,
		},
	}
}

// ErrorHandler is a middleware function for handling errors in the Fiber framework.
// It returns a function that takes a Fiber context and an error, and processes the error
// to generate an appropriate HTTP response.
func (m *fiberFramework) ErrorHandler() func(c *fiber.Ctx, err error) error {
	return func(c *fiber.Ctx, err error) error {
		// Status code defaults to 500
		code := fiber.StatusInternalServerError

		if metaErr, ok := IsMetaError(err); ok {
			return c.Status(metaErr.HttpStatus()).JSON(metaErr)
		}

		// Retrieve the custom status code if it's a *fiber.Error
		var e *fiber.Error
		if errors.As(err, &e) {
			code = e.Code
		}
		if code == fiber.StatusInternalServerError {
			metaErr := NewMetaError(-1000, "the server encountered an internal error or misconfiguration and was unable to complete your request", WithMetaErrorOptionsHttpStatus(code))
			metaErr.AppendError(err)

			if m.isLog {
				log.Errorf("error: %v", err)
			}
			return c.Status(metaErr.HttpStatus()).JSON(metaErr)
		}

		// Set Content-Type: text/plain; charset=utf-8
		c.Set(fiber.HeaderContentType, fiber.MIMETextPlainCharsetUTF8)

		// Return status code with error message
		return c.Status(code).SendString(err.Error())
	}
}
