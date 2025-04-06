package controller

import (
	"go_vdot_api/middleware"
	"go_vdot_api/model"
	"go_vdot_api/pkg/logger"
	"go_vdot_api/usecase"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
)



type IVdotController interface {
	CreateVdot(c echo.Context) error
	GetVdot(c echo.Context) error
	UpdateVdot(c echo.Context) error
	GetUserVdotValue(c echo.Context) error
}

type vdotController struct {
	vu usecase.IVdotUsecase
}

func NewVdotController(vu usecase.IVdotUsecase) IVdotController {
	return &vdotController{vu}
}

func (vc *vdotController) CreateVdot(c echo.Context) error {
	userClaims, err := middleware.GetUserClaims(c)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, err.Error())
	}

	vdot := model.Vdot{}
	if err := c.Bind(&vdot); err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	vdot.UserId = userClaims.UserID
	vdotRes, err := vc.vu.CreateVdot(vdot)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusCreated, vdotRes)
}

func (vc *vdotController) GetVdot(c echo.Context) error {
	userClaims, err := middleware.GetUserClaims(c)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, err.Error())
	}

	vdotRes, err := vc.vu.GetVdot(userClaims.UserID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	logger.Info("vdotRes", "vdotRes", vdotRes)
	return c.JSON(http.StatusOK, vdotRes)
}

func (vc *vdotController) UpdateVdot(c echo.Context) error {
	logger.Info("UpdateVdot")
	userClaims, err := middleware.GetUserClaims(c)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, err.Error())
	}

	id := c.Param("id")
	logger.Info("id", "id", id)
	vdotId, err := strconv.Atoi(id)
	if err != nil {
		return c.JSON(http.StatusBadRequest, "Invalid ID")
	}

	vdot := model.Vdot{}
	if err := c.Bind(&vdot); err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}
	vdot.UserId = userClaims.UserID
	logger.Info("vdot", "vdot", vdot)
	logger.Info("userClaims.UserID", "userClaims.UserID", userClaims.UserID)
	
	vdotRes, err := vc.vu.UpdateVdot(vdot , userClaims.UserID, uint(vdotId))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	logger.Info("vdotRes", "vdotRes", vdotRes)
	return c.JSON(http.StatusOK, vdotRes)
}

func (vc *vdotController) GetUserVdotValue(c echo.Context) error {
	userClaims, err := middleware.GetUserClaims(c)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, err.Error())
	}

	result, err := vc.vu.GetUserVdotValue(userClaims.UserID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, result)
}
