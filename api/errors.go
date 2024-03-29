package api

import (
	"fmt"
	"github.com/go-chi/render"
	"github.com/pkg/errors"
	"net/http"
)

type ErrResponse struct {
	Err            error `json:"-"` // low-level runtime error
	HTTPStatusCode int   `json:"-"` // http response status code

	StatusText string `json:"status"`          // user-level status message
	AppCode    int64  `json:"code,omitempty"`  // application-specific error code
	ErrorText  string `json:"error,omitempty"` // application-level error message, for debugging
}

func (err *ErrResponse) Render(_ http.ResponseWriter, req *http.Request) error {
	logEntry := GetLogEntry(req)
	logEntry.Error(
		"an error occurred while processing the request",
		"error", fmt.Sprintf("%+v", err.Err),
	)

	render.Status(req, err.HTTPStatusCode)

	return nil
}

func ErrInvalidRequest(err error) render.Renderer {
	return &ErrResponse{
		Err:            errors.Wrap(err, "Invalid Request"),
		HTTPStatusCode: 400,
		StatusText:     "Invalid Request",
		ErrorText:      err.Error(),
	}
}

func ErrInternalServerError(err error) render.Renderer {
	return &ErrResponse{
		Err:            errors.Wrap(err, "Internal Server Error"),
		HTTPStatusCode: 500,
		StatusText:     "Internal Server Error",
		ErrorText:      err.Error(),
	}
}
