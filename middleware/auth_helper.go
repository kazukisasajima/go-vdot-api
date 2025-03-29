package middleware

import (
	"errors"
	"github.com/golang-jwt/jwt/v4"
	"github.com/labstack/echo/v4"
)

// UserClaims は JWT から取得するユーザー情報の構造体
type UserClaims struct {
	UserID uint   `json:"user_id"`
	Email  string `json:"email"`
	Name   string `json:"name"`
}

// GetUserClaims は echo.Context から JWT クレームを取得して UserClaims に変換する
func GetUserClaims(c echo.Context) (*UserClaims, error) {
	user, ok := c.Get("user").(*jwt.Token)
	if !ok {
		return nil, errors.New("invalid token format")
	}

	claims, ok := user.Claims.(jwt.MapClaims)
	if !ok || !user.Valid {
		return nil, errors.New("invalid token claims")
	}

	userIdFloat, ok := claims["user_id"].(float64)
	if !ok {
		return nil, errors.New("user_id not found in token")
	}

	return &UserClaims{
		UserID: uint(userIdFloat),
		Email:  claims["email"].(string),
		Name:   claims["name"].(string),
	}, nil
}
