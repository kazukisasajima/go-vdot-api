package validator

import (
	"go_vdot_api/model"
	"regexp"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type IWorkoutValidator interface {
	WorkoutValidate(workout model.Workout) error
}

type workoutValidator struct{}

func NewWorkoutValidator() IWorkoutValidator {
	return &workoutValidator{}
}

func (wv *workoutValidator) WorkoutValidate(workout model.Workout) error {
	return validation.ValidateStruct(&workout,
		// 日付は必須
		validation.Field(
			&workout.Date,
			validation.Required.Error("date is required"),
		),

		// 開始時刻（"HH:mm"形式の正規表現）
		validation.Field(
			&workout.StartTime,
			validation.Required.Error("start_time is required"),
			validation.Match(regexp.MustCompile(`^\d{2}:\d{2}$`)).Error("start_time must be in HH:mm format"),
		),

		// 練習内容
		validation.Field(
			&workout.Workout,
			validation.Required.Error("workout is required"),
			validation.Length(3, 0).Error("workout must be at least 3 characters"),
		),

		// ラップタイム（null許容。指定された場合のみバリデーション）
		validation.Field(
			&workout.LapTime,
			validation.NilOrNotEmpty,
			validation.Match(regexp.MustCompile(`^\[.*\]$`)).Error("lap_time must be a JSON-like array format"),
		),

		// 練習距離（0以上）
		validation.Field(
			&workout.Mileage,
			validation.Required,
			validation.Min(0.0).Error("mileage must be 0 or more"),
		),

		// 距離単位（km または mile）
		validation.Field(
			&workout.MileageUnit,
			validation.Required,
			validation.In("km", "mile").Error("mileage_unit must be 'km' or 'mile'"),
		),

		// 天気（最大20文字）
		validation.Field(
			&workout.Weather,
			validation.Required,
			validation.Length(1, 20).Error("weather must be 1 to 20 characters"),
		),
	)
}
