package middleware

import (
	"context"
	"net/http"

	"github.com/google/uuid"
)

var contextKeyRequestID ContextKey = "contextKeyRequestID"

func GenerateRequestID() string {
	return uuid.New().String()
}

func LoggingMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		reqID, ok := GetRequestID(ctx)
		if !ok || reqID == "" {
			reqID = GenerateRequestID()
			ctx := withRequestID(ctx, reqID)
			r = r.WithContext(ctx)
		}
		next(w, r)
	}
}

func withRequestID(ctx context.Context, reqID string) context.Context {
	return context.WithValue(ctx, contextKeyRequestID, reqID)
}

func GetRequestID(ctx context.Context) (string, bool) {
	reqID, ok := ctx.Value(contextKeyRequestID).(string)
	return reqID, ok
}
