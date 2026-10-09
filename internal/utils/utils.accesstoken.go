package utils

import (
	"errors"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type AccessTokenService struct {
	secret []byte
}

type AccessTokenClaims struct {
	UserID string `json:"user_id"`
	jwt.RegisteredClaims
}

func NewAccessTokenService() (*AccessTokenService, error) {
	secret := os.Getenv("ACCESS_TOKEN_SECRET_KEY")
	if len(secret) < 64 {
		return nil, errors.New("ACCESS_TOKEN_SECRET_KEY must be at least 64 bytes")
	}

	return &AccessTokenService{
		secret: []byte(secret),
	}, nil
}

func (s *AccessTokenService) GenerateAccessToken(userID string) (string, error) {
	if userID == "" {
		return "", errors.New("UserID cannot be empty")
	}

	now := time.Now()

	claims := AccessTokenClaims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			Subject: userID,
			IssuedAt: jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(15 * time.Minute)),
			Issuer: "goAuth",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodES256, claims)

	return token.SignedString(s.secret)
}

func (s *AccessTokenService) VerifyAccessToken(tokenString string) (*AccessTokenClaims, error) {
	claims := &AccessTokenClaims{}
	
	token, err := jwt.ParseWithClaims(
		tokenString,
		claims,
		func(token *jwt.Token) (any, error) {
			if token.Method != jwt.SigningMethodES256 {
				return nil, errors.New("unexpected signing algorithm")
			}
			return s.secret, nil
		},
		jwt.WithIssuer("goAuth"),
		jwt.WithExpirationRequired(),
		jwt.WithIssuedAt(),
	)

	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, errors.New("invalid access token")
	}

	if claims.UserID == "" || claims.Subject != claims.UserID {
		return nil, errors.New("invalid token claims")
	}

	return claims, nil
}
