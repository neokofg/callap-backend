package jwt

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/oklog/ulid/v2"
)

type Token struct {
	UserId    ulid.ULID
	ExpiresAt time.Time
	IssuedAt  time.Time
	Body      map[string]any
}

type Config struct {
	Secret          string
	AccessTokenTTL  int
	RefreshTokenTTL int
	Issuer          string
}

type Service struct {
	config Config
}

func NewService(config Config) *Service {
	return &Service{
		config: config,
	}
}

type jwtToken struct {
	UserID ulid.ULID      `json:"user_id"`
	Body   map[string]any `json:"body,omitempty"`
	jwt.RegisteredClaims
}

func (s *Service) GenerateToken(userID ulid.ULID, body map[string]any) (string, error) {
	ttl := s.config.AccessTokenTTL

	if tokenType, ok := body["type"].(string); ok && tokenType == "refresh" {
		ttl = s.config.RefreshTokenTTL
	}

	now := time.Now()
	tokenClaims := jwtToken{
		UserID: userID,
		Body:   body,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:   s.config.Issuer,
			IssuedAt: jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(
				now.Add(time.Duration(ttl) * time.Second),
			),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, tokenClaims)
	tokenString, err := token.SignedString([]byte(s.config.Secret))
	if err != nil {
		return "", fmt.Errorf("failed to sign token: %w", err)
	}

	return tokenString, nil
}

func (s *Service) ValidateToken(tokenString string) (Token, error) {
	token, err := jwt.ParseWithClaims(
		tokenString,
		&jwtToken{},
		func(token *jwt.Token) (any, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return []byte(s.config.Secret), nil
		},
	)
	if err != nil {
		return Token{}, fmt.Errorf("failed to parse token: %w", err)
	}

	if !token.Valid {
		return Token{}, errors.New("invalid token")
	}

	claims, ok := token.Claims.(*jwtToken)
	if !ok {
		return Token{}, errors.New("invalid token claims")
	}

	return Token{
		UserId:    claims.UserID,
		ExpiresAt: claims.ExpiresAt.Time,
		IssuedAt:  claims.IssuedAt.Time,
		Body:      claims.Body,
	}, nil
}

func (s *Service) RefreshToken(tokenString string) (string, error) {
	tokenClaims, err := s.ValidateToken(tokenString)
	if err != nil {
		return "", fmt.Errorf("failed to validate token: %w", err)
	}

	now := time.Now()
	if now.Sub(tokenClaims.IssuedAt) < time.Duration(s.config.RefreshTokenTTL)*time.Second {
		return s.GenerateToken(tokenClaims.UserId, tokenClaims.Body)
	}

	return "", errors.New("token is too old to refresh")
}

func (s *Service) GetTTL() (int, int) {
	return s.config.AccessTokenTTL, s.config.RefreshTokenTTL
}

func (s *Service) InvalidateToken(token string) error {
	return nil
}
