package middleware

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

// CheckAuth は JWT ミドルウェア後に呼び出されるログイン確認用ハンドラー
func CheckAuth(c echo.Context) error {
	userClaims, err := GetUserClaims(c)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, err.Error())
	}

	return c.JSON(http.StatusOK, echo.Map{
		"user": echo.Map{
			"id":    userClaims.UserID,
			"email": userClaims.Email,
			"name":  userClaims.Name,
		},
	})
}
