package controller

import (
	"fmt"
	"go_vdot_api/model"
	"go_vdot_api/usecase"
	"net/http"
	"strconv"

	"github.com/golang-jwt/jwt/v4"
	"github.com/labstack/echo/v4"
)



type IVdotController interface {
	CreateVdot(c echo.Context) error
	GetVdotByID(c echo.Context) error
	UpdateVdot(c echo.Context) error
}

type vdotController struct {
	vu usecase.IVdotUsecase
}

func NewVdotController(vu usecase.IVdotUsecase) IVdotController {
	return &vdotController{vu}
}

func (vc *vdotController) CreateVdot(c echo.Context) error {
	user := c.Get("user").(*jwt.Token)
	fmt.Printf("User from JWT: %+v\n", user)
	claims := user.Claims.(jwt.MapClaims)
	userId := claims["user_id"].(float64)

	vdot := model.Vdot{}
	if err := c.Bind(&vdot); err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	vdot.UserId = uint(userId)
	vdotRes, err := vc.vu.CreateVdot(vdot)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusCreated, vdotRes)
}

func (vc *vdotController) GetVdotByID(c echo.Context) error {
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	userId := claims["user_id"].(float64)

	id := c.Param("id")
	vdotId, _ := strconv.Atoi(id)
	vdotRes, err := vc.vu.GetVdotByID(uint(userId), uint(vdotId))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusOK, vdotRes)
}

func (vc *vdotController) UpdateVdot(c echo.Context) error {
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	userId := claims["user_id"].(float64)

	// パスから id を取得（重要）
	id := c.Param("id")
	vdotId, err := strconv.Atoi(id)
	if err != nil {
		return c.JSON(http.StatusBadRequest, "Invalid ID")
	}

	vdot := model.Vdot{}
	if err := c.Bind(&vdot); err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}
	
	vdotRes, err := vc.vu.UpdateVdot(vdot ,uint(userId), uint(vdotId))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusOK, vdotRes)
}
