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

type IWorkoutController interface {
	CreateWorkout(c echo.Context) error
	GetWorkoutPerMonth(c echo.Context) error
	UpdateWorkout(c echo.Context) error
}

type workoutController struct {
	wu usecase.IWorkoutUsecase
}

func NewWorkoutController(wu usecase.IWorkoutUsecase) IWorkoutController {
	return &workoutController{wu}
}

func (wc *workoutController) CreateWorkout(c echo.Context) error {
	logger.Info("CreateWorkout called")
	userClaims, err := middleware.GetUserClaims(c)
	if err != nil {
		logger.Error("GetUserClaims error", err)
		return c.JSON(http.StatusUnauthorized, err.Error())
	}

	workout := model.Workout{}
	if err := c.Bind(&workout); err != nil {
		logger.Error("Bind error", err)
		return c.JSON(http.StatusBadRequest, err.Error())
	}
	logger.Info("workout", workout)

	workout.UserId = userClaims.UserID
	workoutRes, err := wc.wu.CreateWorkout(workout)
	if err != nil {
		logger.Error("CreateWorkout error", err)
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	logger.Info("workoutRes", workoutRes)
	return c.JSON(http.StatusCreated, workoutRes)
}

func (wc *workoutController) GetWorkoutPerMonth(c echo.Context) error {
	userClaims, err := middleware.GetUserClaims(c)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, err.Error())
	}

	yearStr := c.QueryParam("year")
	monthStr := c.QueryParam("month")

	year, err := strconv.Atoi(yearStr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, "invalid year format")
	}
	month, err := strconv.Atoi(monthStr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, "invalid month format")
	}	

	workoutRes, err := wc.wu.GetWorkoutPerMonth(userClaims.UserID, year, month)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusOK, workoutRes)
}

func (wc *workoutController) UpdateWorkout(c echo.Context) error {
	userClaims, err := middleware.GetUserClaims(c)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, err.Error())
	}

	workout := model.Workout{}
	if err := c.Bind(&workout); err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	workoutId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, "invalid workout ID")
	}

	workoutRes, err := wc.wu.UpdateWorkout(workout, userClaims.UserID, uint(workoutId))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusOK, workoutRes)
}