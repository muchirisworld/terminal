package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/clerk/clerk-sdk-go/v2"
	"github.com/clerk/clerk-sdk-go/v2/jwks"
	"github.com/clerk/clerk-sdk-go/v2/jwt"
	"github.com/muchirisworld/terminal/internal/config"
)

type contextKey string

const clerkUserIDKey contextKey = "clerk_user_id"

// AuthMiddleware creates a middleware that verifies Clerk JWT tokens
func AuthMiddleware(cfg *config.Config) func(next http.Handler) http.Handler {
	clientConfig := &clerk.ClientConfig{}
	clientConfig.Key = clerk.String(cfg.ClerkSecretKey)
	jwksClient := jwks.NewClient(clientConfig)

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				http.Error(w, "missing authorization header", http.StatusUnauthorized)
				return
			}

			parts := strings.Split(authHeader, " ")
			if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
				http.Error(w, "invalid authorization header", http.StatusUnauthorized)
				return
			}

			token := parts[1]

			unsafeClaims, err := jwt.Decode(r.Context(), &jwt.DecodeParams{
				Token: token,
			})
			if err != nil {
				http.Error(w, "Invalid or expired token", http.StatusUnauthorized)
				return
			}

			jwk, err := jwt.GetJSONWebKey(r.Context(), &jwt.GetJSONWebKeyParams{
				KeyID:      unsafeClaims.KeyID,
				JWKSClient: jwksClient,
			})
			if err != nil {
				http.Error(w, "Invalid or expired token", http.StatusUnauthorized)
				return
			}

			claims, err := jwt.Verify(r.Context(), &jwt.VerifyParams{
				Token: token,
				JWK:   jwk,
			})
			if err != nil {
				http.Error(w, "Invalid or expired token", http.StatusUnauthorized)
				return
			}

			// Attach user ID to context
			ctx := context.WithValue(r.Context(), clerkUserIDKey, claims.Subject)

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// ClerkUserID retrieves the user ID from the context if it exists
func ClerkUserID(ctx context.Context) (string, bool) {
	id, ok := ctx.Value(clerkUserIDKey).(string)
	return id, ok
}
