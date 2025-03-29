package controller

import (
	"go_vdot_api/middleware"
	"go_vdot_api/model"
	"go_vdot_api/pkg/logger"
	"go_vdot_api/usecase"
	"net/http"
	"os"
	"time"

	"github.com/labstack/echo/v4"
)

type IUserController interface {
	SignUp(c echo.Context) error
	LogIn(c echo.Context) error
	LogOut(c echo.Context) error
	CsrfToken(c echo.Context) error

	UpdateUser(c echo.Context) error
	DeleteUser(c echo.Context) error
}

type userController struct {
	uu usecase.IUserUsecase
}

func NewUserController(uu usecase.IUserUsecase) IUserController {
	return &userController{uu}
}

func (uc *userController) SignUp(c echo.Context) error {
	user := model.User{}
	if err := c.Bind(&user); err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}
	userRes, err := uc.uu.SignUp(user)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusCreated, userRes)
}

func (uc *userController) LogIn(c echo.Context) error {
	user := model.User{}
	if err := c.Bind(&user); err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}
	tokenString, err := uc.uu.LogIn(user)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	cookie := new(http.Cookie)
	cookie.Name = "token"
	cookie.Value = tokenString
	cookie.Expires = time.Now().Add(24 * time.Hour)
	cookie.Path = "/"
	cookie.Domain = os.Getenv("API_DOMAIN")
	cookie.Secure = true
	cookie.HttpOnly = true
	cookie.SameSite = http.SameSiteNoneMode
	c.SetCookie(cookie)
	return c.NoContent(http.StatusOK)
}

func (uc *userController) LogOut(c echo.Context) error {
	cookie := new(http.Cookie)
	cookie.Name = "token"
	cookie.Value = ""
	cookie.Expires = time.Now()
	cookie.Path = "/"
	cookie.Domain = os.Getenv("API_DOMAIN")
	cookie.Secure = true
	cookie.HttpOnly = true
	cookie.SameSite = http.SameSiteNoneMode
	c.SetCookie(cookie)
	return c.NoContent(http.StatusOK)
}

func (uc *userController) CsrfToken(c echo.Context) error {
	token := c.Get("csrf").(string)
	return c.JSON(http.StatusOK, echo.Map{
		"csrf_token": token,
	})
}

func (uc *userController) UpdateUser(c echo.Context) error {
	logger.Info("テスト")
	userClaims, err := middleware.GetUserClaims(c)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, err.Error())
	}

	userData := model.User{}
	if err := c.Bind(&userData); err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	userData.ID = userClaims.UserID
	userRes, err := uc.uu.UpdateUser(userData)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	logger.Info("User updated successfully", userData)
	return c.JSON(http.StatusOK, userRes)
}

func (uc *userController) DeleteUser(c echo.Context) error {
	userClaims, err := middleware.GetUserClaims(c)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, err.Error())
	}

	if err := uc.uu.DeleteUser(userClaims.UserID); err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	return c.NoContent(http.StatusOK)
}
