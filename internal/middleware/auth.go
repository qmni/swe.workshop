package middleware

import (
	"context"
	"net/http"
	"os"
	"strings"

	keyfunc "github.com/MicahParks/keyfunc/v3"
	"github.com/golang-jwt/jwt/v5"
	"github.com/gofiber/fiber/v2"
)

var jwksKeyfunc jwt.Keyfunc

// InitKeycloak fetches the JWKS from Keycloak once at startup.
// If KEYCLOAK_JWKS_URL is not set, auth middleware is disabled.
func InitKeycloak() error {
	url := os.Getenv("KEYCLOAK_JWKS_URL")
	if url == "" {
		return nil
	}
	k, err := keyfunc.NewDefaultCtx(context.Background(), []string{url})
	if err != nil {
		return err
	}
	jwksKeyfunc = k.Keyfunc
	return nil
}

// RequireAuth validates Bearer JWT tokens issued by Keycloak.
// If KEYCLOAK_JWKS_URL was not set at startup, all requests pass through.
func RequireAuth(c *fiber.Ctx) error {
	if jwksKeyfunc == nil {
		return c.Next()
	}
	header := c.Get("Authorization")
	if !strings.HasPrefix(header, "Bearer ") {
		return c.Status(http.StatusUnauthorized).JSON(fiber.Map{"error": "missing bearer token"})
	}
	tokenStr := strings.TrimPrefix(header, "Bearer ")
	_, err := jwt.Parse(tokenStr, jwksKeyfunc, jwt.WithValidMethods([]string{"RS256"}))
	if err != nil {
		return c.Status(http.StatusUnauthorized).JSON(fiber.Map{"error": "invalid token"})
	}
	return c.Next()
}
