package repository

import (
	"fmt"
	"go_vdot_api/model"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type IVdotRepository interface {
	CreateVdot(vdot *model.Vdot) error
	GetVdotByID(vdot *model.Vdot, userId uint, vdotId uint) error
	UpdateVdot(vdot *model.Vdot, userId uint, vdotId uint) error
}

type vdotRepository struct {
	db *gorm.DB
}

func NewVdotRepository(db *gorm.DB) IVdotRepository {
	return &vdotRepository{db}
}

func (vr *vdotRepository) CreateVdot(vdot *model.Vdot) error {
	if err := vr.db.Create(vdot).Error; err != nil {
		return err
	}
	return nil
}

func (vr *vdotRepository) GetVdotByID(vdot *model.Vdot, userId uint, vdotId uint) error {
	if err := vr.db.Joins("User").Where("user_id = ?", userId).First(vdot, vdotId).Error; err != nil {
		return err
	}
	return nil
}

func (vr *vdotRepository) UpdateVdot(vdot *model.Vdot, userId uint, vdotId uint) error {
	result := vr.db.Model(vdot).Clauses(clause.Returning{}).Where("id = ? AND user_id = ?", vdotId, userId).Updates(vdot)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected < 1 {
		return fmt.Errorf("object does not exist")
	}
	return nil
}
