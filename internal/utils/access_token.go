package utils

import (
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

const (
	accessTokenIssuer = "goAuth"
	accessTokenTTL    = 5 * time.Minute
)

type TokenService struct {
	secret []byte
}

type AuthenticatedUser struct {
	UserID string
	Role   string
}

type accessTokenClaims struct {
	UserID string `json:"userID"`
	Role   string `json:"role"`
	jwt.RegisteredClaims
}

func NewTokenService() (*TokenService, error) {
	secret := os.Getenv("ACCESS_TOKEN_SECRET_KEY")
	if len(secret) < 64 {
		return nil, errors.New("ACCESS TOKEN SECRET KEY must be at least 64 bytes")
	}

	return &TokenService{
		secret: []byte(secret),
	}, nil
}

func (a *TokenService) GenerateAccessToken(user AuthenticatedUser) (string, error) {
	if a == nil || len(a.secret) == 0 {
		return "", errors.New("access token and service is not initializable")
	}

	if user.UserID == "" || user.Role == "" {
		return "", errors.New("all fields are mandatory")
	}

	now := time.Now().UTC()

	claims := accessTokenClaims{
		UserID: user.UserID,
		Role:   user.Role,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    accessTokenIssuer,
			Subject:   user.UserID,
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(accessTokenTTL)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString(a.secret)
	if err != nil {
		return "", fmt.Errorf("sign access token: %w", err)
	}

	return signedToken, nil
}

func (a *TokenService) VerifyAccessToken(tokenString string) (*AuthenticatedUser, error) {
	if a == nil || len(a.secret) == 0 {
		return nil, errors.New("access token service is not initialized")
	}

	if tokenString == "" {
		return nil, errors.New("access token is required")
	}

	claims := &accessTokenClaims{}

	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (any, error) {
		if token.Method != jwt.SigningMethodHS256 {
			return nil, errors.New("unexpected signing method")
		}
		return a.secret, nil
	},
		jwt.WithValidMethods([]string{"HS256"}),
		jwt.WithIssuer("goAuth"),
		jwt.WithExpirationRequired(),
		jwt.WithIssuedAt(),
	)

	if err != nil {
		return nil, fmt.Errorf("verify access token: %w", err)
	}

	if token == nil || !token.Valid {
		return nil, errors.New("invalid access token")
	}

	if claims.UserID == "" || claims.Role == "" {
		return nil, errors.New("required user claims are missing")
	}

	if claims.Subject != claims.UserID {
		return nil, errors.New("token subject does not match user ID")
	}

	return &AuthenticatedUser{UserID: claims.UserID, Role: claims.Role}, nil
}
