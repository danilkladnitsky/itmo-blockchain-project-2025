package error

import (
	"errors"
)

var ErrorInternalMessage = errors.New("internal error")

type InternalErr struct {
	Message error `json:"error" swaggertype:"string" example:"internal error"`
}

func (err InternalErr) GetInternalError() error {
	return err.Message
}

func (err InternalErr) Error() string {
	return err.Message.Error()
}

var ErrorBadRequestMessage = errors.New("bad request error")

type BadRequestErr struct {
	Message error `json:"error" swaggertype:"string" example:"bad request"`
}

func (err BadRequestErr) GetBadRequestError() error {
	return err.Message
}

func (err BadRequestErr) Error() string {
	return err.Message.Error()
}

type RequestErr struct {
	Message error
	Status  int
}

func (err RequestErr) GetStatus() int {
	return err.Status
}

func (err RequestErr) GetStatusError() error {
	return err.Message
}

func (err RequestErr) Error() string {
	return err.Message.Error()
}

type NotFoundErr struct {
	Message error
}

func (err NotFoundErr) Error() string {
	return err.Message.Error()
}

func (err NotFoundErr) GetNotFoundError() error {
	return err.Message
}

type ForbiddenErr struct {
	Message error `json:"error" swaggertype:"string" example:"forbidden request"`
}

func (err ForbiddenErr) Error() string {
	return err.Message.Error()
}

func (err ForbiddenErr) GetForbiddenError() error {
	return err.Message
}

type UnAuthorizeErr struct {
	Message error `json:"error" swaggertype:"string" example:"unauthorized request"`
}

func (err UnAuthorizeErr) Error() string {
	return err.Message.Error()
}

func (err UnAuthorizeErr) GetUnAuthorizeError() error {
	return err.Message
}
