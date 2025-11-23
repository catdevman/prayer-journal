package middleware

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	jwtmiddleware "github.com/auth0/go-jwt-middleware/v2"
	"github.com/auth0/go-jwt-middleware/v2/jwks"
	"github.com/auth0/go-jwt-middleware/v2/validator"
)

// EnsureValidToken is a middleware that validates the JWT in the Authorization header.
func EnsureValidToken() func(next http.Handler) http.Handler {
	issuerURL, err := url.Parse(os.Getenv("AUTH0_ISSUER"))
	if err != nil {
		panic(fmt.Sprintf("Failed to parse Auth0 Issuer URL: %v", err))
	}

	audience := os.Getenv("AUTH0_AUDIENCE")

	provider := jwks.NewCachingProvider(issuerURL, 5*time.Minute)

	jwtValidator, err := validator.New(
		provider.KeyFunc,
		validator.RS256,
		issuerURL.String(),
		[]string{audience},
		validator.WithCustomClaims(func() validator.CustomClaims {
			return &CustomClaims{}
		}),
		validator.WithAllowedClockSkew(time.Minute),
	)
	if err != nil {
		panic(fmt.Sprintf("Failed to set up JWT validator: %v", err))
	}

	errorHandler := func(w http.ResponseWriter, r *http.Request, err error) {
		slog.Error("JWT Validation Failed", "error", err)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte(`{"message": "Failed to validate JWT."}`))
	}

	middleware := jwtmiddleware.New(
		jwtValidator.ValidateToken,
		jwtmiddleware.WithErrorHandler(errorHandler),
	)

	return func(next http.Handler) http.Handler {
		return middleware.CheckJWT(next)
	}
}

// CustomClaims defines any custom data you expect in the token
type CustomClaims struct {
	Scope string `json:"scope"`
}

func (c *CustomClaims) Validate(ctx context.Context) error {
	return nil
}

// GetUserID extracts the 'sub' (Subject) claim from the context
func GetUserID(ctx context.Context) (string, error) {
	claims, ok := ctx.Value(jwtmiddleware.ContextKey).(*validator.ValidatedClaims)
	if !ok {
		return "", fmt.Errorf("no claims found in context")
	}
	return claims.RegisteredClaims.Subject, nil
}
