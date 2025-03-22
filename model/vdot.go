package model

import "time"

type Vdot struct {
	ID            uint      `json:"id" gorm:"primaryKey"`
	DistanceValue float64   `json:"distance_value"`
	DistanceUnit  string    `json:"distance_unit"`
	Time          string    `json:"time" gorm:"type:time"`
	Elevation     *float64  `json:"elevation"`   // NULL を許容するためポインタ型
	Temperature   *float64  `json:"temperature"` // NULL を許容するためポインタ型
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
	User          User      `json:"user" gorm:"foreignKey:UserId; constraint:OnDelete:CASCADE"`
	UserId        uint      `json:"user_id" gorm:"not null"`
}

type VdotResponse struct {
	ID            uint     `json:"id" gorm:"primaryKey"`
	DistanceValue float64  `json:"distance_value"`
	DistanceUnit  string   `json:"distance_unit"`
	Time          string   `json:"time" gorm:"type:time"`
	Elevation     *float64 `json:"elevation"`
	Temperature   *float64 `json:"temperature"`
}
