package middlewares

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v4"
)

const ctxKeyClaims ctxKey = "claims"

type JWTParser interface {
	ParseToken(tokenString string) (jwt.MapClaims, error)
}

type JWTMiddleware struct {
	jwtParser     JWTParser
	requiredRoles []string
}

func NewJWTMiddleware(jwtParser JWTParser, requiredRoles []string) *JWTMiddleware {
	return &JWTMiddleware{jwtParser: jwtParser, requiredRoles: requiredRoles}
}

func (m *JWTMiddleware) VerifyToken(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {

		authHeader := r.Header.Get("Authorization")
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			http.Error(rw, "missing or invalid Authorization header", http.StatusUnauthorized)
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")

		claims, err := m.jwtParser.ParseToken(tokenString)
		if err != nil {
			http.Error(rw, fmt.Sprintf("invalid token: %v", err), http.StatusUnauthorized)
			return
		}

		roles, ok := claims["roles"].([]any)
		if !ok || !m.hasRequiredRole(roles) {
			http.Error(rw, "forbidden: insufficient permissions", http.StatusForbidden)
			return
		}

		ctx := context.WithValue(r.Context(), ctxKeyClaims, claims)
		next.ServeHTTP(rw, r.WithContext(ctx))
	})
}

func (m *JWTMiddleware) hasRequiredRole(roles []any) bool {
	for _, role := range roles {
		for _, reqRole := range m.requiredRoles {
			if role == reqRole {
				return true
			}
		}
	}
	return false
}
