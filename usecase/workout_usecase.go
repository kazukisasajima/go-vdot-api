package usecase

import (
	"go_vdot_api/model"
	"go_vdot_api/pkg/logger"
	"go_vdot_api/repository"
	"go_vdot_api/validator"
)

type IWorkoutUsecase interface {
	CreateWorkout(workout model.Workout) (model.WorkoutResponse, error)
	GetWorkoutPerMonth(userId uint, year int, month int) ([]model.WorkoutResponse, error)
	UpdateWorkout(workout model.Workout, userId uint, workoutId uint) (model.WorkoutResponse, error)
}

type workoutUsecase struct {
	wr repository.IWorkoutRepository
	wv validator.IWorkoutValidator
}

func NewWorkoutUsecase(wr repository.IWorkoutRepository, wv validator.IWorkoutValidator) IWorkoutUsecase {
	return &workoutUsecase{wr, wv}
}

func (wu *workoutUsecase) CreateWorkout(workout model.Workout) (model.WorkoutResponse, error) {
	logger.Info("workout", "workout", workout)
	if err := wu.wv.WorkoutValidate(workout); err != nil {
		return model.WorkoutResponse{}, err
	}

	if err := wu.wr.CreateWorkout(&workout); err != nil {
		return model.WorkoutResponse{}, err
	}

	resWorkout := model.WorkoutResponse{
		ID:   workout.ID,
		Date: workout.Date,
	}
	return resWorkout, nil
}

func (wu *workoutUsecase) GetWorkoutPerMonth(userId uint, year int, month int) ([]model.WorkoutResponse, error) {
	workouts, err := wu.wr.GetWorkoutPerMonth(userId, year, month)
	if err != nil {
		return nil, err
	}
	resWorkout := make([]model.WorkoutResponse, len(workouts))
	for i, w := range workouts {
		resWorkout[i] = model.WorkoutResponse{
			ID:   w.ID,
			Date: w.Date,
			StartTime:    w.StartTime,
			Workout:      w.Workout,
			LapTime:      w.LapTime,
			Mileage:      w.Mileage,
			MileageUnit:  w.MileageUnit,
			Weather:      w.Weather,			
		}
	}
	return resWorkout, nil
}

func (wu *workoutUsecase) UpdateWorkout(workout model.Workout, userId uint, workoutId uint) (model.WorkoutResponse, error) {
	if err := wu.wv.WorkoutValidate(workout); err != nil {
		return model.WorkoutResponse{}, err
	}

	if err := wu.wr.UpdateWorkout(&workout, userId, workoutId); err != nil {
		return model.WorkoutResponse{}, err
	}

	resWorkout := model.WorkoutResponse{
		ID:   workout.ID,
		Date: workout.Date,
	}
	return resWorkout, nil
}
