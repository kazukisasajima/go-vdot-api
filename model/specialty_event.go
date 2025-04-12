package model

import (
	"go_vdot_api/pkg"
	"time"

	"gorm.io/gorm"
)

type SpecialtyEvent struct {
	ID         uint           `json:"id" gorm:"primaryKey"`
	EventName  string         `json:"event_name" gorm:"type:varchar(20);not null"`
	BestTime   string         `json:"best_time" gorm:"type:varchar(15);not null"`
	RecordedAt pkg.DateOnly   `json:"recorded_at"`
	CreatedAt  time.Time      `json:"created_at"`
	UpdatedAt  time.Time      `json:"updated_at"`
	DeletedAt  gorm.DeletedAt `json:"deleted_at" gorm:"index"`

	User   User `json:"user" gorm:"foreignKey:UserId;constraint:OnDelete:CASCADE"`
	UserId uint `json:"user_id"`
}
