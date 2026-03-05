package auth

import "context"

type contextKey string

const authContextKey contextKey = "auth_context"

// AuthContext holds the verified identity of the current request.
type AuthContext struct {
	UserID    string
	OrgID     string
	SessionID string
	OrgRole   string
}

// WithContext attaches the AuthContext to the provided context.
func WithContext(ctx context.Context, authCtx *AuthContext) context.Context {
	return context.WithValue(ctx, authContextKey, authCtx)
}

// FromContext retrieves the AuthContext from the context.
func FromContext(ctx context.Context) (*AuthContext, bool) {
	authCtx, ok := ctx.Value(authContextKey).(*AuthContext)
	return authCtx, ok
}
