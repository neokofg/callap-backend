package context

import "github.com/neokofg/callap-backend/pkg/jwt"

type (
	JwtService interface {
		GenerateToken(userID string, body map[string]any) (string, error)
		ValidateToken(token string) (jwt.Token, error)
		RefreshToken(token string) (string, error)
		InvalidateToken(token string) error
		GetTTL() (int, int)
	}
)
