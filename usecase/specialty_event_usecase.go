package usecase

import (
	"go_vdot_api/model"
	"go_vdot_api/repository"
	"go_vdot_api/validator"
)

type ISpecialtyEventUsecase interface {
	CreateSpecialtyEvent(specialtyEvent model.SpecialtyEvent) (model.SpecialtyEvent, error)
	GetSpecialtyEvent(userId uint) ([]model.SpecialtyEvent, error)
	UpdateSpecialtyEvent(specialtyEvent model.SpecialtyEvent, userId uint, specialtyEventId uint) (model.SpecialtyEvent, error)
}

type specialtyEventUsecase struct {
	ser repository.ISpecialtyEventRepository
	sev validator.ISpecialtyEventValidator
}

func NewSpecialtyEventUsecase(ser repository.ISpecialtyEventRepository, sev validator.ISpecialtyEventValidator) ISpecialtyEventUsecase {
	return &specialtyEventUsecase{ser, sev}
}

func (seu *specialtyEventUsecase) CreateSpecialtyEvent(specialtyEvent model.SpecialtyEvent) (model.SpecialtyEvent, error) {
	if err := seu.sev.SpecialtyEventValidate(specialtyEvent); err != nil {
		return model.SpecialtyEvent{}, err
	}

	if err := seu.ser.CreateSpecialtyEvent(&specialtyEvent); err != nil {
		return model.SpecialtyEvent{}, err
	}

	return specialtyEvent, nil
}

func (seu *specialtyEventUsecase) GetSpecialtyEvent(userId uint) ([]model.SpecialtyEvent, error) {
	specialtyEvents, err := seu.ser.GetSpecialtyEvent(userId)
	if err != nil {
		return nil, err
	}
	resSpecialtyEvents := make([]model.SpecialtyEvent, len(specialtyEvents))
	for i, se := range specialtyEvents {
		resSpecialtyEvents[i] = model.SpecialtyEvent{
			ID:         se.ID,
			EventName:  se.EventName,
			BestTime:   se.BestTime,
			RecordedAt: se.RecordedAt,
		}
	}
	return specialtyEvents, nil
}

func (seu *specialtyEventUsecase) UpdateSpecialtyEvent(specialtyEvent model.SpecialtyEvent, userId uint, specialtyEventId uint) (model.SpecialtyEvent, error) {
	if err := seu.sev.SpecialtyEventValidate(specialtyEvent); err != nil {
		return model.SpecialtyEvent{}, err
	}

	if err := seu.ser.UpdateSpecialtyEvent(&specialtyEvent, userId, specialtyEventId); err != nil {
		return model.SpecialtyEvent{}, err
	}

	resSpecialEvent := model.SpecialtyEvent{
		ID:         specialtyEvent.ID,
		EventName:  specialtyEvent.EventName,
		BestTime:   specialtyEvent.BestTime,
		RecordedAt: specialtyEvent.RecordedAt,
	}
	return resSpecialEvent, nil
}
