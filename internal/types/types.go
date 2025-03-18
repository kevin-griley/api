package types

type contextKey string

const (
	ContextKeyRequestID contextKey = "requestID"
	ContextKeyUserID    contextKey = "userID"
	ContextKeyClaims    contextKey = "claims"
)
