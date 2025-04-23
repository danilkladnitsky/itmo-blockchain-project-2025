package requestcontext

import (
	"context"
	"net/http"
)

const (
	CtxRequestKey  = "http_request"
	CtxResponseKey = "http_response"
)

// WithRequestContext добавляет http.Request и http.ResponseWriter в контекст.
func WithRequestContext(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := context.WithValue(r.Context(), CtxRequestKey, r)
		ctx = context.WithValue(ctx, CtxResponseKey, w)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// GetRequest изымает http.Request из контекста.
func GetRequest(ctx context.Context) (*http.Request, bool) {
	req, ok := ctx.Value(CtxRequestKey).(*http.Request)
	return req, ok
}

// GetResponseWriter изымает http.ResponseWriter из контекста.
func GetResponseWriter(ctx context.Context) (http.ResponseWriter, bool) {
	writer, ok := ctx.Value(CtxResponseKey).(http.ResponseWriter)
	return writer, ok
}
