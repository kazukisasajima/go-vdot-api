package model

import (
	"time"
	"go_vdot_api/pkg"
)
	

type Workout struct {
	ID          uint      `json:"id" gorm:"primaryKey"`
	Date        pkg.DateOnly `json:"date"`         // 練習日（例：2023-10-01）
	StartTime   string    `json:"start_time"`   // 練習開始時刻（例："07:30"）
	Workout     string    `json:"workout"`      // 練習内容（例：E3.2km, 6x(I800m・レスト2分）, E3.2km）
	LapTime     *string    `json:"lap_time"`     // ラップタイム（例：[3:30, 3:40, 3:50]）
	Mileage     float64   `json:"mileage"`      // 練習距離（例：10, 20.2）
	MileageUnit string    `json:"mileage_unit"` // 練習距離の単位（例：km, mile）
	Weather     string    `json:"weather"`      // 天候（例：晴れ、曇り、雨）
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`

	User   User `json:"user" gorm:"foreignKey:UserId;constraint:OnDelete:CASCADE"`
	UserId uint `json:"user_id"`
}

type WorkoutResponse struct {
	ID          uint      `json:"id"`
	Date        pkg.DateOnly `json:"date"`         // 練習日（例：2023-10-01）
	StartTime   string    `json:"start_time"`   // 練習開始時刻（例："07:30"）
	Workout     string    `json:"workout"`      // 練習内容（例：E3.2km, 6x(I800m・レスト2分）, E3.2km）
	LapTime     *string    `json:"lap_time"`     // ラップタイム（例：[3:30, 3:40, 3:50]）
	Mileage     float64   `json:"mileage"`      // 練習距離（例：10, 20.2）
	MileageUnit string    `json:"mileage_unit"` // 練習距離の単位（例：km, mile）
	Weather     string    `json:"weather"`      // 天候（例：晴れ、曇り、雨）
}
