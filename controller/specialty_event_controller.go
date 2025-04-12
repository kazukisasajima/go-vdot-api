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

type ISpecialtyEventController interface {
	CreateSpecialtyEvent(c echo.Context) error
	GetSpecialtyEvent(c echo.Context) error
	UpdateSpecialtyEvent(c echo.Context) error
}

type specialtyEventController struct {
	seu usecase.ISpecialtyEventUsecase
}

func NewSpecialtyEventController(seu usecase.ISpecialtyEventUsecase) ISpecialtyEventController {
	return &specialtyEventController{seu}
}

func (sec *specialtyEventController) CreateSpecialtyEvent(c echo.Context) error {
	userClaims, err := middleware.GetUserClaims(c)
	if err != nil {
		logger.Error("GetUserClaims error: %v", err)
		return c.JSON(http.StatusUnauthorized, err.Error())
	}

	specialtyEvent := model.SpecialtyEvent{}
	if err := c.Bind(&specialtyEvent); err != nil {
		logger.Error("Bind error: %v", err)
		return c.JSON(http.StatusBadRequest, err.Error())
	}
	specialtyEvent.UserId = userClaims.UserID
	specialtyEventRes, err := sec.seu.CreateSpecialtyEvent(specialtyEvent)
	if err != nil {
		logger.Error("CreateSpecialtyEvent error: %v", err)
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusCreated, specialtyEventRes)
}

func (sec *specialtyEventController) GetSpecialtyEvent(c echo.Context) error {
	userClaims, err := middleware.GetUserClaims(c)
	if err != nil {
		logger.Error("GetUserClaims error: %v", err)
		return c.JSON(http.StatusUnauthorized, err.Error())
	}

	specialtyEvents, err := sec.seu.GetSpecialtyEvent(userClaims.UserID)
	if err != nil {
		logger.Error("GetSpecialtyEvent error: %v", err)
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusOK, specialtyEvents)
}

func (sec *specialtyEventController) UpdateSpecialtyEvent(c echo.Context) error {
	userClaims, err := middleware.GetUserClaims(c)
	if err != nil {
		logger.Error("GetUserClaims error: %v", err)
		return c.JSON(http.StatusUnauthorized, err.Error())
	}

	specialtyEvent := model.SpecialtyEvent{}
	if err := c.Bind(&specialtyEvent); err != nil {
		logger.Error("Bind error: %v", err)
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	specialtyEventId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		logger.Error("strconv.Atoi error: %v", err)
		return c.JSON(http.StatusBadRequest, "invalid id")
	}

	specialtyEventRes, err := sec.seu.UpdateSpecialtyEvent(specialtyEvent, userClaims.UserID, uint(specialtyEventId))
	if err != nil {
		logger.Error("UpdateSpecialtyEvent error: %v", err)
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusOK, specialtyEventRes)
}
