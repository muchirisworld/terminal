package middleware

import (
	"log/slog"
	"net/http"
	"strings"

	"github.com/clerk/clerk-sdk-go/v2"
	"github.com/clerk/clerk-sdk-go/v2/jwks"
	"github.com/clerk/clerk-sdk-go/v2/jwt"
	"github.com/muchirisworld/terminal/internal/auth"
	"github.com/muchirisworld/terminal/internal/config"
	"github.com/muchirisworld/terminal/internal/logger"
)

// AuthMiddleware creates a middleware that verifies Clerk JWT tokens
func AuthMiddleware(cfg *config.Config, _ *slog.Logger) func(next http.Handler) http.Handler {
	clientConfig := &clerk.ClientConfig{}
	clientConfig.Key = clerk.String(cfg.ClerkSecretKey)
	jwksClient := jwks.NewClient(clientConfig)

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				logger.Add(r.Context(), "auth_error", "missing authorization header")
				http.Error(w, "missing authorization header", http.StatusUnauthorized)
				return
			}

			parts := strings.Split(authHeader, " ")
			if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
				logger.Add(r.Context(), "auth_error", "invalid authorization header format")
				http.Error(w, "invalid authorization header", http.StatusUnauthorized)
				return
			}

			token := parts[1]

			unsafeClaims, err := jwt.Decode(r.Context(), &jwt.DecodeParams{
				Token: token,
			})
			if err != nil {
				logger.Add(r.Context(), "auth_error", "failed to decode token")
				logger.Add(r.Context(), "error", err.Error())
				http.Error(w, "Invalid or expired token", http.StatusUnauthorized)
				return
			}

			jwk, err := jwt.GetJSONWebKey(r.Context(), &jwt.GetJSONWebKeyParams{
				KeyID:      unsafeClaims.KeyID,
				JWKSClient: jwksClient,
			})
			if err != nil {
				logger.Add(r.Context(), "auth_error", "failed to get jwk")
				logger.Add(r.Context(), "error", err.Error())
				logger.Add(r.Context(), "kid", unsafeClaims.KeyID)
				http.Error(w, "Invalid or expired token", http.StatusUnauthorized)
				return
			}

			claims, err := jwt.Verify(r.Context(), &jwt.VerifyParams{
				Token: token,
				JWK:   jwk,
			})
			if err != nil {
				logger.Add(r.Context(), "auth_error", "failed to verify token")
				logger.Add(r.Context(), "error", err.Error())
				http.Error(w, "Invalid or expired token", http.StatusUnauthorized)
				return
			}

			orgID := claims.ActiveOrganizationID
			orgRole := claims.ActiveOrganizationRole

			if orgID == "" {
				logger.Add(r.Context(), "auth_error", "token missing active organization")
				logger.Add(r.Context(), "subject", claims.Subject)
				http.Error(w, "organization context required", http.StatusForbidden)
				return
			}

			logger.Add(r.Context(), "user_id", claims.Subject)
			logger.Add(r.Context(), "org_id", orgID)
			logger.Add(r.Context(), "role", orgRole)

			authCtx := &auth.AuthContext{
				UserID:  claims.Subject,
				OrgID:   orgID,
				OrgRole: orgRole,
			}

			ctx := auth.WithContext(r.Context(), authCtx)

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
