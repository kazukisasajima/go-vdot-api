package repository

import (
	"go_vdot_api/model"

	"gorm.io/gorm"
)

type ISpecialtyEventRepository interface {
	CreateSpecialtyEvent(specialtyEvent *model.SpecialtyEvent) error
	GetSpecialtyEvent(userId uint) ([]model.SpecialtyEvent, error)
	UpdateSpecialtyEvent(specialtyEvent *model.SpecialtyEvent, userId uint, specialtyEventId uint) error
}

type specialtyEventRepository struct {
	db *gorm.DB
}

func NewSpecialtyEventRepository(db *gorm.DB) ISpecialtyEventRepository {
	return &specialtyEventRepository{db}
}

// TODO usecaseなどでGetSpecialtyEventを呼び出すよりも、ここで作成後のデータを返す方が効率いいかも
func (ser *specialtyEventRepository) CreateSpecialtyEvent(specialtyEvent *model.SpecialtyEvent) error {
	if err := ser.db.Create(specialtyEvent).Error; err != nil {
		return err
	}
	return nil
}

func (ser *specialtyEventRepository) GetSpecialtyEvent(userId uint) ([]model.SpecialtyEvent, error) {
	specialtyEvents := []model.SpecialtyEvent{}
	if err := ser.db.
	Where("user_id = ?", userId).
	Find(&specialtyEvents).Error; err != nil {
		return nil, err
	}
	return specialtyEvents, nil
}

func (ser *specialtyEventRepository) UpdateSpecialtyEvent(specialtyEvent *model.SpecialtyEvent, userId uint, specialtyEventId uint) error {
	result := ser.db.Model(specialtyEvent).Where("id = ? AND user_id = ?", specialtyEventId, userId).Updates(specialtyEvent)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected < 1 {
		return gorm.ErrRecordNotFound
	}
	return nil
}
