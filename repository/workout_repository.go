package repository

import (
	"go_vdot_api/model"

	"gorm.io/gorm"
)

type IWorkoutRepository interface {
	CreateWorkout(workout *model.Workout) error
	GetWorkoutPerMonth(userId uint, year int, month int) ([]model.Workout, error)
	UpdateWorkout(workout *model.Workout, userId uint, workoutId uint) error
}

type workoutRepository struct {
	db *gorm.DB
}

func NewWorkoutRepository(db *gorm.DB) IWorkoutRepository {
	return &workoutRepository{db}
}

func (wr *workoutRepository) CreateWorkout(workout *model.Workout) error {
	if err := wr.db.Create(workout).Error; err != nil {
		return err
	}
	return nil
}

func (wr *workoutRepository) GetWorkoutPerMonth(userId uint, year int, month int) ([]model.Workout, error) {
	workouts := []model.Workout{}
	if err := wr.db.
		Where("user_id = ? AND YEAR(date) = ? AND MONTH(date) = ?", userId, year, month).
		Find(&workouts).Error; err != nil {
		return nil, err
	}
	return workouts, nil
}

func (wr *workoutRepository) UpdateWorkout(workout *model.Workout, userId uint, workoutId uint) error {
	result := wr.db.Model(workout).Where("id = ? AND user_id = ?", workoutId, userId).Updates(workout)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected < 1 {
		return gorm.ErrRecordNotFound
	}
	return nil
}
