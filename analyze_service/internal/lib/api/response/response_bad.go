package response

import (
	"context"
	"errors"
	"net/http"

	requestcontext "memory-book/internal/app/middleware/request_context"

	"github.com/go-chi/render"
)

type RequestContext struct {
	Writer  http.ResponseWriter
	Request *http.Request
}

func ResponseBad(ctx context.Context, status int, message string) {

	r, ok := requestcontext.GetRequest(ctx)
	if !ok {
		return
	}
	w, ok := requestcontext.GetResponseWriter(ctx)
	if !ok {
		return
	}

	render.Status(r, status)
	render.JSON(w, r, message)
}

type StatusError interface {
	GetStatus() int
	GetStatusError() error
}

type InternalError interface {
	GetInternalError() error
}

type BadRequestError interface {
	GetBadRequestError() error
}

type NotFoundError interface {
	GetNotFoundError() error
}

type ForbiddenError interface {
	GetForbiddenError() error
}

type UnauthorizedError interface {
	GetUnauthorizedError() error
}

func SendBad(ctx context.Context, err error) {
	if e := errors.Unwrap(err); e != nil {
		SendBad(ctx, e)
	}

	if e, ok := err.(StatusError); ok {
		ResponseBad(ctx, e.GetStatus(), e.GetStatusError().Error())

		return
	}

	if _, ok := err.(InternalError); ok {
		ResponseBad(ctx, http.StatusInternalServerError, err.Error())

		return
	}

	if e, ok := err.(BadRequestError); ok {
		if wrapErr := errors.Unwrap(e.GetBadRequestError()); wrapErr != nil {
			status := getFirstErrorStatus(wrapErr)
			ResponseBad(ctx, status, e.GetBadRequestError().Error())

			return
		}

		ResponseBad(ctx, http.StatusBadRequest, e.GetBadRequestError().Error())

		return
	}

	if _, ok := err.(NotFoundError); ok {
		ResponseBad(ctx, http.StatusNotFound, err.Error())

		return
	}

	if e, ok := err.(ForbiddenError); ok {
		ResponseBad(ctx, http.StatusForbidden, e.GetForbiddenError().Error())

		return
	}
	if e, ok := err.(UnauthorizedError); ok {
		ResponseBad(ctx, http.StatusUnauthorized, e.GetUnauthorizedError().Error())

		return
	}

	if err != nil {
		//logger.Error("unhandled knowledgebase fetch error:", sl.Err(err))

		ResponseBad(ctx, http.StatusInternalServerError, err.Error())
	}

}

func getFirstErrorStatus(wrappError error) int {
	e := errors.Unwrap(wrappError)
	if e != nil {
		return getFirstErrorStatus(e)

	}

	switch err := wrappError.(type) {
	case StatusError:
		return err.GetStatus()

	case InternalError:
		return http.StatusInternalServerError

	case BadRequestError:
		return http.StatusBadRequest

	case NotFoundError:
		return http.StatusNotFound

	case ForbiddenError:
		return http.StatusForbidden
	case UnauthorizedError:
		return http.StatusUnauthorized

	default:
		return http.StatusBadRequest
	}
}
