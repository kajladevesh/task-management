package usecase

import (
	"errors"
	"fmt"

	"github.com/golang-jwt/jwt/v4"
)



func (u *UserService) ValidateToken(tokenStr string) (string, error) {
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		return []byte(u.jwtSecret), nil
	})
	if err != nil || !token.Valid {
		return "", errors.New("invalid token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return "", errors.New("invalid token claims")
	}

	// user_id might be float64 if stored as number
	switch id := claims["user_id"].(type) {
	case string:
		return id, nil
	case float64:
		// convert number to string
		return fmt.Sprintf("%.0f", id), nil
	default:
		return "", errors.New("user_id not found in token")
	}
}
