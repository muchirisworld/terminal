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
)

// AuthMiddleware creates a middleware that verifies Clerk JWT tokens
func AuthMiddleware(cfg *config.Config, logger *slog.Logger) func(next http.Handler) http.Handler {
	clientConfig := &clerk.ClientConfig{}
	clientConfig.Key = clerk.String(cfg.ClerkSecretKey)
	jwksClient := jwks.NewClient(clientConfig)

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				logger.DebugContext(r.Context(), "missing authorization header")
				http.Error(w, "missing authorization header", http.StatusUnauthorized)
				return
			}

			parts := strings.Split(authHeader, " ")
			if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
				logger.DebugContext(r.Context(), "invalid authorization header format")
				http.Error(w, "invalid authorization header", http.StatusUnauthorized)
				return
			}

			token := parts[1]

			unsafeClaims, err := jwt.Decode(r.Context(), &jwt.DecodeParams{
				Token: token,
			})
			if err != nil {
				logger.WarnContext(r.Context(), "failed to decode token", "error", err)
				http.Error(w, "Invalid or expired token", http.StatusUnauthorized)
				return
			}

			jwk, err := jwt.GetJSONWebKey(r.Context(), &jwt.GetJSONWebKeyParams{
				KeyID:      unsafeClaims.KeyID,
				JWKSClient: jwksClient,
			})
			if err != nil {
				logger.WarnContext(r.Context(), "failed to get jwk", "error", err, "kid", unsafeClaims.KeyID)
				http.Error(w, "Invalid or expired token", http.StatusUnauthorized)
				return
			}

			claims, err := jwt.Verify(r.Context(), &jwt.VerifyParams{
				Token: token,
				JWK:   jwk,
			})
			if err != nil {
				logger.WarnContext(r.Context(), "failed to verify token", "error", err)
				http.Error(w, "Invalid or expired token", http.StatusUnauthorized)
				return
			}

			orgID := claims.ActiveOrganizationID
			orgRole := claims.ActiveOrganizationRole

			if orgID == "" {
				logger.WarnContext(r.Context(), "token missing active organization", "subject", claims.Subject)
				http.Error(w, "organization context required", http.StatusForbidden)
				return
			}

			logger.DebugContext(r.Context(), "authentication successful", "user_id", claims.Subject, "org_id", orgID, "role", orgRole)

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
