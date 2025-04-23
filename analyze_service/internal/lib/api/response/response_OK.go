package response

import (
	"context"
	"net/http"

	requestcontext "memory-book/internal/app/middleware/request_context"

	"github.com/go-chi/render"
)

type ResponseData interface {
	GetResponseData() interface{}
}

// responseOK model info
type responseOK struct {
	Data interface{} `json:"data" format:"object"`
}

func ResponseOKMessage(ctx context.Context, message string) {
	r, ok := requestcontext.GetRequest(ctx)
	if !ok {
		return
	}
	w, ok := requestcontext.GetResponseWriter(ctx)
	if !ok {
		return
	}
	render.Status(r, http.StatusOK)
	render.JSON(w, r, message)
}

func ResponseOK(ctx context.Context, data ResponseData) {

	r, ok := requestcontext.GetRequest(ctx)
	if !ok {
		return
	}
	w, ok := requestcontext.GetResponseWriter(ctx)
	if !ok {
		return
	}
	render.Status(r, http.StatusOK)
	render.JSON(w, r, responseOK{
		Data: data.GetResponseData(),
	})
}

func ResponseOKWithStatus(ctx context.Context, status int, data ResponseData) {
	r, ok := requestcontext.GetRequest(ctx)
	if !ok {
		return
	}
	w, ok := requestcontext.GetResponseWriter(ctx)
	if !ok {
		return
	}
	render.Status(r, status)
	render.JSON(w, r, responseOK{
		Data: data.GetResponseData(),
	})
}

func HttpResponseOK(ctx context.Context, data any) {
	r, ok := requestcontext.GetRequest(ctx)
	if !ok {
		return
	}
	w, ok := requestcontext.GetResponseWriter(ctx)
	if !ok {
		return
	}
	render.Status(r, http.StatusOK)
	render.JSON(w, r, responseOK{
		Data: data,
	})
}
