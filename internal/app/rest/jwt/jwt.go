// Package jwt provides utilities for JWT token handling and validation.
package jwt

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

type Parser struct {
	config Config
}

func NewJWTParser(config Config) *Parser {
	return &Parser{config: config}
}

func (j *Parser) ParseToken(tokenString string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(j.config.SecretKey), nil
	})
	if err != nil {
		return nil, fmt.Errorf("invalid token: %w", err)
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return nil, fmt.Errorf("invalid token claims")
	}

	if !j.isTrustedIssuer(claims["iss"]) {
		return nil, fmt.Errorf("untrusted issuer")
	}

	if !j.isTokenValid(claims["exp"]) {
		return nil, fmt.Errorf("token expired")
	}

	return claims, nil
}

func (j *Parser) isTrustedIssuer(issuer any) bool {
	if issuerStr, ok := issuer.(string); ok {
		for _, trusted := range j.config.TrustedIssuers {
			if issuerStr == trusted {
				return true
			}
		}
	}
	return false
}

func (j *Parser) isTokenValid(exp any) bool {
	if expFloat, ok := exp.(float64); ok {
		return time.Now().Unix() <= int64(expFloat)
	}
	return false
}
