package providers

import (
	"context"
	"fmt"
	"os"

	"github.com/go-chi/jwtauth/v5"
)

type JWTProvider struct {
	jwt *jwtauth.JWTAuth
	
}

func NewJWTProvider() *JWTProvider {
	jwt := jwtauth.New("HS256", []byte(os.Getenv("JWT_SECRET")), nil)
	return &JWTProvider{jwt: jwt}
}

func (j *JWTProvider) Encode(claims map[string]interface{}) (string, error) {
	_, token, err := j.jwt.Encode(claims)
	if err != nil {
		return "", fmt.Errorf("error while encoding token: %w", err)
	}
	return token, nil
}

func (j *JWTProvider) Decode(t string) (map[string]interface{}, error) {
	token, err := j.jwt.Decode(t)
	if err != nil {
		return nil, fmt.Errorf("error while decoding token: %w", err)
	}
	claims, err := token.AsMap(context.Background())
	if err != nil {
		return nil, fmt.Errorf("error while converting token data to map: %w", err)
	}
	return claims, nil
}

func (j *JWTProvider) ValidateToken(tokenString string) (string, error) {
	claims, err := j.Decode(tokenString)
	if err != nil {
		return "", fmt.Errorf("error while decoding token: %w", err)
	}

	userID, ok := claims["id"].(string)
	if !ok {
		return "", fmt.Errorf("invalid token claims")
	}

	return userID, nil
}